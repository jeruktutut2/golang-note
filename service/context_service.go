package service

import (
	"context"
	"database/sql"
	"fmt"
	"golang-note/exception"
	"golang-note/helper"
	"golang-note/repository"
)

type ContextService interface {
	Timeout(ctx context.Context, requestId string) (str string, err error)
	CheckContext(ctx context.Context, requestId string) (str string, err error)
	CheckContextTx(ctx context.Context, requestId string) (str string, err error)
}

type ContextServiceImplementation struct {
	db                *sql.DB
	ContextRepository repository.ContextRepository
}

func NewContextService(db *sql.DB, contextRepository repository.ContextRepository) ContextService {
	return &ContextServiceImplementation{
		db:                db,
		ContextRepository: contextRepository,
	}
}

func (service *ContextServiceImplementation) Timeout(ctx context.Context, requestId string) (str string, err error) {
	str, err = service.ContextRepository.Timeout(service.db, ctx)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	return
}

func (service *ContextServiceImplementation) CheckContext(ctx context.Context, requestId string) (str string, err error) {
	rowsAffected, err := service.ContextRepository.CreateTable1(service.db, ctx)
	if err != nil {
		fmt.Println("err:", rowsAffected, err)
		return
	}

	str, err = service.ContextRepository.Timeout(service.db, ctx)
	if err != nil {
		fmt.Println("err:", str, err)
		return
	}

	rowsAffected, err = service.ContextRepository.CreateTable2(service.db, ctx)
	if err != nil {
		fmt.Println("err:", rowsAffected, err)
		return
	}

	str, err = service.ContextRepository.Timeout(service.db, ctx)
	if err != nil {
		fmt.Println("err:", str, err)
		return
	}

	rowsAffected, err = service.ContextRepository.CreateTable3(service.db, ctx)
	if err != nil {
		fmt.Println("", rowsAffected, err)
		return
	}
	return
}

func (service *ContextServiceImplementation) CheckContextTx(ctx context.Context, requestId string) (str string, err error) {
	tx, err := service.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return
	}

	rowsAffected, err := service.ContextRepository.CreateTable1Tx(tx, ctx)
	if err != nil {
		fmt.Println("err:", rowsAffected, err)
		return
	}

	str, err = service.ContextRepository.TimeoutTx(tx, ctx)
	if err != nil {
		fmt.Println("err:", str, err)
		return
	}

	rowsAffected, err = service.ContextRepository.CreateTable2Tx(tx, ctx)
	if err != nil {
		fmt.Println("err:", rowsAffected, err)
		return
	}

	str, err = service.ContextRepository.TimeoutTx(tx, ctx)
	if err != nil {
		fmt.Println("err:", str, err)
		return
	}

	rowsAffected, err = service.ContextRepository.CreateTable3Tx(tx, ctx)
	if err != nil {
		fmt.Println("err:", rowsAffected, err)
		return
	}

	return
}
