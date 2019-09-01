package golang_nested_set

import "database/sql"

type AdjacentSetMetadata struct {
	Handle *sql.DB
	TableName string
}

type Category struct {
	Name string
	CategoryId int64
	Parent string
	ParentId int64
}

type AdjacentSetResult struct {
	Id   int64
	Name string
	parentId int64
	tableName *string
	handle *sql.DB
}

type FetchOptions struct {
	Id int64
	Name string
}

type RootCategory struct {
	Name string
	Id int64
}
