package discovery

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/bop0hz/go-yeelight/control"
)

const (
	mcastAddr string = "239.255.255.250"
	mcastPort int    = 1982
	searchMsg string = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb\r\n"
)

// Listener listens for bulbs via multicast
type Listener struct {
	conn  *net.UDPConn
	addrs []net.Addr
	netIf *net.Interface
	mAddr *net.UDPAddr
}

// NewListener creates and initializes Listener on network interface
func NewListener(netInterface string) (l *Listener, err error) {
	ifi, err := net.InterfaceByName(netInterface)
	if err != nil {
		return
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		return
	}

	ip := net.ParseIP(mcastAddr)
	addr := &net.UDPAddr{IP: ip, Port: mcastPort}
	return &Listener{netIf: ifi, mAddr: addr, addrs: addrs}, nil
}

// Listen starts listen multicast
func (l *Listener) Listen() (err error) {
	l.conn, err = net.ListenMulticastUDP("udp", l.netIf, l.mAddr)
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
func (l *Listener) Scan() (bulb *control.Bulb, err error) {
	buffer := make([]byte, 1000)
	n, lAddr, err := l.conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	for _, addr := range l.addrs {
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
	return control.UnmarshalBulb(&req.Header)
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

	_, err = c.Write([]byte(searchMsg))
	if err != nil {
		return
	}
	return net.ResolveUDPAddr("udp", c.LocalAddr().String())
}

// WaitBulbs waits bulbs responses and return found bulbs
func WaitBulbs(addr *net.UDPAddr) (bulbs []*control.Bulb, err error) {
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

	bulb, err := control.UnmarshalBulb(&resp.Header)
	if err != nil {
		return
	}
	bulbs = append(bulbs, bulb)
	return
}
