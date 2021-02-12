package errorhandling

//Check function for checking error
func Check(err error) {
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
}
