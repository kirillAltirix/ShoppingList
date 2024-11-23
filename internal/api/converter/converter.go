package converter

import (
	"ShoppingList/internal/api"
	"ShoppingList/internal/domain/entity"
)

func EntityListToCreateListResponse(list entity.List) api.CreateListResponse {
	res := api.CreateListResponse{
		Key:    list.Key,
		Status: list.Status,
	}
	for _, item := range list.Items {
		res.Items = append(res.Items, api.ShoppingListItemResponse{
			ID:     item.ID,
			Name:   item.Name,
			Status: item.Status,
		})
	}

	return res
}

func EntityListToGetListResponse(list entity.List) api.GetListResponse {
	res := api.GetListResponse{
		Key:    list.Key,
		Status: list.Status,
	}
	for _, item := range list.Items {
		res.Items = append(res.Items, api.ShoppingListItemResponse{
			ID:     item.ID,
			Name:   item.Name,
			Status: item.Status,
		})
	}

	return res
}

func EntityListsToGetListsRespose(lists []entity.List) api.GetListsRespose {
	res := make(api.GetListsRespose, 0, len(lists))
	for _, list := range lists {
		listRes := api.ShoppingListRespose{
			Key:    list.Key,
			Status: list.Status,
		}
		for _, item := range list.Items {
			listRes.Items = append(listRes.Items, api.ShoppingListItemResponse{
				ID:     item.ID,
				Name:   item.Name,
				Status: item.Status,
			})
		}
		res = append(res, listRes)
	}
	return res
}

func EntityListToUpdateListItemStatusResponse(list entity.List) api.ShoppingListRespose {
	res := api.ShoppingListRespose{
		Key:    list.Key,
		Status: list.Status,
	}
	for _, item := range list.Items {
		res.Items = append(res.Items, api.ShoppingListItemResponse{
			ID:     item.ID,
			Name:   item.Name,
			Status: item.Status,
		})
	}

	return res
}
