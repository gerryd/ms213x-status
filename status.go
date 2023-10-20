package main

import (
	"fmt"
	"time"
	"encoding/binary"
	"encoding/json"
)

type Region struct {
	Region string `arg name:"region" help:"Memory region to access."`
	Addr   int    `arg name:"addr" help:"Addresses to access." type:"int"`
}

type StatusCmd struct {
	Json      bool   `optional name:"json" help:"Output JSON instead of the default list."`
	Loop      int    `optional name:"loop" help:"Run in a loop and sleep every N microseconds"`
	Filename  string `optional name:"filename" help:"Output to a file instead of stdout"`
}

type outputData struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Signal string `json:"signal"`
}

func (s *StatusCmd) Run(c *Context) error {
	next := true

	for next {
		output := &outputData{}

		buf := readmem(c, Region { Region: "RAM", Addr: 0x0000e184 }, 2)
		output.Width = int(binary.LittleEndian.Uint16(buf))

		buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e18c }, 2)
		output.Height = int(binary.LittleEndian.Uint16(buf))

		buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e180 }, 1)

		if (len(buf) == 0) {

		} else if (buf[0] > 0) {
			output.Signal = "yes"
		} else {
			output.Signal = "no"
		}

		var p = ""

		if (s.Json) {
			data, _ := json.Marshal(output)
			p = string(data)

		} else {
			p = fmt.Sprintf("width: %d\nheight: %d\nsignal: %s\n", output.Width, output.Height, output.Signal)
		}

		fmt.Print(p)

		if (s.Loop == 0) {
			next = false
		} else {
			time.Sleep(time.Duration(s.Loop) * time.Millisecond)
		}
	}

	return nil
}
