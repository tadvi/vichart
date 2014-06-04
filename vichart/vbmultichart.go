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
	VBMultiLineXYStyle = "stroke:lightgray;stroke-width:2px;"
	VBMultiLineStyle   = "fill:navy;stroke:navy;stroke-width:2px;"
	VBMultiBarStyle1   = "fill:green;stroke:gray;"
	VBMultiBarStyle2   = "fill:yellow;stroke:gray;"
	VBMultiBarStyle3   = "fill:white;stroke:gray;"

	VBMultiGstyle      = "font-family:Calibri; font-size:14"
	VBMultiGutterLeft  = 40
	VBMultiGutterRight = 40
	VBMultiGutterTop   = 40

	VBMultiBarSpacing    = 16
	VBMultiBarWidth      = 15
	VBMultiLegendXOffset = 10
)

type VBMultiChart struct {
	Svg           *svg.SVG
	Width, Height int
	BarValues     []VBMultiChartItem // chart bar values
	LineValues    []int              // chart line values
	MaxBarValue   int                // chart max value, used for scaling all the bar values
	MaxLineValue  int                // chart max value, used for scaling all the line values

	// optional fields below
	BarSpacing int
	BarWidth   int
	LabelsX    []string
	LabelsY1   []string
	LabelsY2   []string

	GutterLeft  int
	GutterRight int // right gutter for the chart, used to fit last bottom label
	GutterTop   int // top gutter for the chart, used top label

	// styles
	Gstyle      string
	LineXYStyle string
	LineStyle   string
	BarStyle1   string
	BarStyle2   string
	BarStyle3   string

	// legend
	BarLegend1 string
	BarLegend2 string
	BarLegend3 string
	LineLegend string

	// legend offset
	LegendXOffset int
}

type VBMultiChartItem struct {
	Bottom, Middle, Top int
}

// Draw produces chart on screen, main entry point.
func (chart *VBMultiChart) Draw() error {
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
	if chart.MaxBarValue == 0 {
		return fmt.Errorf("Missing chart MaxBarValue.")
	}
	if chart.MaxLineValue == 0 && len(chart.LineValues) > 0 {
		return fmt.Errorf("Missing chart MaxLineValue.")
	}
	if len(chart.BarValues) != len(chart.LineValues) {
		return fmt.Errorf("Number of BarValues does not match number of LineValues.")
	}
	// default to sensible constants if value is not set
	if chart.LineXYStyle == "" {
		chart.LineXYStyle = VBMultiLineXYStyle
	}
	if chart.Gstyle == "" {
		chart.Gstyle = VBMultiGstyle
	}
	if chart.LineStyle == "" {
		chart.LineStyle = VBMultiLineStyle
	}
	if chart.BarStyle1 == "" {
		chart.BarStyle1 = VBMultiBarStyle1
	}
	if chart.BarStyle2 == "" {
		chart.BarStyle2 = VBMultiBarStyle2
	}
	if chart.BarStyle3 == "" {
		chart.BarStyle3 = VBMultiBarStyle3
	}
	if chart.GutterLeft == 0 {
		chart.GutterLeft = VBMultiGutterLeft
	}
	if chart.GutterRight == 0 {
		chart.GutterRight = VBMultiGutterRight
	}
	if chart.GutterTop == 0 {
		chart.GutterTop = VBMultiGutterTop
	}
	if chart.BarSpacing == 0 {
		chart.BarSpacing = VBMultiBarSpacing
	}
	if chart.BarWidth == 0 {
		chart.BarWidth = VBMultiBarWidth
	}
	if chart.LegendXOffset == 0 {
		chart.LegendXOffset = VBMultiLegendXOffset
	}

	// start SVG
	canvas.Start(chart.Width, chart.Height)
	canvas.Gstyle(chart.Gstyle)
	x, y := chart.GutterLeft, chart.Height-42
	bHeight := float64(y - chart.GutterTop)
	bWidth := float64(chart.Width - chart.GutterRight - x)

	xoffset := x
	for i, _ := range chart.BarValues {
		yoffset := y + 3
		// scale value to fit in chart pixels
		chartVal := chart.calcBarValue(bHeight, chart.BarValues[i].Bottom)
		chart.drawMeter(xoffset, yoffset, chart.BarWidth, chartVal, chart.BarStyle1)
		yoffset -= chartVal

		chartVal = chart.calcBarValue(bHeight, chart.BarValues[i].Middle)
		chart.drawMeter(xoffset, yoffset, chart.BarWidth, chartVal, chart.BarStyle2)
		yoffset -= chartVal

		chartVal = chart.calcBarValue(bHeight, chart.BarValues[i].Top)
		chart.drawMeter(xoffset, yoffset, chart.BarWidth, chartVal, chart.BarStyle3)
		yoffset -= chartVal

		// draw line on the chart
		if i > 0 && len(chart.LineValues) > 0 {
			valLine1 := float64(chart.LineValues[i-1])
			chartValLine1 := int((valLine1 / float64(chart.MaxLineValue)) * bHeight)
			valLine2 := float64(chart.LineValues[i])
			chartValLine2 := int((valLine2 / float64(chart.MaxLineValue)) * bHeight)

			xpos := xoffset + chart.BarWidth/2
			canvas.Line(xpos-chart.BarSpacing, y-chartValLine1+3, xpos, y-chartValLine2+3, chart.LineStyle)
		}
		xoffset += chart.BarSpacing
	}

	// bottom line markers and labels
	canvas.Line(x, y+12, chart.Width-chart.GutterRight, y+12, chart.LineXYStyle)
	labels := len(chart.LabelsX)
	// display bottom line labels
	for i := 0; i < labels; i++ {
		step := bWidth / float64(labels-1)
		xoffset := int(float64(i) * step)
		canvas.Text(x+xoffset, y+30, chart.LabelsX[i], "font-size:75%;text-anchor:middle;")
		canvas.Line(x+xoffset, y+6, x+xoffset, y+18, chart.LineXYStyle)
	}

	// left vertical Y line
	chart.drawYLine(x, y+2)
	chart.drawYLineText(x-16, y, chart.LabelsY1, true)
	// right vertical Y line
	chart.drawYLine(chart.Width-chart.GutterRight+12, y+2)
	chart.drawYLineText(chart.Width-chart.GutterRight+12, y, chart.LabelsY2, false)

	chart.drawLegend(x)

	canvas.Gend()
	canvas.End()
	return nil
}

