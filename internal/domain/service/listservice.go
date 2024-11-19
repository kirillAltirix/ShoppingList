package service

import (
	"ShoppingList/internal/domain/entity"
	"context"
)

type ListStorage interface {
	Create(ctx context.Context, owner entity.User, list entity.List) (entity.List, error)
	GetByID(ctx context.Context, id int) (entity.List, error)
	GetAll(ctx context.Context) ([]entity.List, error)
	UpdateList(ctx context.Context, list entity.List) error
	UpdateItem(ctx context.Context, item entity.Item) error
	AddSubowner(ctx context.Context, subowner entity.User, list entity.List) error
}

type listService struct {
	storage ListStorage
}

func NewListService(storage ListStorage) *listService {
	return &listService{storage}
}
