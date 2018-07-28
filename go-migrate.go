package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/akaumov/go-migrate/db"
	"path/filepath"

	"github.com/ghodss/yaml"
	"errors"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "migration",
			Usage: "manage migrations",
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "add migrationDescription",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "o",
							Usage: "--o yaml|json",
							Value: "yaml",
						},
					},
					Action: addMigration,
				},
				{
					Name:   "list",
					Usage:  "return migrations",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "o",
							Usage: "--o yaml|json",
							Value: "yaml",
						},
					},
					Action: listMigrations,
				},
				{
					Name:   "snapshot",
					Usage:  "return snapshot",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "o",
							Usage: "--o yaml|json",
							Value: "yaml",
						},
					},
					Action: migrationSnapshot,
				},
				{
					Name:  "table",
					Usage: "operations with tables",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "add tableName",
							Action: addTable,
						},
						{
							Name:   "delete",
							Usage:  "delete tableName",
							Action: deleteTable,
						},
					},
				},
				{
					Name:  "column",
					Usage: "operations with columns of tables",
					Subcommands: []cli.Command{
						{
							Name:  "add",
							Usage: "add tableName columName columnType",
							Flags: []cli.Flag{
								cli.BoolTFlag{
									Name:  "nullable",
									Usage: "isNullable flag, default true",
								},
								cli.StringFlag{
									Name:  "default",
									Usage: "default value",
								},
							},
							Action: addColumn,
						},
						{
							Name:   "delete",
							Usage:  "delete tableName columName",
							Action: deleteColumn,
						},
					},
				},

				{
					Name:  "primary",
					Usage: "operations with primary keys",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "add tableName columnName",
							Action: addPrimaryKey,
						},
						{
							Name:   "delete",
							Usage:  "delete tableName columnName",
							Action: deletePrimaryKey,
						},
					},
				},
				{
					Name:  "sync",
					Usage: "sync migrations",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "user",
							EnvVar: "GO_MIGRATE_USER",
							Usage:  "user name",
						},
						cli.StringFlag{
							Name:   "password",
							EnvVar: "GO_MIGRATE_PASSWORD",
							Usage:  "password",
						},
						cli.StringFlag{
							Name:   "db",
							EnvVar: "GO_MIGRATE_DB",
							Usage:  "db",
						},
						cli.StringFlag{
							Name:   "host",
							EnvVar: "GO_MIGRATE_HOST",
							Usage:  "host",
						},
						cli.IntFlag{
							Name:   "port",
							EnvVar: "GO_MIGRATE_PORT",
							Usage:  "port",
						},
						cli.StringFlag{
							Name:   "migrations-dir",
							EnvVar: "GO_MIGRATE_MIGRATIONS_DIR",
							Usage:  "migrations-dir",
						},
					},
					Action: syncMigrations,
				},
				{
					Name:  "relation",
					Usage: "define table relations",
					Subcommands: []cli.Command{
						{
							Name:      "add",
							ArgsUsage: "relation add relationName relationType tableName remoteTableName 'columnName1:remoteColumnName1;columnName2:remoteColumnName2'",
							Action:    addRelation,
						},
						{
							Name:      "delete",
							ArgsUsage: "relation delete table relationName",
							Action:    deleteRelation,
						},
					},
				},
				{
					Name:  "unique",
					Usage: "define unique constraints",
					Subcommands: []cli.Command{
						{
							Name:      "add",
							ArgsUsage: "unique add constraintName tableName 'columnName1;columnName2'",
							Action:    addUniqueConstraint,
						},
						{
							Name:      "delete",
							ArgsUsage: "unique delete table constraintName",
							Action:    deleteUniqueConstraint,
						},
					},
				},
			},
		},
		{
			Name:  "ping",
			Usage: "ping db",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "user",
					EnvVar: "GO_MIGRATE_USER",
					Usage:  "user name",
				},
				cli.StringFlag{
					Name:   "password",
					EnvVar: "GO_MIGRATE_PASSWORD",
					Usage:  "password",
				},
				cli.StringFlag{
					Name:   "db",
					EnvVar: "GO_MIGRATE_DB",
					Usage:  "db",
				},
				cli.StringFlag{
					Name:   "host",
					EnvVar: "GO_MIGRATE_HOST",
					Usage:  "host",
				},
				cli.IntFlag{
					Name:   "port",
					EnvVar: "GO_MIGRATE_PORT",
					Usage:  "port",
				},
			},
			Action: pingDb,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func addMigration(c *cli.Context) error {
	args := c.Args()
	description := args.Get(0)

	outputFormat := c.String("o")

	switch outputFormat {
	case "json":
	case "yaml":
	default:
		return errors.New("wrong output format")
	}

	migrationFileName, err := db.AddMigration(description, db.Format(outputFormat))
	if err == nil {
		fmt.Println(migrationFileName)
	}

	return err
}

