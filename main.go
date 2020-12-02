package main

import (
	"os"
	"time"
	control "yeelight/control"
	discovery "yeelight/discovery"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func handler(bulbs *control.Bulb) {
	log.Info().Msgf("%+v", bulbs)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Looking for bulbs")

	l := discovery.Listener{Interface: "wlp2s0"}
	if err := l.Listen(); err != nil {
		log.Fatal().Err(err)
	}
	defer l.Close()

	// bulbsChan := make(chan *yeelight.Bulb)
	go func() {
		for {
			bulb, err := l.Scan()
			if err != nil {
				log.Error().Err(err)
			}
			log.Info().Msgf("%+v", bulb)
		}
	}()

	outboundAddr, err := discovery.LookupBulbs()
	if err != nil {
		log.Fatal().Err(err)
	}

	bulbs, err := discovery.WaitBulbs(outboundAddr)
	for _, bulb := range bulbs {
		log.Printf("Found bulb: %+v", bulb)
		err := bulb.Connect()
		if err != nil {
			log.Fatal().Err(err)
		}
		go func() {
			for {
				result, err := bulb.ScanEvents()
				if err != nil {
					log.Fatal().Err(err)
				}
				log.Printf("%+v", result)
			}
		}()
		for i := 0; i < 10; i++ {
			if err := bulb.Toggle(i); err != nil {
				log.Fatal().Err(err)
			}
			time.Sleep(800 * time.Millisecond)
		}
		if err := bulb.Disconnect(); err != nil {
			log.Fatal().Err(err)
		}
	}
}
