package yeelight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func Init(resp *http.Header) (bulb *Bulb, err error) {
	jsonData, err := json.Marshal(resp)
	err = json.Unmarshal(jsonData, &bulb)
	if err != nil {
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

func (b *Bulb) Toggle(id int) (res Result, err error) {
	_, err = fmt.Fprintf(b.channel, "{\"id\":%v,\"method\":\"toggle\",\"params\":[]}\r\n", id)
	if err != nil {
		return
	}
	buffer := make([]byte, 1000)

	n, _ := b.channel.Read(buffer)

	s := bufio.NewScanner(bytes.NewReader(buffer[:n]))
	for s.Scan() {
		err = json.Unmarshal(buffer[:n], &res)
		log.Printf("%s", buffer[:n])
		if err != nil {
			return
		}
	}

	return
}
