package main

import (
	"bytes"
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

type Emulator struct {
	cfg    Config
	client *http.Client
	state  State
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
		cfg:    c,
		client: client,
		state:  initialState,
	}, nil
}

func (e *Emulator) Run() {
	// 1. Discovery (Skipped for brevity, assuming direct paths work through gateway)
	
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
	// In a real setup, we would use the link from DCAP
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
		e.applyControl(c)
	}
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
