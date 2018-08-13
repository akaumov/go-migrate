package lib

import (
	"fmt"
	"database/sql"
	"log"
)

type DBParams struct {
	User string
	Password string
	Name string
	Host string
	Port int
}

func Execute(migrations *[]Migration, dbParams DBParams, handler Handler) (string, error) {

	dbConnectionString := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable",
		dbParams.User,
		dbParams.Password,
		dbParams.Name,
		dbParams.Host,
		dbParams.Port)

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return "", fmt.Errorf("can't connect to db: %v", err)
	}
	defer func() { db.Close() }()

	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("can't connect to db: %v", err)
	}

	log.Println("Connected to db")
	transaction, err := db.Begin()
	if err != nil {
		transaction.Rollback()
		return "", fmt.Errorf("can't start transaction: %v", err)
	}

	err = addMigrationsTableIfNotExist(transaction)
	if err != nil {
		transaction.Rollback()
		return "", fmt.Errorf("can't add migration table: %v", err)
	}

	currentMigrationId, err := getCurrentSyncedMigrationId(transaction)
	if err != nil {
		transaction.Rollback()
		return "", fmt.Errorf("can't read current migration state: %v", err)
	}

	_, err = GetCurrentSnapshot()
	if err != nil {
		return "", err
	}

	var lastMigration *Migration
	isCurrentMigrationPassed := currentMigrationId == ""

	for _, migration := range *migrations {

		if migration.Id == currentMigrationId {
			isCurrentMigrationPassed = true
			continue
		}

		if !isCurrentMigrationPassed {
			continue
		}

		err = handler.BeforeMigration(transaction, &migration)
		if err != nil {
			transaction.Rollback()
			return "", fmt.Errorf("can't apply migration %v: %v\n", migration.Id, err)
		}

		err = applyMigrationActions(transaction, migration, handler)
		if err != nil {
			transaction.Rollback()
			return "", fmt.Errorf("can't apply migration %v: %v\n", migration.Id, err)
		}

		err = addMigrationToMigrationsTable(transaction, migration)
		if err != nil {
			transaction.Rollback()
			return "", fmt.Errorf("can't add migration to migrations table %v: %v\n", migration.Id, err)
		}

		err = handler.AfterMigration(transaction, &migration)
		if err != nil {
			transaction.Rollback()
			return "", fmt.Errorf("can't add migration to migrations table %v: %v\n", migration.Id, err)
		}

		lastMigration = &migration
	}

	err = transaction.Commit()
	if err != nil {
		return "", err
	}

	if lastMigration != nil {
		return lastMigration.Id, nil
	}

	return "", fmt.Errorf("no migrations")
}