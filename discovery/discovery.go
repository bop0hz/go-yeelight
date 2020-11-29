package yeelight

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"
	yeelight "yeelight/control"
)

const (
	searchMsg string = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb\r\n"
	mcastAddr string = "239.255.255.250:1982"
)

func Discover() (err error) {
	ip := net.ParseIP(mcastAddr)
	ifi, _ := net.InterfaceByName("wlp2s0")

	addr := &net.UDPAddr{IP: ip, Port: 1982}
	conn, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024)
	defer conn.Close()
	_, a, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s, %v", buffer, a)
	return
}

func DiscoverBulbs() (addr *net.UDPAddr, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", mcastAddr)
	if err != nil {
		return
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	defer c.Close()

	_, err = c.Write([]byte(searchMsg))
	if err != nil {
		return
	}

	return net.ResolveUDPAddr("udp", c.LocalAddr().String())

}

func WaitBulbs(addr *net.UDPAddr) (bulbs []*yeelight.Bulb, err error) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1000)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	responseString := string(buffer[:n])
	reader := bufio.NewReader(strings.NewReader(responseString + "\r\n"))
	logResp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return
	}

	// var bulb yeelight.Bulb
	bulb, err := yeelight.Init(&logResp.Header)
	if err != nil {
		return
	}
	bulbs = append(bulbs, bulb)
	return
}
