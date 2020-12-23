package task

type NoOpJob struct{}

func (job *NoOpJob) Process() error {
	return nil
}

func (job *NoOpJob) ID() string {
	return "noop"
}
