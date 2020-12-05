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

type Bulb struct {
	Bright   []string
	Location []string
	Addr     string
	Support  []string
	Name     []string
	channel  net.Conn
}

type Error struct {
	Code    int
	Message string
}

type Result struct {
	ID     int
	Result []string
	Error  Error
}

type Notification struct {
	Method string
	Params map[string]string
}

func UnmarshalBulb(resp *http.Header) (bulb *Bulb, err error) {
	jsonData, err := json.Marshal(resp)
	if err = json.Unmarshal(jsonData, &bulb); err != nil {
		return
	}
	bulb.Addr = strings.TrimPrefix(string(bulb.Location[0]), "yeelight://")
	return
}

func (b *Bulb) Connect() (err error) {
	b.channel, err = net.Dial("tcp", b.Addr)
	if err != nil {
		return
	}
	return
}

func (b *Bulb) Disconnect() (err error) {
	return b.channel.Close()
}

func (b *Bulb) Toggle(id int) (err error) {
	_, err = fmt.Fprintf(b.channel, "{\"id\":%v,\"method\":\"toggle\",\"params\":[]}\r\n", id)
	if err != nil {
		return
	}
	return
}

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
