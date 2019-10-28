# go-yit - YAML Iterator

## Introduction

This library compliments [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) by adding
functional style methods for iterating over YAML documents.

## Usage

Import the package
```go
import "github.com/dprotaso/go-yit"
```


Query your YAML document
```go
package main

import (
	"fmt"
	"log"

	"github.com/dprotaso/go-yit"
	"gopkg.in/yaml.v3"
)

var data = `
a: b
c: d
e: f
`

func main() {
	var doc yaml.Node
	err := yaml.Unmarshal([]byte(data), &doc)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	it := yit.FromNode(&doc).
		RecurseNodes().
		Filter(yit.WithKind(yaml.MappingNode)).
		MapKeys()

	for node, ok := it(); ok; node, ok = it() {
		fmt.Println(node.Value)
	}
}
```