func addTable(c *cli.Context) error {
	args := c.Args()
	tableName := args.Get(0)

	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	updatedMigrationId, err := db.AddTable(tableName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func deleteTable(c *cli.Context) error {
	args := c.Args()
	tableName := args.Get(0)

	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	updatedMigrationId, err := db.DeleteTable(tableName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func addColumn(c *cli.Context) error {
	args := c.Args()

	tableName := args.Get(0)
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	columnName := args.Get(1)
	if columnName == "" {
		return fmt.Errorf("column name is required")
	}

	columnType := args.Get(2)
	if columnType == "" {
		return fmt.Errorf("column type is required")
	}

	isNullable := c.BoolT("nullable")
	defaultValue := c.String("default")

	updatedMigrationId, err := db.AddColumn(tableName, columnName, columnType, isNullable, defaultValue)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func deleteColumn(c *cli.Context) error {
	args := c.Args()

	tableName := args.Get(0)
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	columnName := args.Get(1)
	if columnName == "" {
		return fmt.Errorf("column name is required")
	}

	updatedMigrationId, err := db.DeleteColumn(tableName, columnName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func addPrimaryKey(c *cli.Context) error {
	args := c.Args()

	tableName := args.Get(0)
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	columnName := args.Get(1)
	if columnName == "" {
		return fmt.Errorf("column name is required")
	}

	updatedMigrationId, err := db.AddPrimaryKey(tableName, columnName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func deletePrimaryKey(c *cli.Context) error {
	args := c.Args()

	tableName := args.Get(0)
	if tableName == "" {
		return fmt.Errorf("table name is required")
	}

	columnName := args.Get(1)
	if columnName == "" {
		return fmt.Errorf("column name is required")
	}

	updatedMigrationId, err := db.DeletePrimaryKey(tableName, columnName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func listMigrations(c *cli.Context) error {
	migrations, err := db.GetList()
	if err != nil {
		return err
	}

	var packedMigrations []byte

	switch c.String("o") {
	case "yaml":
		packedMigrations, _ = yaml.Marshal(migrations)
	case "json":
		packedMigrations, _ = json.MarshalIndent(migrations, "", "  ")
	default:
		return errors.New("wrong output format")
	}

	fmt.Println(string(packedMigrations))
	return nil
}

func parseColumnsMapping(mappingRaw string) (*[]db.ColumnsMap, error) {
	mapping := []db.ColumnsMap{}

	if mappingRaw != "" {
		for _, rawMap := range strings.Split(mappingRaw, ";") {
			splittedMap := strings.Split(rawMap, ":")

			if len(splittedMap) != 2 {
				return nil, fmt.Errorf("wrong columns mapping: %v\n", rawMap)
			}

			column := splittedMap[0]
			remoteColumn := splittedMap[1]

			mapping = append(mapping, db.ColumnsMap{
				Column:       column,
				RemoteColumn: remoteColumn,
			})
		}
	}

	return &mapping, nil
}

func addRelation(c *cli.Context) error {
	args := c.Args()

	relationName := args.Get(0)
	relationType := args.Get(1)
	table := args.Get(2)
	remoteTable := args.Get(3)
	rawMapping := args.Get(4)

	columnsMapping, err := parseColumnsMapping(rawMapping)
	if err != nil {
		return err
	}

	updatedMigrationId, err := db.AddRelation(relationName, db.RelationType(relationType), table, remoteTable, *columnsMapping)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func deleteRelation(c *cli.Context) error {
	args := c.Args()

	table := args.Get(0)
	relationName := args.Get(1)

	updatedMigrationId, err := db.DeleteRelation(table, relationName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func addUniqueConstraint(c *cli.Context) error {
	args := c.Args()

	constraintName := args.Get(0)
	table := args.Get(1)
	rawColumns := args.Get(2)

	columns := strings.Split(rawColumns, ";")

	updatedMigrationId, err := db.AddUniqueConstraint(constraintName, table, columns)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func deleteUniqueConstraint(c *cli.Context) error {
	args := c.Args()

	table := args.Get(0)
	relationName := args.Get(1)

	updatedMigrationId, err := db.DeleteUniqueConstraint(table, relationName)
	if err != nil {
		return err
	}

	fmt.Println(updatedMigrationId)
	return nil
}

func migrationSnapshot(c *cli.Context) error {
	snapshot, err := db.GetCurrentSnapshot()
	if err != nil {
		return err
	}

	var textSnapshot []byte

	switch c.String("o") {
	case "yaml":
		textSnapshot, _ = yaml.Marshal(*snapshot)
	case "json":
		textSnapshot, _ = json.MarshalIndent(*snapshot, "", "  ")
	default:
		return errors.New("wrong output format")
	}

	fmt.Println(string(textSnapshot))
	return nil
}

func syncMigrations(c *cli.Context) error {

	userName := c.String("user")
	if userName == "" {
		return fmt.Errorf("user name is required")
	}

	password := c.String("password")
	if password == "" {
		return fmt.Errorf("password name is required")
	}

	dbName := c.String("db")
	if dbName == "" {
		return fmt.Errorf("db name is required")
	}

	host := c.String("host")
	if host == "" {
		return fmt.Errorf("host is required")
	}

	port := c.Int("port")
	if port == 0 {
		return fmt.Errorf("port is required")
	}

	migrationsDir := c.String("migrations-dir")
	if migrationsDir == "" {
		migrationsDir = os.Getenv("GO_MIGRATE_MIGRATIONS_DIR")
	}

	if migrationsDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("can't get current dir: %v", err)
		}

		migrationsDir = filepath.Join(currentDir, "migrations")
	}

	lastMigrationId, err := db.Sync(migrationsDir, userName, password, dbName, host, port)
	if err != nil {
		return err
	}

	fmt.Println("DB sync success: " + lastMigrationId)
	return nil
}

func pingDb(c *cli.Context) error {

	userName := c.String("user")
	if userName == "" {
		return fmt.Errorf("user name is required")
	}

	password := c.String("password")
	if password == "" {
		return fmt.Errorf("password name is required")
	}

	dbName := c.String("db")
	if dbName == "" {
		return fmt.Errorf("db name is required")
	}

	host := c.String("host")
	if host == "" {
		return fmt.Errorf("host is required")
	}

	port := c.Int("port")
	if port == 0 {
		return fmt.Errorf("port is required")
	}

	return db.PingDb(userName, password, dbName, host, port)
}
