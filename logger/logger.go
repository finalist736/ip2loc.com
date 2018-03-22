package logger

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
	"time"
	"os"
	"syscall"
	"io"
	"github.com/finalist736/ip2loc.com/config"
)

var stdout, stderr *logrus.Logger

func JsonStdOut(name string, jsn interface{}) {
	ba, _ := json.MarshalIndent(jsn, "", " ")
	StdOut().Printf(name+" %v", ba)
}

func StdOut() *logrus.Logger {
	return stdout
}

func StdErr() *logrus.Logger {
	return stderr
}

func ReloadLogs() {
	logType := config.MustString("logtype")
	switch logType {
	case "std":
		// selected std outs
	case "file":
		logPath := config.MustString("logpath")
		if logPath == "" {
			panic("logger: no log path specified")
		}
		err := os.MkdirAll(logPath, 0764)
		if err != nil {
			panic(err)
		}

		stdoutfile, err := os.OpenFile(logPath+"/stdout.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0764)
		if err != nil {
			panic(err)
		}
		_, err = stdoutfile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err)
		}
		err = syscall.Dup2(int(stdoutfile.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			panic(err)
		}

		stderrfile, err := os.OpenFile(logPath+"/stderr.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0764)
		if err != nil {
			panic(err)
		}
		_, err = stderrfile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err)
		}
		err = syscall.Dup2(int(stderrfile.Fd()), int(os.Stderr.Fd()))
		if err != nil {
			panic(err)
		}
	case "socket":
		logPath := config.MustString("logpath")
		if logPath == "" {
			panic("logger: no log socket specified")
		}
		// open socket, and send all logs there!
	default:
		panic("logger: incorrect logging type")
	}
	logLevel := config.MustString("loglevel")

	var err error
	stdout = logrus.New()
	stdout.Formatter = &logrus.TextFormatter{
		ForceColors:false,
		TimestampFormat:time.RFC3339Nano,
		DisableTimestamp:false,
		FullTimestamp:true,
	}
	stdout.Level, err = logrus.ParseLevel(logLevel)
	if err != nil {
		stdout.Level = logrus.InfoLevel
	}
	stdout.Out = os.Stdout

	stderr = logrus.New()
	stderr.Formatter = &logrus.TextFormatter{
		ForceColors:false,
		TimestampFormat:time.RFC3339Nano,
		DisableTimestamp:false,
		FullTimestamp:true,
	}
	stderr.Level, err = logrus.ParseLevel(logLevel)
	if err != nil {
		stderr.Level = logrus.InfoLevel
	}
	stderr.Out = os.Stderr
}

