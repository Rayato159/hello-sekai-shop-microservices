package payment

type (
	ItemServiceReq struct {
		Items []*ItemServiceReqDatum `json:"items" validate:"required"`
	}

	ItemServiceReqDatum struct {
		ItemId string  `json:"item_id" validate:"required,max=64"`
		Price  float64 `json:"price"`
	}

	PaymentTransferReq struct {
		PlayerId string  `json:"player_id"`
		ItemId   string  `json:"item_id"`
		Amout    float64 `json:"amount"`
	}

	PaymentTransferRes struct {
		InventoryId   string  `json:"inventory_id"`
		TransactionId string  `json:"transaction_id"`
		PlayerId      string  `json:"player_id"`
		ItemId        string  `json:"item_id"`
		Amount        float64 `json:"amount"`
		Error         string  `json:"error"`
	}
)
