package control

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type resultError struct {
	Code    int
	Message string
}

// Result is a response following by the command to the bulb
type Result struct {
	ID     int
	Result []string
	Error  resultError
}

// Notification from the bulb once state is changed
type Notification struct {
	Method string
	Params map[string]string
}

// Bulb struct
type Bulb struct {
	Bright   []string
	Location []string
	Addr     string
	Support  []string
	Name     []string
	channel  net.Conn
}

// UnmarshalBulb unmarshals http.Header to Bulb
func UnmarshalBulb(resp *http.Header) (bulb *Bulb, err error) {
	jsonData, err := json.Marshal(resp)
	if err = json.Unmarshal(jsonData, &bulb); err != nil {
		return
	}
	bulb.Addr = strings.TrimPrefix(string(bulb.Location[0]), "yeelight://")
	return
}

// NewBulb creates bulb
func NewBulb(addr string) (bulb *Bulb, err error) {
	return &Bulb{Addr: addr}, nil
}

// Connect connects to the bulb
func (b *Bulb) Connect() (err error) {
	b.channel, err = net.Dial("tcp", b.Addr)
	if err != nil {
		return
	}
	return
}

// Disconnect from the bulb
func (b *Bulb) Disconnect() (err error) {
	return b.channel.Close()
}

// Toggle bulb state
func (b *Bulb) Toggle(id int) (err error) {
	_, err = fmt.Fprintf(b.channel, "{\"id\":%v,\"method\":\"toggle\",\"params\":[]}\r\n", id)
	if err != nil {
		return
	}
	return
}

// SetName set bulb name
func (b *Bulb) SetName(id int, n string) (err error) {
	_, err = fmt.Fprintf(b.channel, "{\"id\":%v,\"method\":\"set_name\",\"params\":[\"%v\"]}\r\n", id, n)
	if err != nil {
		return
	}
	return
}

// SetBright set bulb brighteness
func (b *Bulb) SetBright(id int, perc uint8, msec int) (err error) {
	_, err = fmt.Fprintf(b.channel, "{\"id\":%v,\"method\":\"set_bright\",\"params\":[%v, \"smooth\", %v]}\r\n", id, perc, msec)
	if err != nil {
		return
	}
	return
}

// ScanEvents scans events from the bulb
func (b *Bulb) ScanEvents() (event []byte, err error) {
	buffer := make([]byte, 1000)
	n, err := b.channel.Read(buffer)
	if err != nil {
		return
	}
	s := bufio.NewScanner(bytes.NewReader(buffer[:n]))
	s.Scan()
	event = s.Bytes()
	return
}
