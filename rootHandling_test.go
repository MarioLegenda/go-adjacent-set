package golang_nested_set

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestRootCreation(t *testing.T) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		t.Errorf("Cannot establish mysql connectiong with message %s", err)
	}

	metadata := AdjacentSetMetadata{
		Handle:    db,
		TableName: tableName,
	}

	rootId, err := metadata.CreateRoot("Root")

	if err != nil {
		t.Errorf("TestRootCreation() failed in creating the root node with message: %s", err)
	}

	if rootId == 0 {
		t.Errorf("TestRootCreation() AdjacentSetMetadata::CreateRoot() cannot return a 0 (zero)")
	}

	rootCategory, err := metadata.Root()
	if err != nil {
		panic(err)
	}

	if rootCategory.Name != "Root" {
		t.Errorf("TestRootCreation() failed in getting the root with AdjacentSetMetadata::Root(). AdjacentSetMetadata::Name must be 'Root'")
	}
}
