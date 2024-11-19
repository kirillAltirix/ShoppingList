package service

import (
	"ShoppingList/internal/domain/entity"
	"context"
)

type UserStorage interface {
	Create(ctx context.Context, user entity.User) (int, error)
	GetByChatID(ctx context.Context, chatID string) (entity.User, error)
}

type userService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) *userService {
	return &userService{storage}
}
