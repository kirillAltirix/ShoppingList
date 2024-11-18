package postgre

import (
	"ShoppingList/internal/domain/entity"
	repository "ShoppingList/internal/repository/model"
	"context"
	"database/sql"
)

type postgreUserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *postgreUserStorage {
	return &postgreUserStorage{db}
}

func (s *postgreUserStorage) Create(ctx context.Context, user entity.User) (int, error) {
	q := "INSERT INTO users (chat_id, username) VALUES ($1, $2) RETURNING user_id"

	var id int
	err := s.db.QueryRowContext(ctx, q, user.ChatID, user.Username).Scan(&id)

	return id, err
}

func (s *postgreUserStorage) GetByChatID(ctx context.Context, chat_id string) (entity.User, error) {
	q := "SELECT * FROM users WHERE chat_id = $1"

	user := repository.User{}
	err := s.db.QueryRowContext(ctx, q, chat_id).Scan(&user.UserID, &user.ChatID, &user.Username)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{ID: user.UserID, ChatID: user.ChatID, Username: user.Username}, nil
}
