package memory

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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

		cert, error := os.ReadFile(path)
		if error != nil {
			return error
		}

		fp := fmt.Sprintf("%X", sha256.Sum256(cert))[0:40]
		e, error := r.NextFreeEntity()
		if error != nil {
			return error
		}

		fmt.Printf("\t%s : %d\n", fp, *e)
		error = r.TagEntity(fp, *e)
		if error != nil {
			return error
		}
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
