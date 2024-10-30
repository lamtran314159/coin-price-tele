package handlers

import (
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type klineData struct {
	Date string     `json:"date"`
	Data [4]float32 `json:"data"`
}

// Create the Kline chart in memory
func klineBase(kd []klineData) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].Date)
		y = append(y, opts.KlineData{Value: kd[i].Data})
	}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Kline On Demand",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 5,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("kline", y)
	return kline
}

// Helper function to convert string to float32
func parseFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}
