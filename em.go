package main

import (
	"fmt"

	aem "github.com/zhnt/ipquant/gozhfetch/em/aem"
)

func main() {

	stockList, err := aem.GetStockList("SZ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stockList)
}
