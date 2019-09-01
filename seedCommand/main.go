package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"database/sql"
	_  "github.com/go-sql-driver/mysql"
	"errors"
)

var (
	username string
	password string
	tableName string
	dbName string
	depth int
	leafs int
)

type validationMetadata struct {
	success bool
	messages []string
}

func init() {
	flag.StringVarP(&username, "user", "u", "", "Specify the MySQL username")
	flag.StringVarP(&password, "password", "p", "", "Specify the MySQL password")
	flag.StringVarP(&tableName, "table", "t", "", "Specify the table name")
	flag.StringVarP(&dbName, "db", "d", "", "Specify the database name")
	flag.IntVarP(&depth, "depth", "h", 5, "Specify the database name")
	flag.IntVarP(&leafs, "leafs", "l", 5, "Specify the database name")
}

func validateArgs() validationMetadata {
	messages := []string{}

	if username == "" {
		messages = append(messages, "Mysql username not present. Use --user or -u to specify the username")
	}

	if password == "" {
		messages = append(messages, "Mysql password not present. Use --password or -p to specify the password")
	}

	if tableName == "" {
		messages = append(messages, "Mysql table name not present. Use --table or -t to specify the table name")
	}

	if dbName == "" {
		messages = append(messages, "Mysql database name not present. Use --db or -d to specify the database name")
	}

	if len(messages) > 0 {
		return validationMetadata{
			success:  false,
			messages: messages,
		}
	}

	return validationMetadata{
		success:  true,
		messages: nil,
	}
}

func showErrors(metadata validationMetadata) {
	if !metadata.success {
		fmt.Println("Invalid usage of seed command. For help in using this command, use --help")
		fmt.Println()

		for i := range metadata.messages {
			msg := metadata.messages[i]

			fmt.Println(msg)
		}

		fmt.Println()
		fmt.Println("Example: ./seed -u user -p password -d dbName -t tableName")
		fmt.Println()
	}
}

func getMysqlHandle() (*sql.DB, error){
	dsn := fmt.Sprintf("%s:%s@/%s", username, password, dbName)
	db, err := sql.Open("mysql", dsn)

	// propagating the error to main to get a sensible formatted string that tells us what happened and where

	if err != nil {
		return db, errors.New(fmt.Sprintf("origin: getMysqlHandle(), cause: Unsuccessful handle creation, message: %s;", err))
	}

	err = db.Ping()

	if err != nil {
		return db, errors.New(fmt.Sprintf("origin: getMysqlHandle(), cause: Unsuccessful handle creation, message: %s", err))
	}

	return db, err
}

func runCommand() error {
	db, err := getMysqlHandle()

	if err != nil {
		return errors.New(err.Error())
	}

	if err != nil {
		return errors.New(err.Error())
	}

	tree := CreateTree("Root")
	tree.Populate(depth, leafs)

	tx, err := db.Begin()
	if err != nil {
		return errors.New(err.Error())
	}

	err = PopulateDb(tableName, tx, &tree)
	if err != nil {
		return errors.New(err.Error())
	}

	err = tx.Commit()

	if err != nil {
		return errors.New(err.Error())
	}

	defer tx.Rollback()

	return nil
}

func main() {
	flag.Parse()

	metadata := validateArgs()

	if !metadata.success {
		showErrors(metadata)
	} else {
		err := runCommand()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(fmt.Sprintf("Command finished successfully"))
	}
}
