package repository

import (
	"database/sql"
	"errors"
	"github.com/codeedu/codebank/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(transation domain.Transation, card domain.CreditCard) error {
	stmt, err := t.db.Prepare(
		"insert into transactions(id, credit_card_id, amount, status, description, store, created_at) values ($1, $2, $3, $4, $5, $6, $7)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transation.ID,
		transation.CreditCardId,
		transation.Amount,
		transation.Status,
		transation.Description,
		transation.Store,
		transation.CreatedAt,
	)
	if err != nil {
		return err
	}
	if transation.Status == "approved" {
		err = t.updateBalance(card)
		if err != nil {
			return err
		}
		return nil
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) updateBalance(card domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2", card.Balance, card.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into credit_cards(id, name, number, expiration_month,expiration_year, CVV,balance, balance_limit) 
								values($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}
