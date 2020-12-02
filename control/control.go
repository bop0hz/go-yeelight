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

func (b *Bulb) ScanEvents() (res Result, err error) {
	buffer := make([]byte, 1000)
	n, _ := b.channel.Read(buffer)
	s := bufio.NewScanner(bytes.NewReader(buffer[:n]))
	for s.Scan() {
		err = json.Unmarshal(s.Bytes(), &res)
		if err != nil {
			return
		}
	}
	return
}
