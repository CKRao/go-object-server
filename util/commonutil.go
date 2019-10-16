package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var lock sync.Mutex

var Ip string = GetIP()

func GetIP() string {
	lock.Lock()
	defer lock.Unlock()

	if Ip == "" {
		Ip = getIp()
	}

	return Ip
}

//获取外网IP
func getIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		log.Fatalln("getIp error : ", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	ip := string(data)

	return ip
}
