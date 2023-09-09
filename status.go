package main

import (
	"fmt"
	"encoding/binary"
)

type Region struct {
	Region string `arg name:"region" help:"Memory region to access."`
	Addr   int    `arg name:"addr" help:"Addresses to access." type:"int"`
}

type StatusCmd struct {
}

func (l *StatusCmd) Run(c *Context) error {
	buf := readmem(c, Region { Region: "RAM", Addr: 0x0000e184 }, 2)
	fmt.Printf("width: %d\n", int(binary.LittleEndian.Uint16(buf)))

	buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e18c }, 2)
	fmt.Printf("height: %d\n", int(binary.LittleEndian.Uint16(buf)))

	buf = readmem(c, Region { Region: "RAM", Addr: 0x0000e180 }, 1)

	if (buf[0] > 0) {
		fmt.Printf("signal: yes\n")
	} else {
		fmt.Printf("signal: no\n")
	}

	return nil
}
