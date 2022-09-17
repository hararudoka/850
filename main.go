package main

import (
	"dict/internal/dict"
)

func main() {
	d, err := dict.New850()
	if err != nil {
		panic(err)
	}

	for i := range d {
		_ = i
	}

	d.ToFile("temporary")
}
