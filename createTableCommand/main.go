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
		fmt.Println("Invalid usage of createTable command. For help in using this command, use --help")
		fmt.Println()

		for i := range metadata.messages {
			msg := metadata.messages[i]

			fmt.Println(msg)
		}

		fmt.Println()
		fmt.Println("Example: ./createTable -u user -p password -d dbName -t tableName")
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

	createTableSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s_id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, parent INT DEFAULT NULL, CONSTRAINT unique_name UNIQUE (name))", tableName, tableName)
	createIndexSql := fmt.Sprintf("CREATE INDEX idx_name_and_parent ON %s (name, parent)", tableName)

	_, err = db.Exec(createTableSql)

	if err != nil {
		return errors.New(fmt.Sprintf("origin: runCommand(), cause: Failed query to create the %s table, message: %s", tableName, err))
	}

	_, err = db.Exec(createIndexSql)

	if err != nil {
		return errors.New(fmt.Sprintf("origin: runCommand(), cause: Failed query to create the %s table, message: %s", tableName, err))
	}

	defer db.Close()

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

		fmt.Println(fmt.Sprintf("Command finished successfully. If the table %s does not exists, it has been created", tableName))
	}
}


