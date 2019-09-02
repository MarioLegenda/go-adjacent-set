1. Overview
2. How to install it
3. Important to know
3. Using built in commands
4. How to use it
5. API

## 1. Overview

This package is a basic adjacency set for MySQL written in Go for representing hierarchical data in MySQL. It is extremely basic and at 
the time of writing these docs, even inefficient. Study the docs carefully if you which to use it.

## 2. Installation

`go get -u github.com/MarioLegenda/go-adjacent-set`

This will install it as a package into your `pkg` directory. After that, import it into your project

````go
import(
	gas "github.com/MarioLegenda/go-adjacent-set"
)

func main() {
    asm := gas.AdjacentSetMetadata{}
}
````

`gas` is a funny but an unfortunate abbreviation. Choose whichever you fancy.

## 3. Important to know before using it

First of all, I just started learning Go. This package is a result of between 20 to 30 hours of work so have that in mind. It still
needs some furbishing but it will get there. Keep that in mind when taking a look at this libray.

This package does not handle connecting to the database. That is your job. You only have to pass the table name information and
the created handle to `AdjacentSetMetadata` before using it. It will not tamper with your connection in any way. 

It also does not give you the full tree as a result. Instead, it lazily loads the results at the time you ask for it. The consequence
is multiple trips to MySQL server but only with single `select` statements. No `joins` were used. This might or might not be a good thing
for you and depends what you need it for. A lot of my experience with SQL tree structures come from the Symfony ecosystem that uses
[this](https://github.com/Atlantic18/DoctrineExtensions/blob/v2.4.x/doc/tree.md) library. I used it mainly for building categories of products,
for example, that were also menus on the frontend. If you need something like that in Go, this package could be useful for you. If you 
need it for something much more complicated, you have two options. Create an SQL query by yourself which is pretty simple to make but can 
also be found with a simple google search. The second option is just to skip this package entirely. 

This package also integrates seamlessly with your database. It creates a table whose name you choose with your already existing database.
The table that it creates has auto increment primary key `id`, a `name` field that must be *unique* and a `parent` field that is also 
an integer that defaults to `NULL`. It also creates an index on name and parent fields which is really important if you alter the table later
on in some way. If you do that, don't forget that there is also an index that needs to reflect your changes.

Also, I haven't written any code to delete a node. I'm still thinking about it. Should I remove all the subtree
nodes of a parent node, delete only the one that you want to delete and keep the children orphans, soft delete them i.e. place a `deletedAt` field
making it unqueriable but still in the database if you wish to restore them later on, or should I rearrange them in a way you see fit? Those
are the questions that I hope to find the answers to. 

I hope you like it.

## 4. Built in commands

There are two built in commands that this package provides: `./createTable` and `./seed`. These commands will not be in your *bin* directory
but in `src/github.com/MarioLegenda/go-adjacent-set` directory. In order to use them, cd into that directory.

The usage is roughly the same for both of them. As I said earlier, this package presumes that you have a working database already
in place. In order to use `./createTable`, you need to pass all the information to connect to the database and the name of the table to use.
The full command looks like this...

````
./createTable -u root -p root -d dbName -t tableName
````

- *-u* is the database username
- *-p* is the database password
- *-d* is the database name
- *-t* is the table name to create

This command creates the table `tableName` with these fields:

- *tableName_id AUTO_INCREMENT PRIMARY KEY*
- *name VARCHAR(255) NOT NULL*
- *parent INT DEFAULT NULL*

It also adds a unique constraint called *unique_name* to (you guessed it) the *name* field. There is also an index created with this
basic statement:

````mysql
CREATE INDEX idx_name_and_parent ON tableName (name, parent)
````

The `./seed` command is the same but it accepts additional *-h* for depth and *-l* for the number of leafs. The seed command
creates some additional data for you to try this package out. 

**Do not run this command on the production server**

- *-h* tells the command how many sublevels the created tree should have. 
- *-l* tells the command the number of leafs every sublevel has. 

If you put *-h 5* and *-l 5*, that basically means *pow(5, 5)* so be careful with this command. I tried *-h 15 -l 10*. It took some time.
Both *-h* and *-l* default to 5. The entire command looks like this

````
./seed -u root -p root -d dbName -t tableName -h 5 -l 5
````

## 5. Usage

The main type to use is the `AdjacentSetMetadata` that accepts an sql connection handle and the table name. In these examples,
I will make the errors `panic` for the sake of brevity but you should handle them however you wish.

````go
import(
	gas "github.com/MarioLegenda/go-adjacent-set"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("root:root@/%s", dbName))

	if err != nil {
		panic(err)
	}

    asm := gas.AdjacentSetMetadata{
        handle: db,
        tableName: "tableName"
    }
}
````

`AdjacentSetMetadata` has all the regular CRUD operations (except for deletion, see *Important to know before using it* section). 

- `CreateRoot` creates the root node 
- `Root` returns the the `RootCategory` struct with the created root node data and an error if an error occurred, `nil` otherwise
- `CreateCategory` accepts a `Category` struct with the neccessary data to create a category. It returns the database id in
   which it is created and an error if an error occurred
- `CategoryExists` accepts the same parameters as `CreateCategory` but returns a bool and an error if an error occurred
- `GetCategory` accepts a `FetchOptions` struct and returns an `AdjacentSetResult` if it found one and an error if an error occurred

We start by creating a root node.

````go
    asm := gas.AdjacentSetMetadata{
        handle: db,
        tableName: "tableName"
    }

    // since the shorthand := assignament creates a basic int, we need to declare the id as int64 before hand
    // since the sql package works only with int64
    var id int64
    // CreateRoot only accepts the root name. Name it whatever you like
    id, error := asm.CreateRoot("Root")

    if error != nil { panic(error) }
````

Under the hood, `CreateRoot` check if the root node already exists. If you know that the root node is already created, you can skip 
the error returned from `CreateRoot` so it is safe to call `CreateRoot` multiple times.

To check if the root node exists, use the `Root` method. One caveat with this method is that if the root node does not exist, it
returns a `RootCategory` struct with `Name` and `Id` fields with default values (empty string and 0). I had to do it like this
because if I return an error object, go complains about memory segmentation for some reason.

````go
root, err := asm.Root()

if root.Name == "" && root.Id == 0 {
    // root does not exist
}
````

After that, you can start created your categories.

`Nodes or categories, name it however you like. The methods on the AdjacentSetMetadata names them categories`

````go
    asm := gas.AdjacentSetMetadata{
        handle: db,
        tableName: "tableName"
    }

    // since the shorthand := assignament creates a basic int, we need to declare the id as int64 before hand
    // since the sql package works only with int64
    var id int64
    // CreateRoot only accepts the root name. Name it whatever you like
    id, error := asm.CreateRoot("Root")

    if error != nil { panic(error) }

    categoryId, err := asm.CreateCategory(gas.Category{
        Name: "First category",
        CategoryId: id
    })
````

And you have created your first category. After you created your first category, create and add more categories to that category

````go
    categoryId, err := asm.CreateCategory(gas.Category{
        Name: "First category",
        // id is the root id from our previous example
        CategoryId: id
    })

    for i := range []int{1, 2, 3, 4, 5} {
    	_, err := asm.CreateCategory(gas.Category{
    		Name: fmt.Sprintf("Category_%d", i),
                // id is the root id from our previous example
            CategoryId: categoryId
        })
    }
````

In the above example, we have just created 5 categories as children to our root category names `Category_{0-4}`. 

### 5.1 Traversing the tree



















