package golang_nested_set

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"database/sql"
	_  "github.com/go-sql-driver/mysql"
)

var (
	username string
	password string
	tableName string
	dbName string
)

func init() {
	flag.StringVar(&username, "username", "", "Database username")
	flag.StringVar(&password, "password", "", "Database Password")
}

func TestMain(m *testing.M) {
	flag.Parse()

	tableName = "adjacent_set_testing_table"
	dbName = "adjacent_set_testing"

	dsn := fmt.Sprintf("%s:%s@/", username, password)
	createDb := fmt.Sprintf("CREATE DATABASE %s", dbName)
	dropDb := fmt.Sprintf("DROP DATABASE %s", dbName)
	createTable := fmt.Sprintf("CREATE TABLE %s (%s_id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, parent INT DEFAULT NULL)", tableName, tableName)
	useDb := fmt.Sprintf("USE %s", dbName)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	_, err = db.Exec(createDb)

	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(useDb)

	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(createTable)

	if err != nil {
		panic(err.Error())
	}

	code := m.Run()

	_, err = db.Exec(dropDb)

	if err != nil {
		panic(err.Error())
	}

	os.Exit(code)
}

func TestNonExistingRoot(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.Root()

	if err == nil {
		t.Errorf("TestNonExistentRoot() should fail but it did not with message %s", err.Error())
	}
}

func TestCreateCategoryNameFail(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CreateCategory(Category{
		Name:       "",
		CategoryId: 0,
		Parent:     "",
		ParentId:   4,
	})

	if err == nil {
		t.Errorf(fmt.Sprintf("TestCreateCategoryNameFail test should have failed but it did not"))
	}
}

func TestCreateCategoryParentFail(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CreateCategory(Category{
		Name:       "someCategory",
		CategoryId: 0,
		Parent:     "",
		ParentId:   0,
	})

	if err == nil {
		t.Errorf(fmt.Sprintf("TestCreateCategoryParentFail test should have failed but it did not"))
	}
}

func TestCategoryCreate(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CreateCategory(Category{
		Name:       "testingCategory",
		CategoryId: 0,
		Parent:     "",
		ParentId:   4,
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCategoryExistInvalidParametersFail(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CategoryExists(Category{
		Name:       "",
		CategoryId: 0,
		Parent:     "",
		ParentId:   0,
	})

	if err == nil {
		t.Errorf("CategoryExists() should have failed because of invalid parameters")
	}
}

func TestCategoryExists(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	categoryName := "TestCategoryExists"

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CreateCategory(Category{
		Name:       categoryName,
		CategoryId: 0,
		Parent:     "",
		ParentId:   4,
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	ok, err := metadata.CategoryExists(Category{
		Name:       categoryName,
		CategoryId: 0,
		Parent:     "",
		ParentId:   0,
	})

	if ok == false {
		if err != nil {
			t.Errorf("TestCategoryExists() failed with error %s", err.Error())

			return
		}

		t.Errorf("TestCategoryExists() failed")
	}
}

func TestGetCategoryFail(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.GetCategory(FetchOptions{
		Id:           0,
		Name:         "",
	})

	if err == nil {
		t.Errorf("TestGetCategoryFail() should have failed but it didn't")
	}
}

func TestGetNonExistentCategory(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connecting with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	result, err := metadata.GetCategory(FetchOptions{
		Id:   500,
		Name: "",
	})

	if err == nil {
		t.Errorf("TestGetNonExistentCategory() should have returned an error")
	}

	if result.Id != 0 {
		t.Errorf("TestGetNonExistentCategory() should have Id with value 0")
	}

	if result.Name != "" {
		t.Errorf("TestGetNonExistentCategory() should have Name as an empty string")
	}

	if result.parentId != 0 {
		t.Errorf("TestGetNonExistentCategory() should have Parent with value nil")
	}
}

func TestGetCategory(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connecting with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	_, err = metadata.CreateCategory(Category{
		Name:       "TestGetCategory",
		CategoryId: 0,
		Parent:     "",
		ParentId:   1,
	})

	if err != nil {
		t.Errorf(err.Error())
	}

	result, err := metadata.GetCategory(FetchOptions{
		Id:   0,
		Name: "TestGetCategory",
	})

	if err != nil {
		t.Errorf("TestGetCategory failed with error: %s", err)
	}

	if result.Id == 0 {
		t.Errorf("TestGetCategory did not return the correct result id")
	}

	if result.Name == "" {
		t.Errorf("TestGetCategory did not return the correct result name")
	}
}


