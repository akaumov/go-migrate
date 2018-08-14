package lib

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Handler interface {
	BeforeMigration(transaction *sql.Tx, migration *Migration) error
	BeforeAction(transaction *sql.Tx, migration *Migration, index int, method string, params interface{}) error
	AfterAction(transaction *sql.Tx, migration *Migration, index int, method string, params interface{}) error
	AfterMigration(transaction *sql.Tx, migration *Migration) error
}
