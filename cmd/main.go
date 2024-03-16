package main

import (
	"flag"
	"goss-thrd-4/shrimp"
	"log"
	"time"
)

func main() {
	config := flag.String("config", "cfg/config.yaml", "path to config file")
	poll := flag.Int("poll", 24, "path to config file")
	flag.Parse()

	shr, err := shrimp.FromFile(*config)
	if err != nil {
		log.Fatal(err)
	}

	for ; true; _ = <-time.Tick(time.Hour * time.Duration(*poll)) {
		start := time.Now()
		if err := shr.Loop(); err != nil {
			log.Println(err)
		}
		log.Print("in", time.Since(start).String())
	}
}
