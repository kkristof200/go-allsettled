package allSettled

type FailedResult struct {
	Result
	Reason error
}
