// ViChart library for Go
// Author: Tad Vizbaras 
// License: http://github.com/tadvi/vichart/blob/master/LICENSE 
//
package vichart

// Must - checks if there are errors and panics if there are any.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
