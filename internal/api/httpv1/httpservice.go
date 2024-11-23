package httpv1

import (
	"ShoppingList/internal/api"
	"ShoppingList/internal/api/converter"
	"ShoppingList/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ShoppingList interface {
	CreateUser(ctx context.Context, user entity.User) error
	CreateList(ctx context.Context, user entity.User, list entity.List) (entity.List, error)
	GetList(ctx context.Context, user entity.User, list entity.List) (entity.List, error)
	GetLists(ctx context.Context, user entity.User) ([]entity.List, error)
	DeleteList(ctx context.Context, user entity.User, list entity.List) error
	CloseList(ctx context.Context, user entity.User, list entity.List) error
	UpdateListItemStatus(ctx context.Context, user entity.User, list entity.List, item entity.Item) (entity.List, error)
}

type httpService struct {
	shoppingList ShoppingList
}

func NewService(shoppingList ShoppingList) *httpService {
	return &httpService{shoppingList}
}

func (h *httpService) Register(router *mux.Router) {
	router.HandleFunc("/user", h.CreateUser).Methods("POST")
	router.HandleFunc("/user/{chat_id}/list", h.CreateList).Methods("POST")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.GetList).Methods("GET")
	router.HandleFunc("/user/{chat_id}/list", h.GetLists).Methods("GET")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.DeleteList).Methods("DELETE")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.CloseList).Methods("PUT")
	router.HandleFunc("/user/{chat_id}/list/{key}/items/{item_id}", h.UpdateListItemStatus).Methods("PUT")
}

func (h *httpService) CreateUser(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	req := api.CreateUserRequest{}
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := entity.User{
		ChatID:   req.ChatID,
		Username: req.Username,
	}
	if err := h.shoppingList.CreateUser(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *httpService) CreateList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]

	d := json.NewDecoder(r.Body)
	req := api.CreateListRequest{}
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	modelList := entity.List{}
	user := entity.User{ChatID: chatID}

	for _, item := range req.Items {
		modelList.Items = append(modelList.Items, entity.Item{
			Name: item,
		})
	}

	list, err := h.shoppingList.CreateList(r.Context(), user, modelList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToCreateListResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) GetList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]
	key := vars["key"]

	user := entity.User{
		ChatID: chatID,
	}
	list := entity.List{
		Key: key,
	}

	list, err := h.shoppingList.GetList(r.Context(), user, list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToGetListResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) GetLists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]

	user := entity.User{
		ChatID: chatID,
	}

	lists, err := h.shoppingList.GetLists(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := converter.EntityListsToGetListsRespose(lists)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) DeleteList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]
	key := vars["key"]

	user := entity.User{
		ChatID: chatID,
	}
	list := entity.List{
		Key: key,
	}

	err := h.shoppingList.DeleteList(r.Context(), user, list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) CloseList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]
	key := vars["key"]

	user := entity.User{
		ChatID: chatID,
	}
	list := entity.List{
		Key: key,
	}

	err := h.shoppingList.CloseList(r.Context(), user, list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) UpdateListItemStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]
	key := vars["key"]
	itemIDStr := vars["item_id"]
	itemID, _ := strconv.Atoi(itemIDStr)

	user := entity.User{
		ChatID: chatID,
	}
	list := entity.List{
		Key: key,
	}
	item := entity.Item{
		ID: itemID,
	}

	list, err := h.shoppingList.UpdateListItemStatus(r.Context(), user, list, item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToUpdateListItemStatusResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
