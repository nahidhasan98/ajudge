package backend

import (
	"fmt"
	"net/http"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, "Hello from test update")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
