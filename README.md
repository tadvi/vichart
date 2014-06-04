# vichart: Go library for SVG charts

vichart (a.k.a Vivid Chart) is simple chart toolkit designed as wrap around SVGo library.
Each chart is more or less independent of other charts and might be used as stand-alone piece.

Toolkit is for Go (golang). It generates SVG as defined by the Scalable Vector Graphics 1.1 Specification
and requires SVGo library to produce charts. Get it here
[http://github.com/ajstarks/svgo](http://github.com/ajstarks/svgo)

Some charts try to display independent information in the single screen space. Examples are
vertical bar charts. They show both bar and line information that does not have to be related.

## Current State

Library is in Alpha state. It is not complete and being worked on. Consider it v0.

## Package Contents

* webcharts: basic webserver to showcase charts
* vichart: charts library

## Try

Try vichart library:

1. Compile it using *go build*
2. Run *webcharts* executable
3. Point your browser to http://localhost:8080/vchart

## Example

There is example of vertical multibar chart. It displays 3 section bars with additional line.
Scale values for bars are on the left while line values are on the right.

![Vertical Multibar Chart](/screen.png)
