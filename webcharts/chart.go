package main

import (
	"github.com/ajstarks/svgo"
	"log"
	"math/rand"
	"net/http"
	"time"
	"vichart"
	"strconv"
)

func main() {
	http.Handle("/piechart", http.HandlerFunc(piechart))
	http.Handle("/hchart", http.HandlerFunc(hchart))
	http.Handle("/vchart", http.HandlerFunc(vchart))
	http.Handle("/vbmultichart", http.HandlerFunc(vbmultichart))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// piechart draws piechart.
func piechart(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canvas := svg.New(w)
	rand.Seed(int64(time.Now().Second()))

	chart := vichart.PieChart{
		Svg:	canvas,
		Width:	650,
		Height: 400,
		PieValues: []int{},
		Labels: []string{},		
		LegendXOffset: 250,
		GutterLeft: 40,
		GutterTop: 40,
	}
	
	for i:=0; i < 8; i++ {
		val := rand.Intn(100)
		chart.PieValues = append(chart.PieValues, val)
		chart.Labels = append(chart.Labels, strconv.Itoa(val))		
	}
	vichart.Must(chart.Draw())	
	
}

// vchart draws bar and line chart.
func vchart(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canvas := svg.New(w)
	rand.Seed(int64(time.Now().Second()))

	chart := vichart.VBarChart{
		Svg:           canvas,
		Width:         650,
		Height:        400,
		LabelsY1:      []string{"Zero", "Half", "Full"},    // bar labels
		LabelsY2:      []string{"Zero2", "Half2", "Full2"}, // line labels
		LabelsX:       []string{"0", "1/4", "1/2", "3/4", "1"},
		BarValues:     []int{},
		LineValues:    []int{},
		MaxBarValue:   3000,
		MaxLineValue:  3000,
		BarLegend:     "Speed",
		LineLegend:    "Rpm",
		BarWidth:      14,
		BarSpacing:    18,  // should be greater than BarWidth
		LegendXOffset: 100, // legend offset from the left
		GutterRight:   60,
		GutterLeft:    65,
	}
	// populate chart with data
	for i := 0; i < 12; i++ {
		val := rand.Intn(chart.MaxBarValue)
		chart.BarValues = append(chart.BarValues, val)
		chart.LineValues = append(chart.LineValues, val)
	}

	// panics if Draw fails
	vichart.Must(chart.Draw())
}

// vbmultichart draws multibar and line chart.
func vbmultichart(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canvas := svg.New(w)
	rand.Seed(int64(time.Now().Second()))

	chart := vichart.VBMultiChart{
		Svg:          canvas,
		Width:        572,
		Height:       400,
		LabelsY1:     []string{"0", "50", "100"},
		LabelsY2:     []string{"0", "500", "1000"},
		LabelsX:      []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
		BarValues:    []vichart.VBMultiChartItem{},
		LineValues:   []int{},
		MaxBarValue:  3000,
		MaxLineValue: 3000,
		BarWidth:     35,
		BarSpacing:   39,
		BarLegend1:   "Driving",
		BarLegend2:   "Idle",
		BarLegend3:   "Off",
		LineLegend:   "Distance",
		GutterRight:  60,
		GutterLeft:   45,
	}
	// populate chart with data
	for i := 0; i < 12; i++ {
		val1 := rand.Intn(chart.MaxBarValue / 2)
		val2 := rand.Intn(chart.MaxBarValue / 3)
		val3 := 3000 - val1 - val2

		chart.BarValues = append(chart.BarValues, vichart.VBMultiChartItem{val1, val2, val3})
		chart.LineValues = append(chart.LineValues, val1*2)
	}

	vichart.Must(chart.Draw())
}

// hchart draws horizontal chart.
func hchart(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canvas := svg.New(w)
	rand.Seed(int64(time.Now().Second()))

	chart := vichart.HBarChart{
		Svg:     canvas,
		Width:   550,
		Height:  200,
		LabelsY: []string{"Cost", "Priorities", "Timing", "Technology"},
		//Spacing: 18,
		LabelsX:   []string{"0", "1/4", "1/2", "3/4", "3000"},
		BarValues: []int{},
		MaxValue:  3000,
	}
	for i := 0; i < len(chart.LabelsY); i++ {
		chart.BarValues = append(chart.BarValues, rand.Intn(chart.MaxValue))
	}

	vichart.Must(chart.Draw())
}
