package main

import (
	"fmt"
)

func hexdump(offset int, data []byte) string {
	var result string

	for len(data) > 0 {
		l := len(data)
		if l > 32 {
			l = 32
		}
		work := data[:l]
		data = data[l:]

		var workHex string
		var workAscii string
		for i := 0; i < 32; i++ {
			m := byte(0)
			valid := i < len(work)

			if valid {
				m = work[i]
			}

			if valid {
				workHex += fmt.Sprintf("%02x ", m)

				if m < 32 || m > 126 {
					m = '.'
				}

				workAscii += fmt.Sprintf("%c", m)
			} else {
				workHex += "   "
				workAscii += " "
			}
			if i%8 == 7 {
				workHex += " "
			}
		}

		result += fmt.Sprintf("%08x  %s|%s|\n", offset, workHex, workAscii)
		offset += l
	}

	return result
}
