package helpers

func Treat(err error) {
	if err != nil {
		panic(err)
	}
}
