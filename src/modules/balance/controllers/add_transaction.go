package balance

import (
	"common/try"
	"modules/balance/acc"
	"modules/misc/trans"
	"time"

	"github.com/gin-gonic/gin"
)

type transactionSinglePayload struct {
	AccountID     int64     `json:"account_id"`
	Amount        acc.Money `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	Description   string    `json:"desciption"`
	Tags          []string  `json:"tags"`
	CorrelationID string    `json:"correlation_id"`
}

type transactionMultiPayload struct {
	Transactions []transactionSinglePayload `json:"transactions"`
}

type transactionSingleResponse struct {
	CorrelationID string `json:"correlation_id"`
	Error         string `json:"error,omitempty"`
	TransactionID int64  `json:"transaction_id,omitempty"`
}

type transactionMultiResponse []transactionSingleResponse

// addTransaction add transaction into system from a user (async)
// @Route {
// 		url = /transactions
//		method = post
//      payload = transactionMultiPayload
//		#resource = resource_name
//      200 = transactionMultiResponse
// }
func (u *Controller) addTransaction(ctx *gin.Context) {
	var err error
	usr := u.MustGetUser(ctx)
	payload := u.MustGetPayload(ctx).(*transactionMultiPayload)

	var resp = make(transactionMultiResponse, len(payload.Transactions))
	var accCache map[int64]*acc.Account
	m := acc.NewAccManager()
	for i, pl := range payload.Transactions {
		acc, ok := accCache[pl.AccountID]
		if !ok {
			acc, err = m.FindAccountByID(pl.AccountID)
			if err != nil {
				resp[i] = transactionSingleResponse{Error: err.Error()}
				continue
			}
			accCache[pl.AccountID] = acc
		}

		if acc.Disabled {
			resp[i] = transactionSingleResponse{
				Error:         trans.T("account_is_disabled"),
				CorrelationID: pl.CorrelationID,
			}
			continue
		}
		// Now we have account and unit. insert transaction
		t, err := m.AddTransaction(
			usr,
			acc,
			pl.Amount,
			pl.CreatedAt,
			pl.Description,
			pl.CorrelationID,
			pl.Tags...,
		)
		if err != nil {
			resp[i] = transactionSingleResponse{
				Error:         try.Try(err).Error(),
				CorrelationID: pl.CorrelationID,
			}
			continue
		}
		resp[i] = transactionSingleResponse{
			TransactionID: t.ID,
			CorrelationID: pl.CorrelationID,
		}
	}

	u.OKResponse(ctx, resp)
}
