package main

import (
	"fmt"
	"encoding/binary"
	"encoding/json"
)

type Region struct {
	Region string `arg name:"region" help:"Memory region to access."`
	Addr   int    `arg name:"addr" help:"Addresses to access." type:"int"`
}

type StatusCmd struct {
	Json   bool   `optional name:"json" help:"Output JSON instead of the default list."`
}

type outputData struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Signal string `json:"signal"`
}

func (s *StatusCmd) Run(c *Context) error {
	output := &outputData{}

	buf := readmem(c, Region { Region: "RAM", Addr: 0x0000e184 }, 2)
	output.Width = int(binary.LittleEndian.Uint16(buf))

	buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e18c }, 2)
	output.Height = int(binary.LittleEndian.Uint16(buf))

	buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e180 }, 1)

	if (len(buf) > 0 && buf[0] > 0) {
		output.Signal = "yes"
	} else {
		output.Signal = "no"
	}

	if (s.Json) {
		data, _ := json.Marshal(output)
		fmt.Println(string(data))
	} else {
		fmt.Printf("width: %d\n", output.Width)
		fmt.Printf("height: %d\n", output.Height)
		fmt.Printf("signal: %s\n", output.Signal)
	}

	return nil
}
