package main

import (
	"fmt"
	"os"

	"github.com/BertoldVdb/ms-tools/gohid"
	"github.com/BertoldVdb/ms-tools/mshal"
	"github.com/alecthomas/kong"
)

type Context struct {
	dev gohid.HIDDevice
	hal *mshal.HAL
}

var CLI struct {
	VID      int    `optional type:"hex" help:"The USB Vendor ID." default:534d`
	VID2     int    `optional type:"hex" help:"The second USB Vendor ID." default:345f`
	PID      int    `optional type:"hex" help:"The USB Product ID."`
	Serial   string `optional help:"The USB Serial."`
	RawPath  string `optional help:"The USB Device Path."`

	ListDev  ListHIDCmd `cmd help:"List devices."`

	Status   StatusCmd `cmd help:"Print status." default:1`
}

func main() {
	k, err := kong.New(&CLI,
		kong.NamedMapper("int", intMapper{}),
		kong.NamedMapper("hex", intMapper{base: 16}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := k.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := &Context{}
	if ctx.Command() != "list-dev" {
		dev, err := OpenDevice()
		if err != nil {
			fmt.Println("Failed to open device", err)
			os.Exit(1)
		}
		defer dev.Close()

		c.dev = dev
		config := mshal.HALConfig{
			// auto-assume EEPROM in mshal
			PatchProbeEEPROM: true,

			LogFunc: func(level int, format string, param ...interface{}) {
				//str := fmt.Sprintf(format, param...)
				//fmt.Printf("HAL(%d): %s\n", level, str)
			},
		}

		c.hal, err = mshal.New(dev, config)
		if err != nil {
			fmt.Println("Failed to create HAL", err)
			os.Exit(1)
		}
	}

	err = ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
