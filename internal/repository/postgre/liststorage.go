package postgre

import (
	"ShoppingList/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
)

const (
	OWNER    = "owner"
	SUBOWNER = "subowner"
)

const (
	STATUS_OPENED  = "opened"
	STATUS_CLOSED  = "closed"
	STATUS_DELETED = "deleted"
)

type postgreListStorage struct {
	db *sql.DB
}

func NewListStorage(db *sql.DB) *postgreListStorage {
	return &postgreListStorage{db}
}

func (s *postgreListStorage) Create(ctx context.Context, owner entity.User, list entity.List) (entity.List, error) {
	q1 := "INSERT INTO lists (key, status) VALUES ($1, $2) RETURNING list_id"
	q2 := "INSERT INTO users_lists (user_id, list_id, ownership) VALUES ($1, $2, $3)"
	q3 := "INSERT INTO items (name, status) VALUES ($1, $2) RETURNING item_id"
	q4 := "INSERT INTO lists_items (list_id, item_id) VALUES ($1, $2)"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return list, err
	}
	defer tx.Rollback()

	var listID int
	if err = tx.QueryRowContext(ctx, q1, list.Key, list.Status).Scan(&listID); err != nil {
		return list, err
	}

	if _, err = tx.ExecContext(ctx, q2, owner.ID, listID, OWNER); err != nil {
		return list, err
	}

	var itemIDs []int
	for _, item := range list.Items {
		var itemID int
		if err = tx.QueryRowContext(ctx, q3, item.Name, item.Status).Scan(&itemID); err != nil {
			return list, err
		}

		if _, err = tx.ExecContext(ctx, q4, listID, itemID); err != nil {
			return list, err
		}

		itemIDs = append(itemIDs, itemID)
	}

	if len(itemIDs) != len(list.Items) {
		return list, errors.New("items insertion failed")
	}

	if err = tx.Commit(); err != nil {
		return list, err
	}

	list.ID = listID
	for i := range list.Items {
		list.Items[i].ID = itemIDs[i]
	}

	return list, nil
}

func (s *postgreListStorage) GetByID(ctx context.Context, id int) (entity.List, error) {
	q := "SELECT * FROM lists WHERE list_id = $1"
	return entity.List{}, nil
}

func (s *postgreListStorage) GetAll(ctx context.Context) ([]entity.List, error) {
	q1 := "SELECT lists.list_id, lists.key, lists.status, users_lists.user_id, users_lists.ownership FROM lists JOIN users_lists ON users_lists.list_id = lists.list_id"
	q2 := "SELECT items.item_id, items.name, items.status FORM items JOIN lists_items ON items.item_id = lists_items.item_id WHERE lists_items.list_id = $1"
	return nil, nil
}

func (s *postgreListStorage) UpdateList(ctx context.Context, list entity.List) error {
	q := "UPDATE lists SET key = $1, status = $2 WHERE list_id = $3"
	return nil
}

func (s *postgreListStorage) UpdateItem(ctx context.Context, item entity.Item) error {
	q := "UPDATE items SET name = $1, status = $2 WHERE item_id = $3"
	return nil
}

func (s *postgreListStorage) AddSubowner(ctx context.Context, subowner entity.User, list entity.List) error {
	q := "INSERT INTO user_lists (user_id, list_id, ownership) VALUES ($1, $2, $3)"
	return nil
}
