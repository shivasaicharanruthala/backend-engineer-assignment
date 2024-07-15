package data

import (
	"github/shivasaicharanruthala/backend-engineer-assignment/model"
)

type Receipts interface {
	Get(receiptID string) (*model.ReceiptGetResponse, error)
	Insert(receipt *model.Receipt) *model.ReceiptPostResponse
}
