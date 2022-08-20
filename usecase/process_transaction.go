package usecase

import (
	"github.com/codeedu/codebank/domain"
	"github.com/codeedu/codebank/dto"
	"time"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transation, error) {
	creditCard := u.hydrateCreditCard(transactionDto)
	creditCardBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transation{}, err
	}
	creditCard.ID = creditCardBalanceAndLimit.ID
	creditCard.Limit = creditCardBalanceAndLimit.Limit
	creditCard.Balance = creditCardBalanceAndLimit.Balance

	t := u.newTransaction(transactionDto, creditCardBalanceAndLimit)
	t.ProcessAndValidate(creditCard)
	err = u.TransactionRepository.SaveTransaction(*t, *creditCard)
	if err != nil {
		return domain.Transation{}, err
	}
	return *t, nil
}

func (u UseCaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.CVV = transactionDto.CVV
	creditCard.ExpirationYear = transactionDto.ExpiratrionYear
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	return creditCard
}

func (u UseCaseTransaction) newTransaction(transactionDto dto.Transaction, card domain.CreditCard) *domain.Transation {
	t := domain.NewTransation()
	t.CreditCardId = card.ID
	t.Amount = transactionDto.Amount
	t.Store = transactionDto.Store
	t.Description = transactionDto.Description
	t.CreatedAt = time.Now()
	return t
}
