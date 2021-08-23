package errorhandling

import (
	"fmt"
	"time"

	"github.com/nahidhasan98/nlogger"
)

var loger = nlogger.NewLogger()

//Check function for checking error
func Check(err error) {
	if err != nil {
		fmt.Println(err)
		//panic(err)

		loger.Error(err, time.Now())
	}
}
