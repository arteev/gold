package main

import (
	"log"

	_ "github.com/arteev/gold"
	"github.com/arteev/gold/display"
	"github.com/arteev/gold/firich"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	config := map[string]interface{}{
		"Name": "/dev/ttyUSB0",
		"Baud": 9600,
	}
	linedsp, err := display.Open("serial", config)
	if err != nil {
		log.Fatal(err)
	}
	dsp, err := linedsp.GetDisplay(&firich.FirichProtocol{})
	if err != nil {
		log.Fatal(err)
	}
	defer dsp.Close()
	dsp.SetEncoding(charmap.CodePage866)
	if err := dsp.Init(); err != nil {
		log.Fatal(err)
	}
	dsp.PrintRow(1, "Price:10$ Quant:2")
	dsp.PrintRow(2, "Total:20$")
}
