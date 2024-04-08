package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"golang-note/exception"
	"golang-note/helper"
	modelentity "golang-note/model/entity"
	modelrequest "golang-note/model/request"
	modelresponse "golang-note/model/response"
	"golang-note/repository"
	"golang-note/util"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type TransactionService interface {
	Purchase(ctx context.Context, requestId string, userId int32, now int64, purchaseTransactionRequest modelrequest.PurchaseTransactionRequest) (purchaseTransactionResponse modelresponse.PurchaseTransactionResponse, err error)
}

type TransactionServiceImplementation struct {
	MysqlUtil                   util.MysqlUtil
	Validate                    *validator.Validate
	UserRepository              repository.UserRepository
	BookRepository              repository.BookRepository
	WalletRepository            repository.WalletRepository
	TransactionRepository       repository.TransactionRepository
	TransactionDetailRepository repository.TransactionDetailRepository
}

func NewTransactionService(mysqlUtil util.MysqlUtil, validate *validator.Validate, userRepository repository.UserRepository, bookRepository repository.BookRepository, walletRepository repository.WalletRepository, transactionRepository repository.TransactionRepository, transactionDetailRepository repository.TransactionDetailRepository) TransactionService {
	return &TransactionServiceImplementation{
		MysqlUtil:                   mysqlUtil,
		Validate:                    validate,
		UserRepository:              userRepository,
		BookRepository:              bookRepository,
		WalletRepository:            walletRepository,
		TransactionRepository:       transactionRepository,
		TransactionDetailRepository: transactionDetailRepository,
	}
}

