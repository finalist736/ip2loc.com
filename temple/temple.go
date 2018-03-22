package temple

import (
	"html/template"
	"sync"
	"github.com/finalist736/ip2loc.com/config"
)


var mux sync.RWMutex
var storage *template.Template

var templatesList = []string{"index.html"}


func Init() {
	mux.Lock()
	defer mux.Unlock()

	var err error
	path := config.MustString("templates")
	fullPath := make([]string, 0)
	for _, iterator := range templatesList {
		fullPath = append(fullPath, path + iterator)
	}

	storage, err = template.New("index.html").ParseFiles(fullPath...)

	//storage, err = template.ParseFiles(fullPath...)
	if err != nil {
		panic(err)
	}
}

func Get() *template.Template {
	mux.RLock()
	defer mux.RUnlock()
	return storage
}
