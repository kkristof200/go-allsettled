package allSettled

type AllSettledConfig[T any] struct {
	OnSuccess      func(int, SuccessfulResult[T])
	OnFailure      func(int, FailedResult)
	MaxConcurrency int
}
