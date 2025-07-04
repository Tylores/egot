package memory

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Tylores/egot/internal/core/repository"
	"github.com/Tylores/egot/internal/sep"
	"github.com/Tylores/egot/internal/uri"
)

type Entity uint32

type Pool struct {
	dcap []sep.DeviceCapability
	edev []sep.EndDevice
	reg  []sep.Registration
}

func NewPool(size Entity) *Pool {
	return &Pool{
		edev: make([]sep.EndDevice, size),
		reg:  make([]sep.Registration, size),
	}
}

type Repository struct {
	sync.RWMutex
	tag_lookup      map[string]Entity
	active_entities []bool
	pool            Pool
}

func NewRepository(size Entity) *Repository {
	return &Repository{
		tag_lookup:      make(map[string]Entity),
		pool:            *NewPool(size),
		active_entities: make([]bool, size),
	}
}

func (r *Repository) NextFreeEntity() (*Entity, error) {
	r.Lock()
	defer r.Unlock()

	for i, e := range r.active_entities {
		if e {
			continue
		}
		r.active_entities[i] = true
		entity := Entity(i)
		return &entity, nil
	}

	return nil, repository.ErrPoolFull
}

func (r *Repository) InitRepository(dir string) {
	println("Initializing ")
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) != ".crt" {
			return nil
		}

		if !strings.Contains(path, "client") {
			return nil
		}

		cert_file, error := os.ReadFile(path)
		if error != nil {
			return error
		}

		block, _ := pem.Decode(cert_file)
		if block == nil {
			return nil
		}

		cert, error := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return error
		}

		lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[:40]
		e, error := r.NextFreeEntity()
		if error != nil {
			return error
		}

		fmt.Printf("\t%s : %d\n", lfdi, *e)
		error = r.TagEntity(lfdi, *e)
		if error != nil {
			return error
		}

		sfdi, error := sep.ToSFDI(lfdi)
		if error != nil {
			return error
		}

		href := strings.ReplaceAll(uri.EndDevice, "{id}", fmt.Sprintf("%d", *e))
		res := sep.NewResource(href)
		adev := sep.NewAbstractDevice(res, sep.SFDIType(sfdi))
		extd := sep.NewExternalDevice(adev)

		extd.FlowReservationRequestListLink = &sep.FlowReservationRequestListLink{
			ListLink: sep.NewListLink(
				sep.NewLink(uri.FlowReservationRequestList), 1),
		}
		extd.FlowReservationResponseListLink = &sep.FlowReservationResponseListLink{
			ListLink: sep.NewListLink(
				sep.NewLink(uri.FlowReservationResponseList), 1),
		}
		r.pool.edev[*e] = *sep.NewEndDevice(extd)

		href = strings.ReplaceAll(uri.Registration, "{id}", fmt.Sprintf("%d", *e))
		res = sep.NewResource(href)
		r.pool.edev[*e].RegistrationLink = &sep.RegistrationLink{
			Link: sep.NewLink(href),
		}

		dt := sep.GetTime()
		var pin sep.PINType = 123455
		r.pool.reg[*e] = *sep.NewRegistration(res, &dt, &pin)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Repository) GetEntity(tag string) (Entity, error) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.tag_lookup[tag]
	if exists {
		return value, nil
	}
	return value, repository.ErrNotFound
}

func (r *Repository) TagEntity(tag string, e Entity) error {
	r.Lock()
	defer r.Unlock()
	_, exists := r.tag_lookup[tag]
	if exists {
		return repository.ErrTagExists
	}
	r.tag_lookup[tag] = e
	return nil
}

func (r *Repository) GetDeviceCapability() sep.DeviceCapability {
	resource := sep.NewResource(uri.DeviceCapability)
	fsab := sep.NewFunctionSetAssignmentsBase(resource)
	fsab.TimeLink = &sep.TimeLink{
		Link: sep.NewLink(uri.Time),
	}
	dcap := sep.NewDeviceCapability(fsab)
	dcap.EndDeviceListLink = &sep.EndDeviceListLink{
		ListLink: sep.NewListLink(
			sep.NewLink(uri.EndDeviceList),
			1),
	}
	return *dcap
}

func (r *Repository) GetTime() sep.Time {
	resource := sep.NewResource(uri.Time)
	return *sep.NewTime(resource)

}

func (r *Repository) GetEndDevice(id Entity) sep.EndDevice {
	r.RLock()
	defer r.RUnlock()
	return r.pool.edev[id]
}

func (r *Repository) PutEndDevice(id Entity, data sep.EndDevice) {
	r.Lock()
	defer r.Unlock()
	r.pool.edev[id] = data
}

func (r *Repository) GetRegistration(id Entity) sep.Registration {
	r.RLock()
	defer r.RUnlock()
	return r.pool.reg[id]
}

func (r *Repository) PutRegistration(id Entity, data sep.Registration) {
	r.Lock()
	defer r.Unlock()
	r.pool.reg[id] = data
}
