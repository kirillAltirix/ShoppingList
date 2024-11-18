package postgre

import (
	"ShoppingList/internal/domain/entity"
	"context"
	"database/sql"
)

// type ListStorage interface {
// 	Create(ctx context.Context, owner entity.User, list entity.List) error
// 	GetByID(ctx context.Context, id int) (entity.List, error)
// 	GetAll(ctx context.Context, owner entity.User) ([]entity.List, error)
// 	Delete(ctx context.Context, list entity.List) error
// 	UpdateList(ctx context.Context, list entity.List) error
// 	UpdateItem(ctx context.Context, item entity.Item) error
// }

const (
	OWNER    = "owner"
	SUBOWNER = "subowner"
)

const (
	STATUS_OPENED = "open"
	STATUS_CLOSED = "close"
)

type postgreListStorage struct {
	db *sql.DB
}

func NewListStorage(db *sql.DB) *postgreListStorage {
	return &postgreListStorage{db}
}

func (s *postgreListStorage) Create(ctx context.Context, owner entity.User, list entity.List) (int, error) {
	q1 := "INSERT INTO lists (key, status) VALUES ($1, $2) RETURNING list_id"
	q2 := "INSERT INTO users_lists (user_id, list_id, ownership) VALUES ($1, $2, $3)"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err = tx.Commit(); err != nil {
		return 0, err
	}
}

func (s *postgreListStorage) GetByID(ctx context.Context, id int) (entity.List, error) {
	return entity.List{}, nil
}

func (s *postgreListStorage) GetAll(ctx context.Context, owner entity.User) ([]entity.List, error) {
	return nil, nil
}

func (s *postgreListStorage) UpdateList(ctx context.Context, list entity.List) error {
	return nil
}

func (s *postgreListStorage) UpdateItem(ctx context.Context, item entity.Item) error {
	return nil
}

func (s *postgreListStorage) AddSubowner(ctx context.Context, subowner entity.User, list entity.List) error {
	return nil
}
