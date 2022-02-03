package backend

import (
	"fmt"
	"net/http"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// out, err := exec.Command("go", "env").CombinedOutput()
	// if err != nil {
	// 	fmt.Println("16", err)
	// }
	// fmt.Println("18#", string(out))

	// out, err = exec.Command("echo", os.Getenv("PATH")).CombinedOutput()
	// if err != nil {
	// 	fmt.Println("22", err)
	// }
	// fmt.Println("24#", string(out))

	// out, err = exec.Command("git", "version").CombinedOutput()
	// if err != nil {
	// 	fmt.Println("28", err)
	// }
	// fmt.Println("30#", string(out))

	fmt.Fprintln(w, "Hello from test")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
