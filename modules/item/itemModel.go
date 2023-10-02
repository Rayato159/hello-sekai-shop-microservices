package item

import "github.com/Rayato159/hello-sekai-shop-tutorial/modules/models"

type (
	CreateItemReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		Damage   int     `json:"damage" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
	}

	ItemShowCase struct {
		ItemId   string  `json:"item_id"`
		Title    string  `json:"title"`
		Price    float64 `json:"price"`
		Damage   int     `json:"damage"`
		ImageUrl string  `json:"image_url"`
	}

	ItemSearchReq struct {
		Title string `query:"title" validate:"max=64"`
		models.PaginateReq
	}

	ItemUpdateReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required"`
	}

	EnableOrDisableItemReq struct {
		UsageStatus bool `json:"usage_status"`
	}
)
