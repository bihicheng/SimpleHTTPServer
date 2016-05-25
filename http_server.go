package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Usage(msg ...string) {
	for _, m := range msg {
		if len(m) > 0 {
			fmt.Println(m)
		}
	}
	if len(msg) == 0 {
		fmt.Printf("Usage: SimpleHTTPServer [port|host:port] (default: 127.0.0.1:8000)\n")
	}
	os.Exit(0)
}

func Dump(params interface{}, exit ...bool) {
	atype := reflect.TypeOf(params)
	fmt.Printf("(%+v)%+v\n", atype, params)
	for _, exit := range exit {
		if exit == true {
			os.Exit(1)
		}
	}
}

func IsIp(ip string) bool {
	reg, err := regexp.Compile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	if err != nil {
		log.Fatal(err)
	}
	var ips []byte
	ips = []byte(ip)
	ret := reg.Match(ips)
	return ret
}

func IsNumberic(param string) (bool, int) {
	if a, err := strconv.Atoi(param); err == nil {
		return true, a
	}
	return false, 0
}

func main() {
	var port string = "8000"
	var host string = "127.0.0.1"
	if len(os.Args) == 2 {
		port = os.Args[1]
		host_port_func := func() (string, string) {
			var host, port string = host, port
			host_port := strings.Split(port, ":")
			if len(host_port) != 2 {
				Usage()
			} else {
				host = host_port[0]
				port = host_port[1]
				if !IsIp(host) {
					Usage("Error: invalidate host " + host)
				}
			}
			return host, port
		}

		b, _ := IsNumberic(port)
		if b == false {
			host, port = host_port_func()
			//Dump(host, true)
		}
	}
	host = host + ":" + port
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	log.Println("Listensing...")
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
