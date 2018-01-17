package main

import (
	"fmt"

	"github.com/LinAnt/godrink/reader"
)

const (
	readerName = "StrongLink USB CardReader"
)

func main() {
	r, err := reader.GetReader(readerName)
	if err != nil {
		panic(err)
	}
	c, err := r.GetCardChannel()
	if err != nil {
		panic(err)
	}
	for id := range c {
		fmt.Println(id)
	}
}
