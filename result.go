package golang_nested_set

import (
	"database/sql"
	"errors"
	"fmt"
)

func (result *AdjacentSetResult) Parent() (AdjacentSetResult, error) {
	fetchParentSql := fmt.Sprintf("SELECT %s_id, name, COALESCE(parent, 0) AS parent FROM %s WHERE %s_id = ?", *result.tableName, *result.tableName, *result.tableName)
	var (
		id int64
		name string
		parent int64
	)
	db := result.handle

	err := db.QueryRow(fetchParentSql, result.parentId).Scan(&id, &name, &parent)

	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Parent with id %d does not exist", result.parentId)
			return AdjacentSetResult{}, errors.New(msg)
		} else {
			return AdjacentSetResult{}, errors.New(err.Error())
		}
	}

	return AdjacentSetResult{
		Id: id,
		Name: name,
		parentId: parent,
		tableName: result.tableName,
		handle: result.handle,
	}, nil
}

func (result AdjacentSetResult) Children() ([]AdjacentSetResult, error) {
	fetchChildrenSql := fmt.Sprintf("SELECT %s_id, name, COALESCE(parent, 0) AS parent FROM %s WHERE parent = ?", *result.tableName, *result.tableName)

	var (
		id int64
		name string
		parent int64
		resultSet []AdjacentSetResult
	)

	db := result.handle

	rows, err := db.Query(fetchChildrenSql, result.parentId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &parent)
		if err != nil {
			rows.Close()

			return nil, errors.New(err.Error())
		}

		resultSet = append(resultSet, AdjacentSetResult{
			Id:        id,
			Name:      name,
			parentId:  parent,
			tableName: result.tableName,
			handle:    result.handle,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return resultSet, nil
}
