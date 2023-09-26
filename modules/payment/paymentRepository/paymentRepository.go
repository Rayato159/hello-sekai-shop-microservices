package paymentRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PaymentRepositoryService interface{}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
	return &paymentRepository{db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}
