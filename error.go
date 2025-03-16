package deenz

func Error(err error) {
	panic(err)
}

func Must[T any](values *T, err error) *T {
	if err != nil {
		Error(err)
	}
	return values
}
