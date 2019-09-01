package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func PopulateDb(tableName string, tx *sql.Tx, tree *Node) error {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s(name, parent) VALUES (?, ?)", tableName))
	if err != nil {
		return errors.New(err.Error())
	}

	res, err := stmt.Exec(tree.name, tree.parent)

	if err != nil {
		return errors.New(err.Error())
	}

	defer stmt.Close()

	parentId, err := res.LastInsertId()

	err = recursiveTreePopulate(tx, tree.children, stmt, parentId)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func recursiveTreePopulate(tx *sql.Tx, tree []Node, stmt *sql.Stmt, parentId int64) error {
	for i := 0; i < len(tree); i++ {
		node := tree[i]

		res, err := stmt.Exec(node.name, parentId)

		if err != nil {
			return errors.New(err.Error())
		}

		parentId, err := res.LastInsertId()

		if node.children != nil {
			err = recursiveTreePopulate(tx, node.children, stmt, parentId)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

