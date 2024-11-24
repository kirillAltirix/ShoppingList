package service

import (
	"ShoppingList/internal/domain/entity"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
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
	op := "service.listService.CreateList"

	list.Key = uuid.NewString()
	list.Status = "opened"
	for i := 0; i < len(list.Items); i++ {
		list.Items[i].Status = "opened"
	}
	list, err := s.storage.Create(ctx, user, list)

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return list, err
}

func (s *listService) GetList(ctx context.Context, user entity.User, list entity.List) (entity.List, error) {
	op := "service.listService.GetList"
	list, err := s.storage.GetByKey(ctx, list.Key)

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return list, err
}

func (s *listService) GetLists(ctx context.Context, user entity.User) ([]entity.List, error) {
	op := "service.listService.GetLists"
	lists, err := s.storage.GetByUserID(ctx, user.ID)

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return lists, err
}

func (s *listService) DeleteList(ctx context.Context, user entity.User, list entity.List) error {
	op := "service.listService.DeleteList"
	list, _ = s.storage.GetByKey(ctx, list.Key)
	if list.OwnerID != user.ID {
		return errors.New("access denied")
	}

	list.Status = "deleted"

	err := s.storage.UpdateList(ctx, list)

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return err
}

func (s *listService) CloseList(ctx context.Context, user entity.User, list entity.List) error {
	op := "service.listService.CloseList"

	list, _ = s.storage.GetByKey(ctx, list.Key)
	if list.Status == "deleted" {
		return errors.New("access denied")
	}
	list.Status = "closed"

	err := s.storage.UpdateList(ctx, list)

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return err
}

func (s *listService) UpdateListItemStatus(ctx context.Context, user entity.User, list entity.List, item entity.Item) (entity.List, error) {
	op := "service.listService.UpdateListItemStatus"
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

	if err != nil {
		slog.Error("error happened", slog.String("error message", err.Error()), slog.String("trace", op))
	}

	return list, err
}
