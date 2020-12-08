package main

import (
	"encoding/json"
	"os"
	control "yeelight/control"
	discovery "yeelight/discovery"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func handler(bulbs *control.Bulb) {
	log.Info().Msgf("%+v", bulbs)
}

func scanBulbs(l discovery.Listener) {
	for {
		bulb, err := l.Scan()
		if err != nil {
			log.Error().Err(err)
		}
		if bulb != nil {
			log.Info().Msgf("Scan found a bulb: %+v", bulb)
		}
	}
}

func scanEvents(b *control.Bulb) {
	for {
		result, err := b.ScanEvents()
		if err != nil {
			log.Fatal().Err(err)
		}
		var r *control.Result
		if err := json.Unmarshal([]byte(result), &r); err != nil {
			log.Fatal().Err(err)
		}
		if r.ID == 0 && r.Result == nil {
			var r *control.Notification
			json.Unmarshal([]byte(result), &r)
			if err := json.Unmarshal([]byte(result), &r); err != nil {
				log.Fatal().Err(err)
			}
			log.Printf("%+v, %T", r, r)
		} else {
			log.Printf("%+v, %T", r, r)
		}
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Starting")

	l := discovery.Listener{Interface: "wlp2s0"}
	if err := l.Listen(); err != nil {
		log.Fatal().Err(err)
	}
	defer l.Close()
	go scanBulbs(l)

	homeBulb, _ := control.NewBulb("192.168.1.53:55443")
	if err := homeBulb.Connect(); err != nil {
		log.Fatal().Err(err)
	}
	// go scanEvents(homeBulb)
	if err := homeBulb.SetBright(1, 1, 10000); err != nil {
		log.Fatal().Err(err)
	}
	// homeBulb.Toggle(1)
	defer homeBulb.Disconnect()

	outboundAddr, err := discovery.LookupBulbs()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("Asking any bulbs online...")
	bulbs, err := discovery.WaitBulbs(outboundAddr)
	for _, bulb := range bulbs {
		log.Info().Msgf("Found a bulb on discovery request: %+v", bulb)
		if err := bulb.Connect(); err != nil {
			log.Fatal().Err(err)
		}
	}

	// 	for i := 0; i < 6; i++ {
	// 		if err := bulb.Toggle(i); err != nil {
	// 			log.Fatal().Err(err)
	// 		}
	// 		time.Sleep(1000 * time.Millisecond)
	// 	}
	// 	defer bulb.Disconnect()
	// }
}
