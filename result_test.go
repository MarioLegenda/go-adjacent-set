package golang_nested_set

import (
	"database/sql"
	"fmt"
	"testing"
	"github.com/bxcodec/faker"
)

type domainName struct {
	DomainName string `faker:"domain_name"`
}

func TestGettingParent(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	createRoot("Root")
	createTestData(metadata)
	createTestData(metadata)
	createTestData(metadata)

	result, err := metadata.GetCategory(FetchOptions{
		Id:   20,
		Name: "",
	})

	if err != nil {
		t.Errorf("TestGettingParent() failed with GetCategory with message: %s", err)
	}

	if result.Name == "" {
		t.Errorf("TestGettingParent() called GetCategory that returned an invalid Name field with empty string")
	}

	if result.Id != 20 {
		t.Errorf("TestGettingParent() called GetCategory that returned an invalid Id field")
	}

	if result.parentId == 0 {
		t.Errorf("TestGettingParent() called GetCategory that returned an invalid parentId field")
	}

	parent, err := result.Parent()

	if err != nil {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent() method with message: %s", err)
	}

	if parent.Id == 0 {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent(). AdjacentSetResult::Id is 0 (zero)")
	}

	if parent.Name == "" {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent(). AdjacentSetResult::Name is an empty string")
	}

	if parent.parentId == 0 {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent(). AdjacentSetResult::parentId cannot be 0 (zero)")
	}

	parentsParent, err := parent.Parent()

	if err != nil {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent() in fetching the parents parent with message: %s", err)
	}

	if parentsParent.Id == 0 {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent(). AdjacentSetResult::Id from parents parent is 0 (zero)")
	}

	if parentsParent.Name == "" {
		t.Errorf("TestGettingParent() failed with AdjacentSetResult::Parent(). AdjacentSetResult::Name from parents parent is an empty string")
	}
}

func TestGettingChildren(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	result, err := metadata.GetCategory(FetchOptions{
		Id:   15,
		Name: "",
	})

	if err != nil {
		t.Errorf("TestGettingChildren() failed with GetCategory with message: %s", err)
	}

	if result.Name == "" {
		t.Errorf("TestGettingChildren() called GetCategory that returned an invalid Name field with empty string")
	}

	if result.parentId == 0 {
		t.Errorf("TestGettingChildren() called GetCategory that returned an invalid parentId field")
	}

	children, err := result.Children()

	if len(children) == 0 {
		t.Errorf("TestingGettingChildren() called AdjacentResultSet::Children() but it did not return the correct number of children")
	}
}

func TestAdjacentResultSetFull(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:            db,
		TableName:         tableName,
	}

	var searchId int64 = 15

	result, err := metadata.GetCategory(FetchOptions{
		Id:   searchId,
		Name: "",
	})

	if err != nil {
		panic(err)
	}

	parent, err := result.Parent()
	if err != nil {
		panic(err)
	}

	if parent.Id == 1 || parent.Id == 0 {
		t.Errorf("TestAdjacentResultSetFull() returned an invalid parent")
	}

	children, err := parent.Children()
	if err != nil {
		panic(err)
	}

	if len(children) == 0 {
		t.Errorf("TestAdjacentResultSetFull() returned an invalid number of children")
	}

	originalFound := false

	for idx := range children {
		original := children[idx]

		if original.Id == original.Id {
			originalFound = true

			break
		}
	}

	if !originalFound {
		t.Errorf("TestAdjacentResultFull was not able to do a complete voyage original->parent->children->original")
	}
}

func createRoot(name string) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		panic(err)
	}

	metadata := AdjacentSetMetadata{
		Handle:    db,
		TableName: tableName,
	}

	_, err = metadata.CreateRoot("Root")

	if err != nil {
		panic(err)
	}
}

func createTestData(metadata AdjacentSetMetadata) {
	fakeName := domainName{}
	faker.FakeData(&fakeName)

	var (
		categoryId int64
		childCategoryId int64
	)

	categoryId, err := metadata.CreateCategory(Category{
		Name:       fakeName.DomainName,
		CategoryId: 0,
		Parent:     "Root",
		ParentId:   0,
	})

	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		fakeName := domainName{}
		faker.FakeData(&fakeName)

		childCategoryId, err = metadata.CreateCategory(Category{
			Name:       fakeName.DomainName,
			CategoryId: 0,
			Parent:     "",
			ParentId:   categoryId,
		})

		if err != nil {
			panic(err)
		}

		for a := 0; a < 5; a++ {
			fakeName := domainName{}
			faker.FakeData(&fakeName)

			_, err = metadata.CreateCategory(Category{
				Name:       fakeName.DomainName,
				CategoryId: 0,
				Parent:     "",
				ParentId:   childCategoryId,
			})

			if err != nil {
				panic(err)
			}
		}
	}
}
