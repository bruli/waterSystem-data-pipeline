package executedlogs

import "time"

type ExecutedLog struct {
	zone       string
	seconds    int
	executedAt time.Time
}

func (e ExecutedLog) Zone() string {
	return e.zone
}

func (e ExecutedLog) Seconds() int {
	return e.seconds
}

func (e ExecutedLog) ExecutedAt() time.Time {
	return e.executedAt
}

func New(zone string, seconds int, executedAt time.Time) *ExecutedLog {
	return &ExecutedLog{zone: zone, seconds: seconds, executedAt: executedAt}
}
