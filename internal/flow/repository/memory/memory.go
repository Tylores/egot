package memory

import (
	"slices"
	"sync"

	"github.com/Tylores/egot/internal/flow/repository"
	"github.com/Tylores/egot/internal/sep"
	"github.com/Tylores/egot/internal/uri"
)

type Entity uint32
type Entities = []Entity

type Pool struct {
	active []bool
	frq    []sep.FlowReservationRequest
	frp    []sep.FlowReservationResponse
	frprsp []sep.FlowReservationResponseResponse
}

func NewPool(size Entity) *Pool {
	return &Pool{
		active: make([]bool, size),
		frq:    make([]sep.FlowReservationRequest, size),
		frp:    make([]sep.FlowReservationResponse, size),
		frprsp: make([]sep.FlowReservationResponseResponse, size),
	}
}

type Repository struct {
	sync.RWMutex
	tag_lookup map[string]Entities
	pool       Pool
}

func NewRepository(size Entity) *Repository {
	return &Repository{
		tag_lookup: make(map[string]Entities),
		pool:       *NewPool(size),
	}
}

func (r *Repository) NextFreeEntity() (*Entity, error) {
	r.Lock()
	defer r.Unlock()

	for i, e := range r.pool.active {
		if e {
			continue
		}
		r.pool.active[i] = true
		entity := Entity(i)
		return &entity, nil
	}

	return nil, repository.ErrPoolFull
}

func (r *Repository) InitRepository(dir string) {
	// cant think of anything this repo needs in init yet
	println("Initializing ")
}

func (r *Repository) GetEntities(tag string) (Entities, error) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.tag_lookup[tag]
	if exists {
		return value, nil
	}
	return value, repository.ErrNotFound
}

func (r *Repository) TagEntity(tag string, e Entity) {
	r.Lock()
	defer r.Unlock()
	// entities will not overlap by controling active array
	r.tag_lookup[tag] = append(r.tag_lookup[tag], e)
}

func (r *Repository) FreeEntity(tag string, e Entity) error {
	r.Lock()
	defer r.Unlock()
	var entities Entities
	found := false
	for idx := range r.tag_lookup[tag] {
		if Entity(idx) != e {
			entities = append(entities, r.tag_lookup[tag][idx])
		} else {
			found = true
		}
	}
	if !found {
		return repository.ErrEntityNotTagged
	}
	r.tag_lookup[tag] = entities
	r.pool.active[e] = false
	return nil
}

func (r *Repository) GetFlowReservationRequests(lfdi string) (sep.FlowReservationRequestList, error) {
	r.Lock()
	defer r.Unlock()
	entities, err := r.GetEntities(lfdi)
	if err != nil {
		return sep.FlowReservationRequestList{}, repository.ErrNotFound
	}

	all := sep.UInt32(len(entities))
	result := all
	list := sep.NewList(sep.NewResource(uri.FlowReservationRequestList), all, result)
	frql := sep.NewFlowReservationRequestList(list)

	for _, e := range entities {
		frql.FlowReservationRequest = append(frql.FlowReservationRequest, &r.pool.frq[e])
	}

	return *frql, nil
}

func (r *Repository) GetFlowReservationRequest(lfdi string, e Entity) (sep.FlowReservationRequest, error) {
	r.Lock()
	defer r.Unlock()
	entities, err := r.GetEntities(lfdi)
	if err != nil {
		return sep.FlowReservationRequest{}, repository.ErrNotFound
	}

	if slices.Index(entities, e) == -1 {
		return sep.FlowReservationRequest{}, repository.ErrNotFound
	}

	return r.pool.frq[e], nil
}

func (r *Repository) PostFlowReservationRequest(lfdi string, data sep.FlowReservationRequest) (Entity, error) {
	r.Lock()
	defer r.Unlock()
	e, err := r.NextFreeEntity()
	if err != nil {
		return 0, repository.ErrPoolFull

	}
	r.TagEntity(lfdi, *e)
	r.pool.frq[*e] = data
	return *e, nil
}

func (r *Repository) PutFlowReservationRequest(id Entity, data sep.FlowReservationRequest) error {
	r.Lock()
	defer r.Unlock()
	if r.pool.frq[id].MRID != data.MRID {
		return repository.ErrWrongMRID
	}

	if r.pool.frq[id].RequestStatus.RequestStatus == data.RequestStatus.RequestStatus {
		return repository.ErrBadData
	}
	r.pool.active[id] = false
	return nil
}
