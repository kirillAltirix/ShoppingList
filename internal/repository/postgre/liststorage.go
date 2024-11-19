package postgre

import (
	"ShoppingList/internal/domain/entity"
	repository "ShoppingList/internal/repository/model"
	"context"
	"database/sql"
	"errors"
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
	q2 := "INSERT INTO lists_owners (list_id, user_id) VALUES ($1, $2)"
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

	if _, err = tx.ExecContext(ctx, q2, owner.ID, listID); err != nil {
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
	q := "SELECT lists.list_id, lists.key, lists.status, lists_owners.user_id FROM lists JOIN lists_owners ON lists.list_id = lists_owners.list_id WHERE lists.list_id = $1"

	list := repository.List{}
	var ownerID int
	err := s.db.QueryRowContext(ctx, q, id).Scan(&list.ListID, &list.Key, &list.Status, &ownerID)
	if err != nil {
		return entity.List{}, err
	}

	items, err := s.getItemsByListID(ctx, list.ListID)
	if err != nil {
		return entity.List{}, err
	}

	return entity.List{
		ID:      list.ListID,
		Key:     list.Key,
		Status:  list.Status,
		OwnerID: ownerID,
		Items:   items,
	}, nil
}

func (s *postgreListStorage) GetByKey(ctx context.Context, key string) (entity.List, error) {
	q := "SELECT lists.list_id, lists.key, lists.status, lists_owners.user_id FROM lists JOIN lists_owners ON lists.list_id = lists_owners.list_id WHERE lists.key = $1"

	list := repository.List{}
	var ownerID int
	err := s.db.QueryRowContext(ctx, q, key).Scan(&list.ListID, &list.Key, &list.Status, &ownerID)
	if err != nil {
		return entity.List{}, err
	}

	items, err := s.getItemsByListID(ctx, list.ListID)
	if err != nil {
		return entity.List{}, err
	}

	return entity.List{
		ID:      list.ListID,
		Key:     list.Key,
		Status:  list.Status,
		OwnerID: ownerID,
		Items:   items,
	}, nil
}

func (s *postgreListStorage) GetByUserID(ctx context.Context, id int) ([]entity.List, error) {
	q1 := "SELECT list_id FROM lists JOIN lists_owners ON lists.list_id = lists_owners.list_id WHERE lists_owners.user_id = $1"
	q2 := "SELECT list_id FROM lists JOIN lists_subowners ON lists.list_id = lists_subowners.list_id WHERE lists_subowners.user_id = $1"

	listIDs := []int{}
	rows, err := s.db.QueryContext(ctx, q1)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var listID int
		if err := rows.Scan(&listID); err != nil {
			return nil, nil
		}
		listIDs = append(listIDs, listID)
	}

	rows2, err := s.db.QueryContext(ctx, q2)
	if err != nil {
		return nil, nil
	}
	defer rows2.Close()

	for rows2.Next() {
		var listID int
		if err := rows2.Scan(&listID); err != nil {
			return nil, nil
		}
		listIDs = append(listIDs, listID)
	}

	modelList := []entity.List{}

	for _, listID := range listIDs {
		if list, err := s.GetByID(ctx, listID); err != nil {
			return nil, nil
		} else {
			modelList = append(modelList, list)
		}
	}

	return modelList, nil
}

func (s *postgreListStorage) UpdateList(ctx context.Context, list entity.List) error {
	q := "UPDATE lists SET key = $1, status = $2 WHERE list_id = $3"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, q, list.Key, list.Status)
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		if err := s.UpdateItem(ctx, item); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *postgreListStorage) UpdateItem(ctx context.Context, item entity.Item) error {
	q := "UPDATE items SET name = $1, status = $2 WHERE item_id = $3"

	_, err := s.db.ExecContext(ctx, q, item.Name, item.Status)

	return err
}

func (s *postgreListStorage) AddSubowner(ctx context.Context, subowner entity.User, list entity.List) error {
	q := "INSERT INTO lists_subowners (list_id, user_id) VALUES ($1, $2)"

	_, err := s.db.ExecContext(ctx, q, list.ID, subowner.ID)

	return err
}

func (s *postgreListStorage) getItemsByListID(ctx context.Context, listID int) ([]entity.Item, error) {
	q := "SELECT items.item_id, items.name, items.status FORM items JOIN lists_items ON items.item_id = lists_items.item_id WHERE lists_items.list_id = $1"

	rows, err := s.db.QueryContext(ctx, q, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []repository.Item{}
	for rows.Next() {
		item := repository.Item{}
		if err := rows.Scan(&item.ItemID, &item.Name, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	modelItems := []entity.Item{}
	for _, item := range items {
		modelItems = append(modelItems, entity.Item{
			ID:     item.ItemID,
			Name:   item.Name,
			Status: item.Status,
		})
	}

	return modelItems, nil
}
