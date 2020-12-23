package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/bop0hz/go-yeelight/control"
	"github.com/bop0hz/go-yeelight/discovery"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func scanBulbs(l *discovery.Listener, bulbC chan *control.Bulb) {
	for {
		bulb, err := l.Scan()
		if err != nil {
			log.Error().Err(err)
		}
		if bulb != nil {
			bulbC <- bulb
		}
	}
}

func askBulbs(bulbC chan *control.Bulb) {
	outboundAddr, err := discovery.LookupBulbs()
	if err != nil {
		log.Fatal().Err(err)
	}
	bulbs, err := discovery.WaitBulbs(outboundAddr)
	for _, bulb := range bulbs {
		bulbC <- bulb
	}
}

func scanEvents(b *control.Bulb) {
	for {
		result, err := b.ScanEvents()
		if err != nil {
			log.Fatal().Err(err)
		}
		var r *control.Result
		log.Debug().Msgf("%+v", result)
		if err := json.Unmarshal([]byte(result), &r); err != nil {
			log.Fatal().Err(err)
		}
		log.Debug().Msgf("%+v", r)
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

	if len(os.Args) < 2 {
		log.Fatal().Msgf("Set network interface name as argument, e.g:\n%v eth0", os.Args[0])
	}
	l, err := discovery.NewListener(os.Args[1])
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create listener:")
	}
	if err := l.Listen(); err != nil {
		log.Fatal().Err(err).Msg("Could not listen:")
	}
	defer l.Close()

	bulbC := make(chan *control.Bulb)
	log.Info().Msg("Starting discovery in the background")
	go scanBulbs(l, bulbC)
	log.Info().Msg("Asking any bulbs online")
	go askBulbs(bulbC)
	for {
		select {
		case bulb := <-bulbC:
			log.Info().Msgf("Found a bulb: %+v", bulb)
			log.Debug().Msgf("Connecting to the bulb %s", bulb.Location[0])
			if err := bulb.Connect(); err != nil {
				log.Fatal().Err(err).Msgf("Could not connect to the bulb %+v", bulb.Location)
			}
			defer bulb.Disconnect()
			for i := 1; i <= 2; i++ {
				log.Debug().Msgf("Toggle the bulb %+v", bulb)
				if err := bulb.Toggle(i); err != nil {
					log.Fatal().Err(err).Msgf("Could not toggle the bulb %+v", bulb)
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}

}
