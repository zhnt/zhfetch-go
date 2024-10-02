package main

import (
	"fmt"

	aem "github.com/zhnt/zhfetch-go/em"
)

func main() {

	stockList, err := aem.GetStockList("SZ")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stockList)
}
