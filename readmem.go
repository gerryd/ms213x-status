package main

import (
	"github.com/BertoldVdb/ms-tools/mshal"
)

func readmem(c *Context, request Region, amount int) []byte {
	region := c.hal.MemoryRegionGet(mshal.MemoryRegionNameType(request.Region))

	if amount == 0 {
		amount = region.GetLength()
	}

	buf := make([]byte, amount)
	n, err := region.Access(false, request.Addr, buf)

	if err != nil {
		return []byte{}
	}

	buf = buf[:n]

	return buf
}
