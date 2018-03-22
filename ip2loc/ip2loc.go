package ip2loc

import (
	"github.com/finalist736/ip2location-go"
	"github.com/gocraft/health"
	"github.com/finalist736/ip2loc.com/config"
	"github.com/finalist736/ip2loc.com/logger"
)

func Init() {
	file := config.MustString("ip2location")
	logger.StdOut().Debugf("opening ip2location database: %v", file)
	var err = ip2location.Open(file)
	if err != nil {
		panic(err)
	}
}

func Close() {
	ip2location.Close()
}

func Get(ip string) *ip2location.IP2Locationrecord {
	job := logger.JobsStream().NewJob("ip2loc")
	job.EventKv("ip", health.Kvs{"addr": ip})
	defer job.Complete(health.Success)
	record := ip2location.Get_all(ip)
	return &record
}
