package aem

import (
	"fmt"

	em "github.com/zhnt/zhfetch-go/em"
)

func TestGetStockList() {
	// 准备测试数据
	expectedStocks := []em.Stock{
		{Code: "600000", Name: "浦发银行"},
		{Code: "000001", Name: "平安银行"},
	}

	// 调用被测试的方法
	stocks, _ := em.GetStockList()

	// 检查返回的股票列表长度
	if len(stocks) != len(expectedStocks) {
		fmt.Errorf("GetStockList() 返回的股票数量不正确, 期望 %d, 实际 %d", len(expectedStocks), len(stocks))
	}

	// 检查返回的股票数据
	for i, stock := range stocks {
		if stock.Code != expectedStocks[i].Code || stock.Name != expectedStocks[i].Name {
			fmt.Errorf("股票 #%d 不匹配, 期望 %+v, 实际 %+v", i, expectedStocks[i], stock)
		}
	}
}