func (service *TransactionServiceImplementation) Purchase(ctx context.Context, requestId string, userId int32, now int64, purchaseTransactionRequest modelrequest.PurchaseTransactionRequest) (purchaseTransactionResponse modelresponse.PurchaseTransactionResponse, err error) {
	err = service.Validate.Struct(purchaseTransactionRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, purchaseTransactionRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}

	if len(purchaseTransactionRequest.BookPurchaseTransactionRequests) < 1 {
		err = exception.NewBadRequestException("didn't purchase anything")
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	tx, err := service.MysqlUtil.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	defer func() {
		errCommitOrRollback := service.MysqlUtil.CommitOrRollback(tx, err)
		if errCommitOrRollback != nil {
			helper.PrintLogToTerminal(err, requestId)
			purchaseTransactionResponse = modelresponse.PurchaseTransactionResponse{}
			err = exception.NewInternalServerErrorException()
		}
	}()

	var user modelentity.User
	user.Id = sql.NullInt32{Valid: true, Int32: userId}
	err = service.UserRepository.GetByIdForUpdate(tx, ctx, &user)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	} else if err == sql.ErrNoRows {
		err = exception.NewNotFoundException(fmt.Sprintf("cannot find user with id: %d", userId))
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	var bookIds []int32
	for _, bookPurchaseTransactionRequest := range purchaseTransactionRequest.BookPurchaseTransactionRequests {
		// isContain := helper.CheckArrayInt32(bookIds, bookPurchaseTransactionRequest.BookId)
		isContain := slices.Contains(bookIds, bookPurchaseTransactionRequest.BookId)
		if !isContain {
			bookIds = append(bookIds, bookPurchaseTransactionRequest.BookId)
		}
	}

	var books []modelentity.Book
	err = service.BookRepository.GetBookByInIdsForUpdate(tx, ctx, bookIds, &books)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if len(books) != len(bookIds) {
		err = exception.NewBadRequestException("books data didn't match")
		helper.PrintLogToTerminal(err, requestId)
		return
	}
	var total decimal.Decimal
	fmt.Println("serv:", books)
	for i := 0; i < len(books); i++ {
		fmt.Println(1)
		if books[i].Stock.Int16 < purchaseTransactionRequest.BookPurchaseTransactionRequests[i].Quantity {
			fmt.Println(2)
			err = exception.NewBadRequestException(fmt.Sprintf("insufficient book stock: book stock: %d , quantity request: %d ", books[i].Stock.Int16, purchaseTransactionRequest.BookPurchaseTransactionRequests[i].Quantity))
			helper.PrintLogToTerminal(err, requestId)
			return
		}
		total = total.Add(books[i].Price.Decimal.Mul(decimal.NewFromInt(int64(purchaseTransactionRequest.BookPurchaseTransactionRequests[i].Quantity))))
	}
	fmt.Println(3)

	var wallet modelentity.Wallet
	wallet.UserId = sql.NullInt32{Valid: true, Int32: userId}
	err = service.WalletRepository.GetByUserIdForUpdate(tx, ctx, &wallet)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	} else if err == sql.ErrNoRows {
		err = exception.NewNotFoundException(fmt.Sprintf("cannot find user wallet with user id: %d", userId))
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	if wallet.Balance.Decimal.LessThan(total) {
		err = exception.NewBadRequestException("wallet balance less than total")
		helper.PrintLogToTerminal(err, requestId)
		return
	}
	transaction := &modelentity.Transaction{
		UserId:        sql.NullInt32{Valid: true, Int32: userId},
		Username:      user.Username,
		UserEmail:     user.Email,
		WalletId:      wallet.Id,
		WalletUserId:  wallet.UserId,
		WalletBalance: wallet.Balance,
		Paid:          sql.NullInt16{Valid: true, Int16: 0},
		CreatedAt:     sql.NullInt64{Valid: true, Int64: now},
	}
	rowsAffected, transactionId, err := service.TransactionRepository.Create(tx, ctx, transaction)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected is not one: " + strconv.Itoa(int(rowsAffected)))
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}
	if transactionId < 1 {
		err = errors.New("last insert id less than 1")
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}

	// transactionDetail
	var transactionDetails []*modelentity.TransactionDetail
	for _, bookPurchaseTransactionRequest := range purchaseTransactionRequest.BookPurchaseTransactionRequests {
		var transactionDetail modelentity.TransactionDetail
		transactionDetail.TransactionId = sql.NullInt32{Valid: true, Int32: int32(transactionId)}
	booksLoop:
		for _, book := range books {
			if bookPurchaseTransactionRequest.BookId == book.Id.Int32 {
				transactionDetail.BookId = book.Id
				transactionDetail.BookName = book.Name
				transactionDetail.BookPrice = book.Price
				transactionDetail.Quantity = sql.NullInt16{Valid: true, Int16: bookPurchaseTransactionRequest.Quantity}
				transactionDetail.CreatedAt = sql.NullInt64{Valid: true, Int64: now}
				transactionDetails = append(transactionDetails, &transactionDetail)
				break booksLoop
			}
		}
	}

	rowsAffected, err = service.TransactionDetailRepository.CreateMany(tx, ctx, &transactionDetails)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if len(transactionDetails) != int(rowsAffected) {
		err = errors.New("rows affected transaction detail not equal to number of transaction detail data")
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}

	wallet.Balance = decimal.NullDecimal{Valid: true, Decimal: wallet.Balance.Decimal.Sub(total)}
	rowsAffected, err = service.WalletRepository.UpdateUserWalletBalance(tx, ctx, &wallet)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected is not one: " + strconv.Itoa(int(rowsAffected)))
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}

	for i := 0; i < len(books); i++ {
	bookPurchaseTransactionRequestLoop:
		for _, bookPurchaseTransactionRequest := range purchaseTransactionRequest.BookPurchaseTransactionRequests {
			if books[i].Id.Int32 == bookPurchaseTransactionRequest.BookId {
				books[i].Stock.Int16 -= bookPurchaseTransactionRequest.Quantity
				break bookPurchaseTransactionRequestLoop
			}
		}
	}
	rowsAffected, err = service.BookRepository.UpdateManyStock(tx, ctx, &books)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if int(rowsAffected) != len(books) {
		err = errors.New("rows affected id not one: " + strconv.Itoa(int(rowsAffected)))
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}
	purchaseTransactionResponse = modelresponse.ToPurchaseTransactionResponse(int32(transactionId))
	return
}
