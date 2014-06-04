// ViChart library for Go
// Author: Tad Vizbaras 
// License: http://github.com/tadvi/vichart/blob/master/LICENSE 
//
package vichart

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

const (
	HBarLineXStyle  = "stroke:lightgray;stroke-width:2px;"
	HBarGstyle      = "font-family:Calibri; font-size:14"
	HBarGutterLeft  = 100
	HBarGutterRight = 20
	HBarSpacing     = 18
)

type HBarChart struct {
	Svg           *svg.SVG
	Width, Height int
	BarValues     []int // chart values
	MaxValue      int   // chart max value, used for scaling all the display values
	LabelsY       []string

	// optional fields below
	BarSpacing  int
	LabelsX     []string
	GutterLeft  int
	GutterRight int // right gutter for the chart, used to fit last bottom label

	// styles
	Gstyle     string
	LineXStyle string
}

// Draw main entry.
func (chart *HBarChart) Draw() error {
	canvas := chart.Svg
	if chart.Svg == nil {
		return fmt.Errorf("Missing pointer to svg.SVG in field Svg.")
	}
	if chart.Width < 10 || chart.Height < 10 {
		return fmt.Errorf("Incorrect Width or Height value.")
	}
	if len(chart.BarValues) == 0 {
		return fmt.Errorf("Missing BarValues for the chart.")
	}
	if chart.MaxValue == 0 {
		return fmt.Errorf("Missing chart MaxValue.")
	}
	if len(chart.BarValues) != len(chart.LabelsY) {
		return fmt.Errorf("Number of BarValues does not match number of LabelY.")
	}
	// default to sensible constants if value is not set
	if chart.LineXStyle == "" {
		chart.LineXStyle = HBarLineXStyle
	}
	if chart.Gstyle == "" {
		chart.Gstyle = HBarGstyle
	}
	if chart.GutterRight == 0 {
		chart.GutterRight = HBarGutterRight
	}
	if chart.BarSpacing == 0 {
		chart.BarSpacing = HBarSpacing
	}
	if chart.GutterLeft == 0 {
		chart.GutterLeft = HBarGutterLeft
	}

	// start SVG
	canvas.Start(chart.Width, chart.Height)
	canvas.Gstyle(chart.Gstyle)
	x, y := chart.GutterLeft, 5
	bWidth := float64(chart.Width - chart.GutterRight - x)

	for i, data := range chart.LabelsY {
		// scale value to fit in chart pixels
		val := float64(chart.BarValues[i])
		chartVal := int((val / float64(chart.MaxValue)) * bWidth)
		fmt.Printf("%f scaled to %d\n", val, chartVal)
		chart.drawMeter(x, y, chart.Width-x, chart.BarSpacing, chartVal,
			chart.BarValues[i], data)
		y += chart.BarSpacing
	}

	// bottom line markers and label
	canvas.Line(x, y+12, chart.Width-chart.GutterRight, y+12, chart.LineXStyle)
	step := bWidth / 10
	pos := 0
	for i := 0.0; i <= bWidth; i += step {
		marker := x + int(i)
		if pos == 0 || pos == 5 || pos == 10 {
			canvas.Line(marker, y+6, marker, y+18, chart.LineXStyle)
		} else {
			canvas.Line(marker, y+9, marker, y+15, chart.LineXStyle)
		}
		pos += 1
	}
	// display bottom line labels
	labels := len(chart.LabelsX)
	for i := 0; i < labels; i++ {
		step := bWidth / float64(labels-1)
		xoffset := int(float64(i) * step)
		canvas.Text(x+xoffset, y+30, chart.LabelsX[i], "font-size:75%;text-anchor:middle;")
	}
	canvas.Gend()
	canvas.End()
	return nil
}

// drawMeter draw bar on screen.
func (chart *HBarChart) drawMeter(x, y, w, h, value, origValue int, label string) {
	canvas := chart.Svg
	corner := h / 2
	inset := corner / 2
	canvas.Text(x-5, y+h/2, label, "text-anchor:end;baseline-shift:-33%")
	canvas.Roundrect(x, y+inset, value, h-(inset*2),
		inset, inset, "fill:darkgray")
	if value > 9 {
		// draw inset circle only if value is not too small
		canvas.Circle(x+inset+value-corner, y+corner, inset, "fill:red;fill-opacity:0.3")
	}
	canvas.Text(x+inset+value+2, y+h/2, fmt.Sprintf("%-3d", origValue),
		"font-size:75%;text-anchor:start;baseline-shift:-33%")
}
