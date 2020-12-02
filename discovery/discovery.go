package discovery

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"strings"
	yeelight "yeelight/control"
)

const (
	mcastAddr string = "239.255.255.250"
	mcastPort int    = 1982
)

type Listener struct {
	conn      *net.UDPConn
	Addrs     []net.Addr
	Interface string
}

// Listen init multicast listener
func (l *Listener) Listen() (err error) {
	ip := net.ParseIP(mcastAddr)
	ifi, err := net.InterfaceByName(l.Interface)
	if err != nil {
		return
	}
	l.Addrs, err = ifi.Addrs()
	if err != nil {
		return
	}
	addr := &net.UDPAddr{IP: ip, Port: mcastPort}
	l.conn, err = net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return
	}
	return
}

// Close closes listener connection
func (l *Listener) Close() (err error) {
	return l.conn.Close()
}

// Scan scans for bulbs going online
func (l *Listener) Scan() (bulb *yeelight.Bulb, err error) {
	buffer := make([]byte, 1000)
	n, lAddr, err := l.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	for _, addr := range l.Addrs {
		if strings.Contains(addr.String(), lAddr.IP.String()) {
			return
		}
	}

	msgString := string(buffer[:n])
	reader := bufio.NewReader(strings.NewReader(msgString + "\r\n"))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return
	}
	return yeelight.UnmarshalBulb(&req.Header)
}

// LookupBulbs sends request if there are any bulbs online
func LookupBulbs() (addr *net.UDPAddr, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", mcastAddr+":"+strconv.Itoa(mcastPort))
	if err != nil {
		return
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	defer c.Close()

	_, err = c.Write([]byte(SearchMsg))
	if err != nil {
		return
	}
	return net.ResolveUDPAddr("udp", c.LocalAddr().String())
}

// WaitBulbs waits bulbs responses and return found bulbs
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
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return
	}

	bulb, err := yeelight.UnmarshalBulb(&resp.Header)
	if err != nil {
		return
	}
	bulbs = append(bulbs, bulb)
	return
}
