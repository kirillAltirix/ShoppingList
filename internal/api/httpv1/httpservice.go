package httpv1

import (
	"ShoppingList/internal/api"
	"ShoppingList/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ShoppingList interface {
	CreateUser(ctx context.Context, user entity.User) error
	CreateList(ctx context.Context, user entity.User, list entity.List) (entity.List, error)
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

	resList, err := h.shoppingList.CreateList(r.Context(), user, modelList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := api.CreateListResponse{
		Key:    resList.Key,
		Status: resList.Status,
	}
	for _, item := range resList.Items {
		res.Items = append(res.Items, api.ShoppingListItemResponse{
			ID:     item.ID,
			Name:   item.Name,
			Status: item.Status,
		})
	}

	e := json.NewEncoder(w)
	if err := e.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpService) GetList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["chat_id"]
	_ = vars["key"]
	w.WriteHeader(http.StatusOK)
}

func (h *httpService) GetLists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["chat_id"]
	w.WriteHeader(http.StatusOK)
}

func (h *httpService) DeleteList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["chat_id"]
	_ = vars["key"]
	w.WriteHeader(http.StatusOK)
}

func (h *httpService) CloseList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["chat_id"]
	_ = vars["key"]
	w.WriteHeader(http.StatusOK)
}

func (h *httpService) UpdateListItemStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["chat_id"]
	_ = vars["key"]
	_ = vars["item_id"]
	w.WriteHeader(http.StatusOK)
}
