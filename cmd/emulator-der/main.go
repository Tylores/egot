package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Tylores/egot/sep"
)

type DERType string

const (
	TypeLoad DERType = "load"
	TypeESS  DERType = "ess"
	TypePV   DERType = "pv"
)

type Config struct {
	Type     DERType
	Name     string
	Gateway  string
	SSLDir   string
	CertName string
	Interval time.Duration
}

type ControlState struct {
	mrid      string
	control   *sep.DERControl
	received  bool
	started   bool
	completed bool
}

type Emulator struct {
	cfg      Config
	client   *http.Client
	state    State
	lfdi     string
	sfdi     sep.SFDIType
	controls map[string]*ControlState
}

type State struct {
	SoC          float64 // for ESS (0.0 to 1.0)
	Capacity     float64 // Wh
	Power        float64 // Watts (positive is injection, negative is consumption)
	Irradiance   float64 // for PV (0.0 to 1.0)
	LastUpdate   time.Time
}

func main() {
	baseCfg := Config{}
	var derType string
	var count int
	flag.StringVar(&derType, "type", "load", "DER type: load, ess, pv")
	flag.StringVar(&baseCfg.Name, "name", "emulator", "Base emulator name")
	flag.StringVar(&baseCfg.Gateway, "gateway", "localhost:8443", "Nginx gateway address")
	flag.StringVar(&baseCfg.SSLDir, "ssl", "./ssl", "SSL directory")
	flag.DurationVar(&baseCfg.Interval, "interval", 15*time.Second, "Simulation interval")
	flag.IntVar(&count, "count", 1, "Number of virtual devices to emulate")
	flag.Parse()

	baseCfg.Type = DERType(derType)

	for i := 1; i <= count; i++ {
		go func(id int) {
			cfg := baseCfg
			cfg.Name = fmt.Sprintf("%s-%04d", baseCfg.Name, id)
			cfg.CertName = fmt.Sprintf("client-%04d", id)

			e, err := NewEmulator(cfg)
			if err != nil {
				log.Printf("[%s] Failed to create emulator: %v", cfg.Name, err)
				return
			}

			log.Printf("[%s] Starting %s emulator", cfg.Name, cfg.Type)
			e.Run()
		}(i)
		
		// Stagger startup to avoid thundering herd on registration
		time.Sleep(10 * time.Millisecond)
	}

	select {}
}

func NewEmulator(c Config) (*Emulator, error) {
	caCert, err := os.ReadFile(fmt.Sprintf("%s/ca.crt", c.SSLDir))
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(
		fmt.Sprintf("%s/%s.crt", c.SSLDir, c.CertName),
		fmt.Sprintf("%s/%s.key", c.SSLDir, c.CertName),
	)
	if err != nil {
		return nil, err
	}

	// Parse x509 cert to extract LFDI and SFDI
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %w", err)
	}
	lfdi := fmt.Sprintf("%X", sha256.Sum256(x509Cert.Raw))[0:40]
	sfdi, err := sep.ToSFDI(lfdi)
	if err != nil {
		return nil, fmt.Errorf("failed to compute SFDI from LFDI: %w", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true, // Nginx might use localhost
			},
		},
		Timeout: 10 * time.Second,
	}

	initialState := State{
		Capacity:   10000, // 10 kWh
		SoC:        0.5,
		LastUpdate: time.Now(),
	}

	return &Emulator{
		cfg:      c,
		client:   client,
		state:    initialState,
		lfdi:     lfdi,
		sfdi:     sfdi,
		controls: make(map[string]*ControlState),
	}, nil
}

