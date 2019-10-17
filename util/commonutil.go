package util

import (
	"fmt"
	"go-object-server/commons"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var lock sync.Mutex

var ip string

func GetIP() string {
	lock.Lock()
	defer lock.Unlock()

	if ip == "" {
		ip = getIp()
	}

	return ip
}

func GetAddress() string {
	port := commons.GetConfigIns().Server.Port
	address := fmt.Sprintf("%s:%s", GetIP(), port)

	return address
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
