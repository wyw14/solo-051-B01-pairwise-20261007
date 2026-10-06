package run

import "github.com/wyw14/cry-161/internal/model"

type Result struct {
	Run      *model.Run
	Released bool
	Message  string
}

func BuildResult(run *model.Run, released bool) Result {
	message := "stopped cleanly"
	if run.Error != "" {
		message = run.Error
	}
	return Result{Run: run, Released: released, Message: message}
}
