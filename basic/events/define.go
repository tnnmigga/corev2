package events

type SystemEvent string

func (e SystemEvent) Name() string {
	return string(e)
}

const (
	SystemAfterRun   SystemEvent = "system-after-run"
	SystemBeforeStop SystemEvent = "system-before-stop"
)
