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

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByChatID(ctx context.Context, user entity.User) (entity.User, error)
}

type ListService interface {
	CreateList(ctx context.Context, user entity.User, list entity.List) (entity.List, error)
	GetList(ctx context.Context, user entity.User, list entity.List) (entity.List, error)
	GetLists(ctx context.Context, user entity.User) ([]entity.List, error)
	DeleteList(ctx context.Context, user entity.User, list entity.List) error
	CloseList(ctx context.Context, user entity.User, list entity.List) error
	UpdateListItemStatus(ctx context.Context, user entity.User, list entity.List, item entity.Item) (entity.List, error)
}

type httpService struct {
	UserService UserService
	ListService ListService
}

func NewService() *httpService {
	return &httpService{}
}

func (h *httpService) Register(router *mux.Router) {
	router.HandleFunc("/user", h.CreateUser).Methods("POST")
	router.HandleFunc("/user/{chat_id}/list", h.CreateList).Methods("POST")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.GetList).Methods("GET")
	router.HandleFunc("/user/{chat_id}/list", h.GetLists).Methods("GET")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.DeleteList).Methods("DELETE")
	router.HandleFunc("/user/{chat_id}/list/{key}", h.CloseList).Methods("PUT")
	router.HandleFunc("/user/{chat_id}/list/{key}/item/{item_id}", h.UpdateListItemStatus).Methods("PUT")
	router.HandleFunc("/healthcheck", h.Healthcheck).Methods("GET")
}

func (h *httpService) CreateUser(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	req := api.CreateUserRequest{}
	if err := d.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := entity.User{
		ChatID:   req.ChatID,
		Username: req.Username,
	}
	if err := h.UserService.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelList := entity.List{}
	user := entity.User{ChatID: chatID}

	for _, item := range req.Items {
		modelList.Items = append(modelList.Items, entity.Item{
			Name: item,
		})
	}

	user, err := h.UserService.GetUserByChatID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list, err := h.ListService.CreateList(r.Context(), user, modelList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToCreateListResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	list, err := h.ListService.GetList(r.Context(), user, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToGetListResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *httpService) GetLists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"]

	user := entity.User{
		ChatID: chatID,
	}

	user, err := h.UserService.GetUserByChatID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lists, err := h.ListService.GetLists(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := converter.EntityListsToGetListsRespose(lists)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	user, err := h.UserService.GetUserByChatID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ListService.DeleteList(r.Context(), user, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	user, err := h.UserService.GetUserByChatID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ListService.CloseList(r.Context(), user, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	user, err := h.UserService.GetUserByChatID(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list, err = h.ListService.UpdateListItemStatus(r.Context(), user, list, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := converter.EntityListToUpdateListItemStatusResponse(list)

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *httpService) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
