package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type TransactionRepository interface {
	SaveTransaction(transation Transation, card CreditCard) error
	GetCreditCard(card CreditCard) (CreditCard, error)
	CreateCreditCard(card CreditCard) error
}

type Transation struct {
	ID           string
	Amount       float64
	Status       string
	Description  string
	Store        string
	CreditCardId string
	CreatedAt    time.Time
}

func NewTransation() *Transation {
	t := &Transation{}
	t.ID = uuid.NewV4().String()
	t.CreatedAt = time.Now()
	return t
}

func (t *Transation) ProcessAndValidate(card *CreditCard) {
	if t.Amount+card.Balance > card.Limit {
		t.Status = "rejected"
	} else {
		t.Status = "approved"
		card.Balance = card.Balance + t.Amount
	}
}
