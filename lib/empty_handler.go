package lib

import "database/sql"

type EmptyHandler struct {
}

func (*EmptyHandler) BeforeMigration(transaction *sql.Tx, migration *Migration) error {
	return nil
}

func (*EmptyHandler) BeforeAction(transaction *sql.Tx, migration *Migration, method string, params interface{}) error {
	return nil
}

func (*EmptyHandler) AfterAction(transaction *sql.Tx, migration *Migration, method string, params interface{}) error {
	return nil
}

func (*EmptyHandler) AfterMigration(transaction *sql.Tx, migration *Migration) error {
	return nil
}

var _ Handler = (*EmptyHandler)(nil)
