package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/BertoldVdb/ms-tools/mshal/ms213x"
)

// local copy of calcSum because I can't figure out how to call it in the ms213x
// package
func calcSum(f []byte) uint16 {
	var csum uint16
	for _, m := range f {
		csum += uint16(m)
	}
	return csum
}

type FirmwareCmd struct {
	Mode      string `name:"mode" help:"What to do (check or write)" default:check`
	Filename  string `name:"filename" help:"Filename to do things with"`
}

func (f *FirmwareCmd) Run(c *Context) error {
	if (f.Mode == "write") {
		fw, err := os.ReadFile(f.Filename)
		if err != nil {
			return err
		}

		codeLen := int(binary.BigEndian.Uint16(fw[2:]))
		end := 0x30 + codeLen

		hdrSum := calcSum(fw[2:12]) + calcSum(fw[16:0x30])
		codeSum := calcSum(fw[0x30:end])

		hdrSumByteArray := make([]byte, 2)
		codeSumByteArray := make([]byte, 2)

		binary.BigEndian.PutUint16(hdrSumByteArray, uint16(hdrSum))
		binary.BigEndian.PutUint16(codeSumByteArray, uint16(codeSum))

		copy(fw[end:], hdrSumByteArray)
		copy(fw[end+2:], codeSumByteArray)

		// running our modified fw through a check because belt and suspenders
		check := ms213x.CheckImage(fw)

		if (check != nil) {
			return fmt.Errorf("calculated checksum does not match, not writing: %x vs %x", hdrSum, codeSum)
		}

		err = os.WriteFile(f.Filename, fw, 0644)

		if (err != nil) {
			return err
		}

		fmt.Println("firmware file with new checksums written")
	}

	fw, err := os.ReadFile(f.Filename)
	if err != nil {
		return err
	}

	check := ms213x.CheckImage(fw)

	if (check == nil) {
		fmt.Println("checksums are correct")
	} else {
		return fmt.Errorf("checksums are wrong")
	}

	return nil
}
