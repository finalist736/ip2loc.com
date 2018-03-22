package logger

import (
	"github.com/gocraft/health"
	"os"
)

var jobsHealthStream *health.Stream

func JobsStream() *health.Stream {
	if jobsHealthStream == nil {
		jobsHealthStream = health.NewStream()
		jobsHealthStream.AddSink(&health.WriterSink{os.Stdout})
	}
	return jobsHealthStream
}
