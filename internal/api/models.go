package api

type ShoppingListItemResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type shoppingListRespose struct {
	Key    string                     `json:"key"`
	Status string                     `json:"status"`
	Items  []ShoppingListItemResponse `json:"items"`
}

type CreateUserRequest struct {
	ChatID   string `json:"chat_id"`
	Username string `json:"username"`
}

type CreateUserResponse struct {
}

type CreateListRequest struct {
	Items []string `json:"items"`
}

type CreateListResponse shoppingListRespose

type GetListRequest struct {
}

type GetListResponse shoppingListRespose

type GetListsRequest struct {
}

type GetListsRespose []shoppingListRespose

type DeleteListRequest struct {
}

type DeleteListResponse struct {
}

type CloseListRequest struct {
}

type CloseListResponse struct {
}

type UpdateListItemStatusRequest struct {
}

type UpdateListItemStatusResponse shoppingListRespose
