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

