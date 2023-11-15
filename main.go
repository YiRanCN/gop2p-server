package main

import (
	"net"
	"os"
	"strconv"
)

var clients = make(map[string]string)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":11194")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	listen, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println("udp server started")

	for {
		var buf [128]byte
		len, addr, err := listen.ReadFromUDPAddrPort(buf[:])

		if err != nil {
			println(err.Error())
			continue
		}

		ip := addr.Addr().String()
		port := addr.Port()
		msg := string(buf[:len])

		println(ip, ":", port, " -> ", msg)

		if msg[0:1] == "a" {
			handleA(msg[2:], ip, port)
		} else if msg[0:1] == "c" {
			handleC(listen, ip, port, msg[2:])
		}
	}

}

func handleA(cid string, ip string, port uint16) {
	client := ip + ":" + strconv.FormatUint(uint64(port), 10)
	clients[cid] = client
	println("clients append ", cid, ",", client)
}

func handleC(listen *net.UDPConn, myIp string, myPort uint16, otherCid string) {
	client, ok := clients[otherCid]
	if !ok {
		return
	}
	println("server -> ", client)

	addr, err := net.ResolveUDPAddr("udp4", client)
	if err != nil {
		println(err.Error())
		return
	}
	msg := "cc," + myIp + ":" + strconv.FormatUint(uint64(myPort), 10)
	_, err = listen.WriteToUDP([]byte(msg), addr)
	if err != nil {
		println(err.Error())
	}
	println("msg -> ", msg)
}