func (chart *VBMultiChart) calcBarValue(bHeight float64, value int) int {
	val := float64(value)
	chartVal := int((val / float64(chart.MaxBarValue)) * bHeight)
	return chartVal
}

// drawLegend produces legend on the chart.
func (chart *VBMultiChart) drawLegend(x int) {
	canvas := chart.Svg
	canvas.Rect(x+chart.LegendXOffset, 10, 40, 10, chart.BarStyle1)
	canvas.Text(x+chart.LegendXOffset+50, 20, chart.BarLegend1, "font-size:75%;")

	canvas.Rect(x+chart.LegendXOffset+120, 10, 40, 10, chart.BarStyle2)
	canvas.Text(x+chart.LegendXOffset+170, 20, chart.BarLegend2, "font-size:75%;")

	canvas.Rect(x+chart.LegendXOffset+230, 10, 40, 10, chart.BarStyle3)
	canvas.Text(x+chart.LegendXOffset+280, 20, chart.BarLegend3, "font-size:75%;")

	if chart.LineLegend != "" {
		xpos := x + chart.LegendXOffset + 340
		canvas.Line(xpos, 15, xpos+40, 15, chart.LineStyle)
		canvas.Text(x+chart.LegendXOffset+390, 20, chart.LineLegend, "font-size:75%;")
	}
}

// drawYLineText draws Y line text.
func (chart *VBMultiChart) drawYLineText(x, h int, labels []string, left bool) {
	canvas := chart.Svg

	style := "font-size:75%;text-anchor:start;baseline-shift:-75%"
	if left {
		style = "font-size:75%;text-anchor:end;baseline-shift:-75%"
	}
	labelsCount := len(labels)
	for i := 0; i < labelsCount; i++ {
		step := float64(h-chart.GutterTop) / float64(labelsCount-1)
		yoffset := int(float64(i) * step)
		canvas.Text(x, yoffset+chart.GutterTop, labels[labelsCount-i-1], style)
	}
}

// drawYLine draws Y line on screen.
func (chart *VBMultiChart) drawYLine(x, y int) {
	canvas := chart.Svg
	canvas.Line(x-8, chart.GutterTop, x-8, y, chart.LineXYStyle)

	height := float64(y - chart.GutterTop)
	step := height / 10
	pos := 0
	for i := 0.0; i <= height; i += step {
		marker := int(height-i) + 1 + chart.GutterTop
		if pos == 0 || pos == 5 || pos == 10 {
			canvas.Line(x-2, marker, x-14, marker, chart.LineXYStyle)
		} else {
			canvas.Line(x-5, marker, x-11, marker, chart.LineXYStyle)
		}
		pos += 1
	}
}

// drawMeter draws bar on chart.
func (chart *VBMultiChart) drawMeter(x, y, w, value int, barStyle string) {
	canvas := chart.Svg
	corner := w
	canvas.Roundrect(x, y-value, corner, value, 0, 0, barStyle)
}
