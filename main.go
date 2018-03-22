package main

import (
	"flag"
	"fmt"
	"github.com/finalist736/ip2loc.com/config"
	"github.com/finalist736/ip2loc.com/ctx"
	"github.com/finalist736/ip2loc.com/ip2loc"
	"github.com/finalist736/ip2loc.com/logger"
	"github.com/gocraft/web"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"github.com/finalist736/ip2loc.com/index"
	"github.com/finalist736/ip2loc.com/temple"
)

var config_path = flag.String("config", "config.ini", "config file location")
var listen_file_descriptor *int = flag.Int("fd", 0, "Server socket descriptor")
var listener1 net.Listener
var file1 *os.File = nil

func main() {

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	err := config.Init(config.NewFileProvider(config_path))
	if err != nil {
		panic(err)
	}
	logger.ReloadLogs()

	logger.StdOut().Print("-------------- SERVER STARTING ----------------------")

	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.StdOut().Errorf("Error Getting Rlimit: %v ", err)
	}
	logger.StdOut().Warnf("current RLIMIT: %+v", rLimit)

	rLimit.Max = 1048000
	rLimit.Cur = rLimit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.StdOut().Errorf("Error Setting Rlimit: %v", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		logger.StdOut().Errorf("Error Getting Rlimit: %v", err)
	}
	logger.StdOut().Warnf("final RLIMIT: %+v", rLimit)

	temple.Init()

	ip2loc.Init()
	defer ip2loc.Close()

	router := web.New(ctx.Ctx{})
	router.Get("/", index.Index)

	serv := &http.Server{Addr: config.MustString("listen"),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serv.SetKeepAlivesEnabled(false)

	logger.StdOut().Infof("Listen on %s", serv.Addr)
	if *listen_file_descriptor != 0 {
		logger.StdOut().Info("Starting with FD ", *listen_file_descriptor)
		file1 = os.NewFile(uintptr(*listen_file_descriptor), config.MustString("socket_path"))
		listener1, err = net.FileListener(file1)
		if err != nil {
			panic(fmt.Sprintf("fd listener failed: %s", err))
		}
	} else {
		logger.StdOut().Info("Virgin Start")
		listener1, err = net.Listen("tcp", serv.Addr)
		if err != nil {
			panic(fmt.Sprintf("listener failed: %s", err))
		}
	}
	go serv.Serve(listener1)

	usr1sigChannel := make(chan os.Signal)
	usr2sigChannel := make(chan os.Signal)
	interruptChannel := make(chan os.Signal)
	signal.Notify(usr1sigChannel, syscall.SIGUSR1)
	signal.Notify(usr2sigChannel, syscall.SIGUSR2)
	signal.Notify(interruptChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case killSignal := <-interruptChannel:
			logger.StdOut().Debug("Got signal:", killSignal)
			logger.StdOut().Debug("Stoping listening on ", serv.Addr)
			listener1.Close()
			if file1 != nil {
				file1.Close()
			}
			if killSignal == os.Interrupt {
				logger.StdOut().Warn("Daemon was interrupted by system signal")
				return
			}
			logger.StdOut().Debug("Daemon was killed")
			return
		case <-usr1sigChannel:
			logger.ReloadLogs()
		case usr2Signal := <-usr2sigChannel:
			logger.StdOut().Debug("Got signal:", usr2Signal)
			logger.StdOut().Debug("Grace restarting")
			listener2 := listener1.(*net.TCPListener)
			file2, err := listener2.File()
			if err != nil {
				logger.StdErr().Warn("get file2 from listener", err)
			}
			fd1 := int(file2.Fd())
			fd2, err := syscall.Dup(fd1)
			if err != nil {
				logger.StdErr().Warn("Dup error:", err)
			}
			listener1.Close()
			if file1 != nil {
				file1.Close()
			}
			command_line := fmt.Sprintf("%s", os.Args[0])
			fd_param := fmt.Sprintf("%d", fd2)
			cmd := exec.Command(command_line, "-config", *config_path, "-fd", fd_param)
			err = cmd.Start()
			if err != nil {
				panic("subprocess run error: " + err.Error())
			}
			time.Sleep(time.Second * 10)
			return
		}
		runtime.Gosched()
	}

}
