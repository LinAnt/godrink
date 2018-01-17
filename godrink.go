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
	c, err := r.Read()
	if err != nil {
		panic(err)
	}
	for k := range c {
		if k.Type == reader.EvKEY {
			fmt.Println(k.KeyString())
		}
	}
}
