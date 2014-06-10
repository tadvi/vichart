// ViChart library for Go
// Author: Tad Vizbaras 
// License: http://github.com/tadvi/vichart/blob/master/LICENSE 
//
package vichart

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"math"
)

const (
	PieGstyle      = "font-family:Calibri; font-size:14"
	PieStyle = "fill:white;stroke:black;stroke-width:2px;"
	
	PieFillStyle1    = "fill:red;stroke:gray;"
	PieFillStyle2    = "fill:green;stroke:gray;"
	PieFillStyle3    = "fill:navy;stroke:gray;"
	PieFillStyle4    = "fill:orange;stroke:gray;"
	PieFillStyle5    = "fill:gray;stroke:gray;"	
	PieFillStyle6    = "fill:white;stroke:gray;"	
	PieFillStyle7    = "fill:blue;stroke:gray;"	
	PieFillStyle8    = "fill:yellow;stroke:gray;"		
			
	PieGutterLeft  = 40
	PieGutterTop   = 40
	
	PieRadius = 80

	PieLegendXOffset = 40
)

type PieChart struct {
	Svg           *svg.SVG
	Width, Height int
	PieValues     []int // chart bar values
	Labels []string
	Radius int
	
	// optional fields below
	FillStyles   []string	
	
	GutterLeft  int // left gutter for the chart, used to fit left labels
	GutterTop   int // top gutter for the chart, used top label

	// styles
	Gstyle      string
	PieStyle string
	//LineStyle   string

	// legend related
	Legend  string
	// legend offset
	LegendXOffset int
}

// Draw produces chart on screen, main entry point.
func (chart *PieChart) Draw() error {
	canvas := chart.Svg
	if chart.Svg == nil {
		return fmt.Errorf("Missing pointer to svg.SVG in field Svg.")
	}
	if chart.Width < 10 || chart.Height < 10 {
		return fmt.Errorf("Incorrect Width or Height value.")
	}
	if len(chart.PieValues) == 0 {
		return fmt.Errorf("Missing PieValues for the chart.")
	}	
	if len(chart.PieValues) != len(chart.Labels) {
		return fmt.Errorf("Number of PieValues does not match number of Labels.")
	}
	
	if len(chart.FillStyles) == 0 {	// fill styles not set use defaults
		chart.FillStyles = append(chart.FillStyles, PieFillStyle1)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle2)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle3)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle4)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle5)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle6)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle7)
		chart.FillStyles = append(chart.FillStyles, PieFillStyle8)		
	}
	// default to sensible constants if value is not set
	if chart.Gstyle == "" {
		chart.Gstyle = PieGstyle
	}
	if chart.PieStyle == "" {
		chart.PieStyle = PieStyle
	}
	if chart.Radius == 0 {
		chart.Radius = PieRadius
	}
	
	if chart.GutterTop == 0 {
		chart.GutterTop = PieGutterTop
	}	
	if chart.GutterLeft == 0 {
		chart.GutterLeft = PieGutterLeft
	}
	
	if chart.LegendXOffset == 0 {
		chart.LegendXOffset = PieLegendXOffset
	}
	
	// convert values into degrees
	sum := 0.0
	for _, val := range chart.PieValues {
		sum += float64(val)
	}
	angles := make([]float64, len(chart.PieValues), len(chart.PieValues))
	for i, val := range chart.PieValues {
		angles[i] = float64(val) * 360 / sum
	}

	// start SVG
	canvas.Start(chart.Width, chart.Height)
	canvas.Gstyle(chart.Gstyle)

	// cx, cy - center of the pie
	cx := chart.GutterLeft + chart.Radius
	cy := chart.GutterTop + chart.Radius
		
	var startAngle, endAngle float64
	// draw each slice in the loop
	for i, val := range angles {
		startAngle = endAngle
		endAngle = startAngle + val
		
		radius := float64(chart.Radius)
		bx := cx + int(radius * math.Cos(math.Pi * startAngle / 180))
		by := cy + int(radius * math.Sin(math.Pi * startAngle / 180))
		endx := cx + int(radius * math.Cos(math.Pi * endAngle / 180))
		endy := cy + int(radius * math.Sin(math.Pi * endAngle / 180))
		
		path := fmt.Sprintf("M%d,%d  L%d,%d  A%d,%d 0 0,1 %d,%d z", cx, cy, bx, by, chart.Radius, chart.Radius, endx, endy)
		canvas.Path(path, chart.FillStyles[i])
	}		

	// labels
	labels := len(chart.Labels)
	y := chart.GutterTop;
	// display bottom line labels
	for i := 0; i < labels; i++ {
		yoffset := int(float64(i) * 15)
		canvas.Text(chart.LegendXOffset+50, y+yoffset, chart.Labels[i], "font-size:75%;text-anchor:middle;")
		canvas.Rect(chart.LegendXOffset, y+yoffset-8, 30, 10, chart.FillStyles[i])
	}
	
	canvas.Gend()
	canvas.End()
	return nil
}


