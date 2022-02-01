package backend

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/nahidhasan98/nlogger"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// path, err := exec.LookPath("go")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Fprintln(w, path)

	logger := nlogger.NewLogger()
	logger.Warn("test 23#", time.Now())

	out, err := exec.Command("systemctl", "restart", "ajudge.service").Output()
	fmt.Println("test 25#", string(out), err)
	if err != nil {
		logger.Warn("test 28#"+err.Error(), time.Now())
		//fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, "Hello from test")
	// fmt.Fprintln(w, string(out))

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
