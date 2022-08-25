package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

var (
	port       = 7160
	host       = "127.0.0.1"
	disabledIp = "localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;172.32.*;192.168.*"
)

func setProxy(enable bool) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		log.Fatal("open reg:", err)
	}
	defer k.Close()

	var enabled uint32 = 0
	var server = ""
	var ignoreIp = ""

	if enable {
		enabled = 1
		server = fmt.Sprintf("%s:%d", host, port)
		ignoreIp = disabledIp
	}

	if err := k.SetDWordValue("ProxyEnable", enabled); err != nil {
		log.Fatal("set proxy:", err)
	}
	if err := k.SetStringValue("ProxyServer", server); err != nil {
		log.Fatal("set proxy:", err)
	}
	if err := k.SetStringValue("ProxyOverride", ignoreIp); err != nil {
		log.Fatal("set proxy:", err)
	}
}
