package main

import (
	"fmt"
	"os"
	"time"

	"encoding/json"

	"github.com/google/renameio"
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

		buf := readmem(c, Region { Region: "RAM", Addr: 0x0000e180 }, 16)

		if (len(buf) == 0) {
			// failed read in a loop should just continue and hope...
			if (s.Loop == 0) {
				fmt.Println("Read nothing from RAM, exiting")
				os.Exit(1)
			} else {
				time.Sleep(time.Duration(s.Loop) * time.Millisecond)
				continue
			}
		}

		// signal is weird and seems to be different depending on (non-)interlaced input
		// observed:
		// - no signal: 0x00
		// - no signal after non-interlaced: 0x08
		// - non-interlaced signal: 0x07
		// - interlaced signal: 0x0f
		signal := int(buf[0:1][0])

		output.Width = (int(buf[5]) * 256) + int(buf[4])
		output.Height = (int(buf[13]) * 256) + int(buf[12])

		// looks like 0x0f/15 means we have an interlaced signal, so we read half the height
		if (signal == 15) {
			output.Height = 2 * output.Height
		}

		if (signal > 0 && signal != 8) {
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

		if (s.Filename == "") {
			fmt.Print(p)
		} else {
			renameio.WriteFile(s.Filename, []byte(p), os.FileMode(int(0644)))
		}

		if (s.Loop == 0) {
			next = false
		} else {
			time.Sleep(time.Duration(s.Loop) * time.Millisecond)
		}
	}

	return nil
}
