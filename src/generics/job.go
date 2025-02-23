package generics

func JobThreeArgs[T any, K any, M, O any](fn func(T, K, M) O, args ...any) func() O {
	return func() O {
		return fn(args[0].(T), args[1].(K), args[2].(M))
	}
}

func JobTwoArgs[T any, K any, O any](fn func(T, K) O, args ...any) func() O {
	return func() O {
		return fn(args[0].(T), args[1].(K))
	}
}

func JobOneArg[T any, O any](fn func(T) O, args ...any) func() O {
	return func() O {
		return fn(args[0].(T))
	}
}

func JobInfiniteArgs(fn func(args ...any) error, args ...any) func() error {
	return func() error {
		return fn(args...)
	}
}
