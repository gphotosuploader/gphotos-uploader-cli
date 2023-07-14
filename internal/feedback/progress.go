package feedback

import "github.com/schollz/progressbar/v3"

// Progress is a progress bar to show the status of a task.
type Progress struct {
	bar *progressbar.ProgressBar
}

// NewTaskProgressBar returns a progress bar for tasks that outputs to
// the terminal.
func NewTaskProgressBar(desc string, steps int, visibility bool) *Progress {
	bar := progressbar.NewOptions(steps,
		progressbar.OptionSetWriter(feedbackOut),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetVisibility(visibility),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	return &Progress{
		bar: bar,
	}
}

func (pb *Progress) Add(num int) {
	_ = pb.bar.Add(num)
}

func (pb *Progress) Finish() {
	_ = pb.bar.Finish()
}
