package service

import (
	"ShoppingList/internal/domain/entity"
	"context"
	"log/slog"
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

func (s *userService) CreateUser(ctx context.Context, user entity.User) error {
	_, err := s.storage.Create(ctx, user)
	return err
}

func (s *userService) GetUserByChatID(ctx context.Context, user entity.User) (entity.User, error) {
	user, err := s.storage.GetByChatID(ctx, user.ChatID)

	slog.Debug("GetUserByChatID", slog.Group("user", slog.Int("id", user.ID), slog.String("chat_id", user.ChatID)))

	return user, err
}
