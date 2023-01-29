package handlers

import (
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
)

// CreateCollection is the struct used for that will implement the Handler interface for the creation
type CreateCollection struct {
	req     entity.CollectionDescription
	res     entity.Collection
	Service interfaces.IService
}

type DeleteCollection struct {
	Service interfaces.IService
}

// GetCollection In case some response type or sth similar is needed in the future
type GetCollection struct {
	res     entity.Collection
	Service interfaces.IService
}

// ListCollections In case some response type or sth similar is needed in the future
type ListCollections struct {
	res     []*entity.Collection
	Service interfaces.IService
}

// UpdateCollection represents the struct that implement the handler for update
type UpdateCollection struct {
	req     entity.CollectionDescription
	res     entity.Collection
	Service interfaces.IService
}
