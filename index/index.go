package index

import (
	"fmt"
	"github.com/finalist736/ip2loc.com/ctx"
	"github.com/finalist736/ip2loc.com/logger"
	"github.com/gocraft/web"
	"net"
	"github.com/finalist736/ip2loc.com/ip2loc"
	"github.com/finalist736/ip2loc.com/temple"
	"github.com/finalist736/ip2location-go"
	"sync"
	"github.com/finalist736/ip2loc.com/config"
	"strings"
)

type IpData struct {
	*ip2location.IP2Locationrecord
	IP string
	GoogleApiKey string
}

var ipdataPool = sync.Pool{New: func() interface{} {
	return &IpData{GoogleApiKey: config.MustString("apikey")}
}}

func Index(c *ctx.Ctx, rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()
	ipdata := ipdataPool.Get().(*IpData)
	defer ipdataPool.Put(ipdata)
	ipdata.IP = ""
	ipdata.IP = req.FormValue("ip")
	logger.StdOut().Debugf("ip str: %+v", ipdata.IP)
	if ipdata.IP == "" {
		ipdata.IP = req.Header.Get("X-Real-IP")
		logger.StdOut().Debugf("ip real: %+v", ipdata.IP)
		if ipdata.IP == "" {
			ipdata.IP = req.RemoteAddr
			logger.StdOut().Debugf("ip remote: %+v", ipdata.IP)
		}
	}
	if strings.Contains(ipdata.IP, ":") {
		ipdata.IP = strings.Split(ipdata.IP, ":")[0]
		logger.StdOut().Debugf("ip clerified: %+v", ipdata.IP)
	}
	ip := net.ParseIP(ipdata.IP)
	//if ip == nil {
	//	logger.StdErr().Infof("no correct ip address")
	//	fmt.Fprint(rw, "no ip address :(")
	//	return
	//}
	ipdata.IP2Locationrecord = ip2loc.Get(ip.String())

	logger.StdOut().Debugf("ip data: %+v", ipdata)

	t := temple.Get()
	err := t.Execute(rw, ipdata)
	if err != nil {
		logger.StdErr().Debugf("tpl error: %v", err)
		fmt.Fprint(rw, "internal error with templates :(")
	}
}