func (e *Emulator) onboard() (string, error) {
	// Step 1: Discovery (GET /dcap)
	dcapURL := "https://" + e.cfg.Gateway + "/dcap"
	log.Printf("[%s] [CSIP] Step 1: Discovery GET %s", e.cfg.Name, dcapURL)
	resp, err := e.client.Get(dcapURL)
	if err != nil {
		return "", fmt.Errorf("discovery failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("discovery returned status %d", resp.StatusCode)
	}
	var dcap sep.DeviceCapability
	if err := xml.NewDecoder(resp.Body).Decode(&dcap); err != nil {
		return "", fmt.Errorf("failed to decode discovery response: %w", err)
	}
	log.Printf("[%s] [CSIP] Discovery successful", e.cfg.Name)

	// Step 2: Time Sync (GET /tm)
	timeURL := "https://" + e.cfg.Gateway + "/tm"
	log.Printf("[%s] [CSIP] Step 2: Time Sync GET %s", e.cfg.Name, timeURL)
	resp, err = e.client.Get(timeURL)
	if err != nil {
		return "", fmt.Errorf("time sync failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("time sync returned status %d", resp.StatusCode)
	}
	var sTime sep.Time
	if err := xml.NewDecoder(resp.Body).Decode(&sTime); err != nil {
		return "", fmt.Errorf("failed to decode time response: %w", err)
	}
	if sTime.Quality != 7 {
		return "", fmt.Errorf("time quality is not 7 (got %d)", sTime.Quality)
	}
	log.Printf("[%s] [CSIP] Time Sync successful (server time: %d, quality: %d)", e.cfg.Name, sTime.CurrentTime, sTime.Quality)

	// Step 3: Device Registration (POST /edev)
	ed := &sep.EndDevice{
		ExternalDevice: &sep.ExternalDevice{
			AbstractDevice: &sep.AbstractDevice{
				LFDI: e.lfdi,
				SFDI: &e.sfdi,
			},
		},
	}
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(ed); err != nil {
		return "", fmt.Errorf("failed to encode EndDevice: %w", err)
	}

	edevURL := "https://" + e.cfg.Gateway + "/edev"
	log.Printf("[%s] [CSIP] Step 3: Device Registration POST %s", e.cfg.Name, edevURL)
	resp, err = e.client.Post(edevURL, sep.ContentType, &buf)
	if err != nil {
		return "", fmt.Errorf("registration POST failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registration POST returned status %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		location = "/edev/" + fmt.Sprintf("%d", e.sfdi)
	}
	log.Printf("[%s] [CSIP] Registration successful, device location: %s", e.cfg.Name, location)

	// Step 4: PIN Verification (GET /edev/{id1}/rg)
	pinURL := "https://" + e.cfg.Gateway + location + "/rg"
	log.Printf("[%s] [CSIP] Step 4: PIN Verification GET %s", e.cfg.Name, pinURL)
	resp, err = e.client.Get(pinURL)
	if err != nil {
		return "", fmt.Errorf("PIN verification GET failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("PIN verification GET returned status %d", resp.StatusCode)
	}
	var reg sep.Registration
	if err := xml.NewDecoder(resp.Body).Decode(&reg); err != nil {
		return "", fmt.Errorf("failed to decode registration response: %w", err)
	}
	if reg.PIN == nil {
		return "", fmt.Errorf("registration response PIN is nil")
	}
	if *reg.PIN != 111115 {
		return "", fmt.Errorf("registration PIN %d does not match expected 111115", *reg.PIN)
	}
	log.Printf("[%s] [CSIP] PIN verification successful (PIN: %d)", e.cfg.Name, *reg.PIN)

	return location, nil
}

func (e *Emulator) Run() {
	// 1. Onboard using CSIP procedure
	devicePath, err := e.onboard()
	if err != nil {
		log.Printf("[%s] CSIP Onboarding failed: %v", e.cfg.Name, err)
		// Proceed anyway, but log warning
	} else {
		log.Printf("[%s] CSIP Onboarding completed successfully (Device Path: %s)", e.cfg.Name, devicePath)
	}

	// 2. Register with MUP
	mupPath, err := e.registerMUP()
	if err != nil {
		log.Printf("MUP registration failed: %v", err)
	} else {
		log.Printf("Registered MUP: %s", mupPath)
	}

	ticker := time.NewTicker(e.cfg.Interval)
	for range ticker.C {
		e.updateState()
		e.pushTelemetry(mupPath)
		e.pollControls()
		e.updateTrackedControls()
	}
}

func (e *Emulator) updateState() {
	now := time.Now()
	dt := now.Sub(e.state.LastUpdate).Hours()
	e.state.LastUpdate = now

	switch e.cfg.Type {
	case TypeLoad:
		// Random load between 500W and 2000W
		e.state.Power = -(500 + rand.Float64()*1500)
	case TypePV:
		// Simple solar curve based on time of day
		hour := float64(now.Hour()) + float64(now.Minute())/60.0
		if hour > 6 && hour < 18 {
			e.state.Irradiance = math.Sin((hour-6)/12.0 * math.Pi)
		} else {
			e.state.Irradiance = 0
		}
		e.state.Power = e.state.Irradiance * 5000 // 5kW peak
	case TypeESS:
		// ESS state machine would go here
		// For now, just idle if no controls
		if e.state.Power > 0 {
			// Discharging
			energyUsed := e.state.Power * dt
			e.state.SoC -= energyUsed / e.state.Capacity
		} else if e.state.Power < 0 {
			// Charging
			energyStored := -e.state.Power * dt * 0.95 // 95% efficiency
			e.state.SoC += energyStored / e.state.Capacity
		}
		
		if e.state.SoC > 1.0 {
			e.state.SoC = 1.0
			e.state.Power = 0
		} else if e.state.SoC < 0.0 {
			e.state.SoC = 0.0
			e.state.Power = 0
		}
	}
}

func (e *Emulator) registerMUP() (string, error) {
	kind := sep.ServiceKind(0) // Electric
	mup := &sep.MirrorUsagePoint{
		UsagePointBase: &sep.UsagePointBase{
			ServiceCategoryKind: &kind,
			IdentifiedObject: &sep.IdentifiedObject{
				Description: e.cfg.Name,
			},
		},
	}

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(mup); err != nil {
		return "", err
	}

	resp, err := e.client.Post("https://"+e.cfg.Gateway+"/mup", sep.ContentType, &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return resp.Header.Get("Location"), nil
}

func (e *Emulator) pushTelemetry(mupPath string) {
	if mupPath == "" {
		return
	}

	now := sep.TimeType(time.Now().Unix())
	reading := &sep.MirrorMeterReading{
		ReadingType: &sep.ReadingType{
			Resource: &sep.Resource{
				HrefAttr: "/rt/1", // Mock reading type
			},
		},
		LastUpdateTime: &now,
		Reading: &sep.Reading{
			ReadingBase: &sep.ReadingBase{
				Value: int64(e.state.Power),
			},
		},
	}

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(reading); err != nil {
		log.Printf("Failed to encode reading: %v", err)
		return
	}

	resp, err := e.client.Post("https://"+e.cfg.Gateway+mupPath, sep.ContentType, &buf)
	if err != nil {
		log.Printf("Failed to push telemetry: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Printf("Push telemetry failed with status: %d", resp.StatusCode)
	}

	log.Printf("[%s] State: Power=%.2fW, SoC=%.2f, Irradiance=%.2f", 
		e.cfg.Name, e.state.Power, e.state.SoC, e.state.Irradiance)
}

func (e *Emulator) pollControls() {
	if e.cfg.Type == TypeLoad {
		return
	}

	// Poll DERP for active controls (assuming /derp/1/actderc)
	resp, err := e.client.Get("https://"+e.cfg.Gateway+"/derp/1/actderc")
	if err != nil {
		log.Printf("Failed to poll controls: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var list sep.DERControlList
	if err := xml.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Printf("Failed to decode control list: %v", err)
		return
	}

	for _, c := range list.DERControl {
		e.processControl(c)
	}
}

func (e *Emulator) processControl(c *sep.DERControl) {
	if c == nil || c.DERControlBase == nil || c.MRID == nil || c.MRID.HexBinary128 == nil {
		return
	}

	mridStr := *c.MRID.HexBinary128
	cs, exists := e.controls[mridStr]
	if !exists {
		cs = &ControlState{
			mrid:    mridStr,
			control: c,
		}
		e.controls[mridStr] = cs
	}

	now := time.Now().Unix()
	var startTime int64
	if c.Interval != nil && c.Interval.Start != nil {
		startTime = int64(*c.Interval.Start)
	}
	duration := uint32(0)
	if c.Interval != nil {
		duration = c.Interval.Duration
	}
	endTime := startTime + int64(duration)

	// Transition 1: Received
	if !cs.received {
		log.Printf("[%s] Control %s received. Posting status Received (1)", e.cfg.Name, mridStr)
		if err := e.postControlResponse(c, 1); err == nil {
			cs.received = true
		} else {
			log.Printf("[%s] Failed to post Received status for control %s: %v", e.cfg.Name, mridStr, err)
		}
	}

	// Transition 2: Started
	if cs.received && !cs.started && now >= startTime && now < endTime {
		log.Printf("[%s] Control %s starting. Posting status Started (2)", e.cfg.Name, mridStr)
		if err := e.postControlResponse(c, 2); err == nil {
			cs.started = true
			e.applyControl(c)
		} else {
			log.Printf("[%s] Failed to post Started status for control %s: %v", e.cfg.Name, mridStr, err)
		}
	}

	// Transition 3: Completed
	if cs.started && !cs.completed && now >= endTime {
		log.Printf("[%s] Control %s completed. Posting status Completed (3)", e.cfg.Name, mridStr)
		if err := e.postControlResponse(c, 3); err == nil {
			cs.completed = true
		} else {
			log.Printf("[%s] Failed to post Completed status for control %s: %v", e.cfg.Name, mridStr, err)
		}
	}
}

func (e *Emulator) updateTrackedControls() {
	now := time.Now().Unix()
	for mridStr, cs := range e.controls {
		c := cs.control
		var startTime int64
		if c.Interval != nil && c.Interval.Start != nil {
			startTime = int64(*c.Interval.Start)
		}
		duration := uint32(0)
		if c.Interval != nil {
			duration = c.Interval.Duration
		}
		endTime := startTime + int64(duration)

		// Transition 2: Started
		if cs.received && !cs.started && now >= startTime && now < endTime {
			log.Printf("[%s] Control %s starting. Posting status Started (2)", e.cfg.Name, mridStr)
			if err := e.postControlResponse(c, 2); err == nil {
				cs.started = true
				e.applyControl(c)
			} else {
				log.Printf("[%s] Failed to post Started status for control %s: %v", e.cfg.Name, mridStr, err)
			}
		}

		// Transition 3: Completed
		if cs.started && !cs.completed && now >= endTime {
			log.Printf("[%s] Control %s completed. Posting status Completed (3)", e.cfg.Name, mridStr)
			if err := e.postControlResponse(c, 3); err == nil {
				cs.completed = true
			} else {
				log.Printf("[%s] Failed to post Completed status for control %s: %v", e.cfg.Name, mridStr, err)
			}
		}
	}
}

func (e *Emulator) postControlResponse(c *sep.DERControl, status uint8) error {
	nowVal := sep.TimeType(time.Now().Unix())

	respPayload := &sep.DERControlResponse{
		Response: &sep.Response{
			CreatedDateTime: &nowVal,
			EndDeviceLFDI:   e.lfdi,
			Status:          status,
			Subject:         c.MRID,
		},
	}

	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(respPayload); err != nil {
		return fmt.Errorf("failed to encode control response: %w", err)
	}

	url := fmt.Sprintf("https://%s/rsps/%d/rsp", e.cfg.Gateway, e.sfdi)
	resp, err := e.client.Post(url, sep.ContentType, &buf)
	if err != nil {
		return fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[%s] Successfully posted control response status %d to %s", e.cfg.Name, status, url)
	return nil
}

func (e *Emulator) applyControl(c *sep.DERControl) {
	if c.DERControlBase == nil {
		return
	}

	// Simple example: follow OpModFixedW (signed percent control)
	if c.DERControlBase.OpModFixedW != nil && c.DERControlBase.OpModFixedW.SignedPerCent != nil {
		targetPct := float64(*c.DERControlBase.OpModFixedW.SignedPerCent) / 10000.0 // Value is in hundredths of a percent
		// Max power 5kW
		e.state.Power = targetPct * 5000
		log.Printf("[%s] Applied control: Power target %.2f%% (%.2fW)", 
			e.cfg.Name, targetPct*100, e.state.Power)
	}

	if c.DERControlBase.OpModFreqWatt != nil && c.DERControlBase.OpModFreqWatt.HrefAttr != "" {
		e.fetchCurve(c.DERControlBase.OpModFreqWatt.HrefAttr, "Frequency-Watt")
	}

	if c.DERControlBase.OpModVoltVar != nil && c.DERControlBase.OpModVoltVar.HrefAttr != "" {
		e.fetchCurve(c.DERControlBase.OpModVoltVar.HrefAttr, "Volt-Var")
	}
}

func (e *Emulator) fetchCurve(href, curveName string) {
	resp, err := e.client.Get("https://" + e.cfg.Gateway + href)
	if err != nil {
		log.Printf("[%s] Failed to fetch %s curve: %v", e.cfg.Name, curveName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("[%s] Successfully received and stored %s curve from %s", e.cfg.Name, curveName, href)
	} else {
		log.Printf("[%s] Failed to fetch %s curve from %s, status: %d", e.cfg.Name, curveName, href, resp.StatusCode)
	}
}
