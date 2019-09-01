package golang_nested_set

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	createCategorySql string
	categoryExistsSql string
	fetchRootSql string
	fetchCategorySql string
	createRootSql string
	db *sql.DB
)

func (metadata AdjacentSetMetadata) CreateRoot(name string) (int64, error) {
	createRootSql = fmt.Sprintf("INSERT INTO %s(name, parent) VALUES(?, NULL)", metadata.TableName)

	db = metadata.Handle

	var lastInsertId int64

	stmt, err := db.Prepare(createRootSql)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	res, err := stmt.Exec(name)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	lastInsertId, err = res.LastInsertId()

	return lastInsertId, nil
}

func (metadata AdjacentSetMetadata) CreateCategory(c Category) (int64, error) {
	if c.Name == "" {
		return 0, errors.New(fmt.Sprintf("Invalid parameters for CategoryCreate(). Category::name must not be an empty string"))
	}

	if c.ParentId == 0 && c.Parent == "" {
		return 0, errors.New(fmt.Sprintf("Invalid parameters for CategoryCreate(). Category::parentId has to be a non 0 id or Category::parent has to be a string"))
	}

	createCategorySql = fmt.Sprintf("INSERT INTO %s(name, parent) VALUES(?, ?)", metadata.TableName)

	db = metadata.Handle

	stmt, err := db.Prepare(createCategorySql)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	res, err := stmt.Exec(c.Name, c.ParentId)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	defer stmt.Close()

	// error is ignored since the statement did not fail so we can expect that insert statement must have a last insert id
	lastInsertId, err := res.LastInsertId()

	return lastInsertId, nil
}

func (metadata AdjacentSetMetadata) CategoryExists(c Category) (bool, error) {
	if c.CategoryId == 0 && c.Name == "" {
		return false, errors.New("Invalid parameters to CategoryExists(). Category::categoryId or Category::name has to be set")
	}

	var categoryId int64
	var categoryName string

	if c.CategoryId != 0 {
		categoryExistsSql = fmt.Sprintf("SELECT %s_id FROM %s WHERE %s_id = ?", metadata.TableName, metadata.TableName, metadata.TableName)
		categoryId = c.CategoryId
	} else {
		categoryExistsSql = fmt.Sprintf("SELECT %s_id FROM %s WHERE name = ?", metadata.TableName, metadata.TableName)
		categoryName = c.Name
	}

	db = metadata.Handle

	var id int64
	var err error

	if categoryId != 0 {
		err = db.QueryRow(categoryExistsSql, categoryId).Scan(&id)
	} else {
		err = db.QueryRow(categoryExistsSql, categoryName).Scan(&id)
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (metadata AdjacentSetMetadata) GetCategory(opts FetchOptions) (AdjacentSetResult, error) {
	if opts.Name == "" && opts.Id == 0 {
		return AdjacentSetResult{}, errors.New("Invalid FetchOptions to GetCategory(). FetchOptions::Name or FetchOptions::Id must be provided")
	}

	db = metadata.Handle

	var (
		err error
		field string
		categoryId int64
		categoryName string
		id int64
		name string
		parentField int64
	)

	if opts.Id != 0 {
		categoryId = opts.Id
		field = fmt.Sprintf("%s_id", metadata.TableName)
	} else {
		field = "name"
		categoryName = opts.Name
	}

	fetchCategorySql = fmt.Sprintf("SELECT %s_id, name, COALESCE(parent, 0) AS parent FROM %s WHERE %s = ?", metadata.TableName, metadata.TableName, field)

	if categoryId != 0 {
		err = db.QueryRow(fetchCategorySql, categoryId).Scan(&id, &name, &parentField)
	} else {
		err = db.QueryRow(fetchCategorySql, categoryName).Scan(&id, &name, &parentField)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return AdjacentSetResult{}, errors.New("Adjacent set result does not exist")
		} else {
			return AdjacentSetResult{}, errors.New(err.Error())
		}
	}

	return AdjacentSetResult{
		Id:     id,
		Name:   name,
		parentId: parentField,
		tableName: &metadata.TableName,
		handle: metadata.Handle,
	}, nil
}

func (metadata AdjacentSetMetadata) Root() (RootCategory, error) {
	fetchRootSql = fmt.Sprintf("SELECT %s_id, name FROM %s WHERE parent IS NULL", metadata.TableName, metadata.TableName)
	var (
		id int64
		name string
	)

	db = metadata.Handle

	err := db.QueryRow(fetchRootSql).Scan(&id, &name)

	if err != nil {
		if err == sql.ErrNoRows {
			return RootCategory{}, errors.New("Root category does not exist")
		} else {
			return RootCategory{}, errors.New(err.Error())
		}
	}

	return RootCategory{
		Name:   name,
		Id:     id,
	}, nil
}

