package main

import (
	"log"
	"time"
	yeelight "yeelight/discovery"
)

func main() {
	log.Println("Looking for bulbs")

	outboundAddr, err := yeelight.DiscoverBulbs()
	if err != nil {
		log.Fatal(err)
	}

	bulbs, err := yeelight.WaitBulbs(outboundAddr)
	for _, bulb := range bulbs {
		log.Printf("Found bulb: %+v", bulb)
		err := bulb.Connect()
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			result, err := bulb.Toggle(i)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%+v", result)
			time.Sleep(800 * time.Millisecond)
		}
		bulb.Disconnect()
	}

}
