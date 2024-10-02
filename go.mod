module github.com/zhnt/ipquant/zhfetch

require (
	github.com/zhnt/ipquant/gozhfetch v0.0.0
)

replace github.com/zhnt/ipquant/gozhfetch => ./gozhfetch

go 1.22.7
