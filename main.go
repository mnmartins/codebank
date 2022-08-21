package main

import (
	"database/sql"
	"fmt"
	"github.com/codeedu/codebank/domain"
	"github.com/codeedu/codebank/infrastructure/repository"
	"github.com/codeedu/codebank/usecase"
	_ "github.com/lib/pq"
	"log"
	"moul.io/banner"
)

func main() {
	fmt.Println(banner.Inline("codebank"))
	db := setupDb()
	defer db.Close()

	//TESTE
	creditCard := domain.NewCreditCard()
	creditCard.Number = "1234"
	creditCard.Name = "Nikita"
	creditCard.ExpirationYear = 2027
	creditCard.ExpirationMonth = 12
	creditCard.CVV = 123
	creditCard.Limit = 1000
	creditCard.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*creditCard)
	if err != nil {
		fmt.Println(err)
	}
	//TESTE
}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	usecase := usecase.NewUseCaseTransaction(transactionRepository)
	return usecase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	return db
}
