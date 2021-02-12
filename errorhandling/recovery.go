package errorhandling

import "fmt"

//Recovery function for recovering from panic
func Recovery() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}
