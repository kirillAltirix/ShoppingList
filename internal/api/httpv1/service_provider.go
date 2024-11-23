package httpv1

import (
	"ShoppingList/internal/domain/service"
	"ShoppingList/internal/repository/postgre"
	"database/sql"
)

func Build(db *sql.DB) *httpService {
	listStorage := postgre.NewListStorage(db)
	userStorage := postgre.NewUserStorage(db)
	listService := service.NewListService(listStorage)
	userService := service.NewUserService(userStorage)
	s := NewService()
	s.ListService = listService
	s.UserService = userService
	return s
}
