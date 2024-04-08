package exception

import (
	"context"
)

// var (
// 	ErrRowsAffectedNotOne      = errors.New("rows affected not one")
// 	ErrRowsAffectedDoesntMatch = errors.New("rows affected doesnt match")
// 	// cycllic redudancy
// 	// ErrConvertDecimalToFloat   = errors.New("cannot convert from decimal to float64")
// )

func CheckError(err error) error {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return NewTimeoutCancelException()
	}
	return NewInternalServerErrorException()

}
