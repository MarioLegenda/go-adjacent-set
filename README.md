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

This package does not handle connecting to the database. That is your job. You only have to pass the table name information and
the created handle to `AdjacentSetMetadata` before using it. It will not tamper with your connection in any way. 

It also does not give you the full tree as a result. Instead, it lazily loads the results at the time you ask for it. The consequence
is multiple trips to MySQL server but only with single `select` statements. No `joins` were used. This might or might not be a good thing
for you and depends what you need it for. A lot of my experience with SQL tree structures come from the Symfony ecosystem that uses
[this](https://github.com/Atlantic18/DoctrineExtensions/blob/v2.4.x/doc/tree.md) library. I used it mainly for building categories of products,
for example, that were also menus on the frontend. If you need something like that in Go, this package could be useful for you. If you 
need it for something much more complicated, you have two options. Create an SQL query by yourself which is pretty simple to make but can 
also be found with a simple google search. The second option is just to skip this package entirely. 

The choice is yours. I hope you like it.

## 4. Built in commands





