package utils

import (
	"github.com/go-sql-driver/mysql"
)

type EnumMySQLError int32

const (
	EnumMySQLError_DUPLICATE_ENTRY EnumMySQLError = 1062
)

func IsMySQLDuplicateError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)

	return ok && mysqlErr != nil &&
		mysqlErr.Number == uint16(EnumMySQLError_DUPLICATE_ENTRY)
}
