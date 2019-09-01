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









