package discovery

import (
	"encoding/json"
	"net/http"
)

const (
	SearchMsg string = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb\r\n"
)

type Notify struct{}

func Unmarshal(h *http.Header, msg *Notify) (err error) {
	jsonData, err := json.Marshal(h)
	return json.Unmarshal(jsonData, &msg)

}
