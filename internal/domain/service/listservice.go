package service

import (
	"ShoppingList/internal/domain/entity"
	"context"
	"errors"
)

type ListStorage interface {
	Create(ctx context.Context, owner entity.User, list entity.List) (entity.List, error)
	GetByID(ctx context.Context, id int) (entity.List, error)
	GetByKey(ctx context.Context, key string) (entity.List, error)
	GetByUserID(ctx context.Context, id int) ([]entity.List, error)
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

func (s *listService) CreateList(ctx context.Context, user entity.User, list entity.List) (entity.List, error) {
	list, err := s.storage.Create(ctx, user, list)

	return list, err
}

func (s *listService) GetList(ctx context.Context, user entity.User, list entity.List) (entity.List, error) {
	list, err := s.storage.GetByKey(ctx, list.Key)

	return list, err
}

func (s *listService) GetLists(ctx context.Context, user entity.User) ([]entity.List, error) {
	lists, err := s.storage.GetByUserID(ctx, user.ID)

	return lists, err
}

func (s *listService) DeleteList(ctx context.Context, user entity.User, list entity.List) error {
	list, _ = s.storage.GetByKey(ctx, list.Key)
	if list.OwnerID != user.ID {
		return errors.New("access denied")
	}

	list.Status = "deleted"

	err := s.storage.UpdateList(ctx, list)

	return err
}

func (s *listService) CloseList(ctx context.Context, user entity.User, list entity.List) error {
	list.Status = "closed"

	err := s.storage.UpdateList(ctx, list)

	return err
}

func (s *listService) UpdateListItemStatus(ctx context.Context, user entity.User, list entity.List, item entity.Item) (entity.List, error) {
	list, _ = s.storage.GetByKey(ctx, list.Key)
	if list.Status != "opened" {
		return entity.List{}, errors.New("access denied")
	}

	var targetItem *entity.Item
	for i := 0; i < len(list.Items); i++ {
		if list.Items[i].ID == item.ID {
			targetItem = &list.Items[i]
			break
		}
	}

	if targetItem == nil {
		return entity.List{}, errors.New("item not found")
	}

	if targetItem.Status == "closed" {
		targetItem.Status = "opened"
	} else if targetItem.Status == "opened" {
		targetItem.Status = "closed"
	}

	err := s.storage.UpdateList(ctx, list)

	return list, err
}
