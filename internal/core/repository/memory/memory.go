package memory

import (
	"sync"

	"github.com/Tylores/egot/internal/core/repository"
	"github.com/Tylores/egot/internal/sep"
)

type Entity uint32

type Pool struct {
	dcap []sep.DeviceCapability
	edev []sep.EndDevice
	reg  []sep.Registration
}

func NewPool(size Entity) *Pool {
	return &Pool{
		dcap: make([]sep.DeviceCapability, size),
		edev: make([]sep.EndDevice, size),
		reg:  make([]sep.Registration, size),
	}
}

type Repository struct {
	sync.RWMutex
	tag_lookup map[string]Entity
	pool       Pool
}

func NewRepository(size Entity) *Repository {
	return &Repository{
		pool: *NewPool(size),
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

func (r *Repository) TagEntity(tag string, id Entity) error {
	r.Lock()
	defer r.Unlock()
	_, exists := r.tag_lookup[tag]
	if exists {
		return repository.ErrTagExists
	}
	r.tag_lookup[tag] = id
	return nil
}

func (r *Repository) GetDeviceCapability(id Entity) sep.DeviceCapability {
	r.RLock()
	defer r.RUnlock()
	return r.pool.dcap[id]
}

func (r *Repository) PutDeviceCapability(id Entity, data sep.DeviceCapability) {
	r.Lock()
	defer r.Unlock()
	r.pool.dcap[id] = data
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
