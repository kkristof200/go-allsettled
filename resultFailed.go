package allSettled

type SuccessfulResult[T any] struct {
	Result
	Value *T
}
