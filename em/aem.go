package em

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ut                 = "6d2ffaa6a585d612eda28417681d58fb"
	shortChgHqFields   = "f3,f12,f13"
	shortCloseHqFields = "f2,f3,f12,f13,f18"
	shortOhlcHqFields  = "f1,f2,f3,f5,f12,f13,f14,f17,f15,f16,f17,f18"
	longHqFields       = "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f31,f32,f33,f41,f46,f62,f100,f115"
	finHqFields        = "f1,f9,f12,f13,f14,f20,f21,f41,f46,f100,f115"
)

var hqColDict = map[string][]interface{}{
	"-":    {"无名", "string", false},
	"f1":   {"zoom", "float64", false},
	"f2":   {"close", "float64", true},
	"f3":   {"pct_chg", "float64", true},
	"f4":   {"change", "float64", true},
	"f5":   {"vol", "float64", false},
	"f6":   {"amount", "float64", false},
	"f7":   {"zhenfu", "float64", true},
	"f8":   {"huanshoulv", "float64", true},
	"f9":   {"pe", "float64", true},
	"f10":  {"liangbi", "float64", true},
	"f11":  {"chg5m", "float64", false},
	"f12":  {"symbol", "string", false},
	"f13":  {"market", "string", false},
	"f14":  {"name", "string", false},
	"f15":  {"high", "float64", true},
	"f16":  {"low", "float64", true},
	"f17":  {"open", "float64", true},
	"f18":  {"pre_close", "float64", true},
	"f20":  {"tmv", "float64", false},
	"f21":  {"fmv", "float64", false},
	"f22":  {"zhangsu", "float64", true},
	"f23":  {"pbr", "float64", false},
	"f24":  {"chg60d", "float64", true},
	"f25":  {"chgnch", "float64", true},
	"f31":  {"buy", "float64", true},
	"f32":  {"sell", "float64", true},
	"f33":  {"weibi", "float64", true},
	"f41":  {"trqchg", "float64", true},
	"f46":  {"npqchg", "float64", true},
	"f62":  {"zhulijlr", "float64", false},
	"f100": {"industry", "string", false},
	"f115": {"pettm", "float64", true},
}

func GetCode(x string) string {
	if strings.HasPrefix(x, "6") {
		return x + ".SH"
	} else if strings.HasPrefix(x, "8") || strings.HasPrefix(x, "4") {
		return x + ".BJ"
	}
	return x + ".SZ"
}

type StockData struct {
	TsCode   string  `json:"ts_code"`
	Name     string  `json:"name"`
	Close    float64 `json:"close"`
	PctChg   float64 `json:"pct_chg"`
	Vol      float64 `json:"vol"`
	Amount   float64 `json:"amount"`
	Industry string  `json:"industry"`
}

func GetStockList(market string) ([]StockData, error) {
	url := fmt.Sprintf("http://80.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=5000&po=1&np=1&ut=%s&fltt=2&invt=2&fid=f3&fs=m:%s&fields=%s", ut, market, finHqFields)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].(map[string]interface{})
	diff := data["diff"].([]interface{})

	var stockList []StockData
	for _, item := range diff {
		stock := item.(map[string]interface{})
		stockData := StockData{
			TsCode:   GetCode(stock["f12"].(string)),
			Name:     stock["f14"].(string),
			Close:    stock["f2"].(float64),
			PctChg:   stock["f3"].(float64),
			Vol:      stock["f5"].(float64),
			Amount:   stock["f6"].(float64),
			Industry: stock["f100"].(string),
		}
		stockList = append(stockList, stockData)
	}

	return stockList, nil
}

func GetStockHq(codes []string, fields string) (map[string]map[string]interface{}, error) {
	secids := make([]string, len(codes))
	for i, code := range codes {
		if strings.HasSuffix(code, ".SH") {
			secids[i] = "1." + strings.TrimSuffix(code, ".SH")
		} else if strings.HasSuffix(code, ".SZ") {
			secids[i] = "0." + strings.TrimSuffix(code, ".SZ")
		} else if strings.HasSuffix(code, ".BJ") {
			secids[i] = "2." + strings.TrimSuffix(code, ".BJ")
		}
	}

	params := url.Values{}
	params.Add("ut", ut)
	params.Add("fltt", "2")
	params.Add("invt", "2")
	params.Add("fields", fields)
	params.Add("secids", strings.Join(secids, ","))

	url := "http://push2.eastmoney.com/api/qt/ulist/get?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].(map[string]interface{})
	diff := data["diff"].([]interface{})

	hqData := make(map[string]map[string]interface{})
	for _, item := range diff {
		stock := item.(map[string]interface{})
		code := GetCode(stock["f12"].(string))
		hqData[code] = make(map[string]interface{})
		for k, v := range stock {
			if colInfo, ok := hqColDict[k]; ok {
				colName := colInfo[0].(string)
				colType := colInfo[1].(string)
				if colType == "float64" {
					hqData[code][colName] = v.(float64)
				} else {
					hqData[code][colName] = v.(string)
				}
			}
		}
	}

	return hqData, nil
}

func GetStockKline(code string, period string, limit int) ([]map[string]interface{}, error) {
	var secid string
	if strings.HasSuffix(code, ".SH") {
		secid = "1." + strings.TrimSuffix(code, ".SH")
	} else if strings.HasSuffix(code, ".SZ") {
		secid = "0." + strings.TrimSuffix(code, ".SZ")
	} else if strings.HasSuffix(code, ".BJ") {
		secid = "2." + strings.TrimSuffix(code, ".BJ")
	}

	params := url.Values{}
	params.Add("ut", ut)
	params.Add("fields1", "f1,f2,f3,f4,f5,f6")
	params.Add("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61")
	params.Add("klt", period)
	params.Add("fqt", "1")
	params.Add("secid", secid)
	params.Add("beg", "0")
	params.Add("end", "20500101")
	params.Add("lmt", strconv.Itoa(limit))

	url := "http://push2his.eastmoney.com/api/qt/stock/kline/get?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].(map[string]interface{})
	klines := data["klines"].([]interface{})

	var klineData []map[string]interface{}
	for _, kline := range klines {
		items := strings.Split(kline.(string), ",")
		klineItem := map[string]interface{}{
			"date":   items[0],
			"open":   parseFloat(items[1]),
			"close":  parseFloat(items[2]),
			"high":   parseFloat(items[3]),
			"low":    parseFloat(items[4]),
			"vol":    parseFloat(items[5]),
			"amount": parseFloat(items[6]),
		}
		klineData = append(klineData, klineItem)
	}

	return klineData, nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func GetTradeDate() (string, error) {
	now := time.Now()
	if now.Weekday() == time.Saturday {
		now = now.AddDate(0, 0, -1)
	} else if now.Weekday() == time.Sunday {
		now = now.AddDate(0, 0, -2)
	}

	if now.Hour() < 15 {
		now = now.AddDate(0, 0, -1)
	}

	return now.Format("20060102"), nil
}
