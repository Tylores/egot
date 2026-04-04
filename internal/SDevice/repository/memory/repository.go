package memory

import "errors"

// Repository is an in-memory data store
type Repository struct {
entities map[string]interface{}
}

// NewRepository creates a new memory repository
func NewRepository() *Repository {
return &Repository{
entities: make(map[string]interface{}),
}
}

// GetEntity retrieves an entity by LFDI
func (r *Repository) GetEntity(lfdi string) (interface{}, error) {
if entity, ok := r.entities[lfdi]; ok {
return entity, nil
}
return nil, errors.New("not found")
}

// SetEntity stores an entity by LFDI
func (r *Repository) SetEntity(lfdi string, entity interface{}) {
r.entities[lfdi] = entity
}
