package deenz

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](values *T, err error) *T {
	if err != nil {
		HandleError(err)
	}
	return values
}

func CatchError(f func(err error)) {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			f(e)
		} else {
			panic(err)
		}
	}
}
