package drivers

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	REGISTER_DISPLAY_SETUP = 0x80
	REGISTER_SYSTEM_SETUP  = 0x20
	REGISTER_DIMMING       = 0xE0
	BLINKRATE_OFF          = 0x00
	BLINKRATE_2HZ          = 0x01
	BLINKRATE_1HZ          = 0x02
	BLINKRATE_HALFHZ       = 0x03
)

var digits = map[rune]uint16{
	'0': 63,
	'1': 6,
	'2': 219,
	'3': 143,
	'4': 230,
	'5': 105, //
	'6': 253,
	'7': 7,
	'8': 255,
	'9': 239,
	'A': 247,
	'B': 143,
	'C': 57,
	'D': 15,
	'E': 249,
	'F': 113,
	'G': 189,
	'H': 246,
	'I': 0,
	'J': 30,
	'K': 112,
	'L': 56,
	'M': 54,
	'N': 54,
	'O': 63,
	'P': 243,
	'Q': 63,
	'R': 243,
	'S': 237,
	'T': 1,
	'U': 62,
	'V': 48,
	'W': 54,
	'X': 0,
	'Y': 0,
	'Z': 9,
}

type HT16K33 struct {
}

func Demo() error {
	bus, err := New()
	if err != nil {
		return err
	}

	defer bus.Close()
	addr := byte(0x70)
	if err := bus.WriteToReg(addr, REGISTER_SYSTEM_SETUP|0x01, []byte{0x00}); err != nil {
		return err
	}

	if err := bus.WriteToReg(addr, REGISTER_DIMMING|5, []byte{0x00}); err != nil {
		return err
	}
	runes := []rune{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'B',
		'C',
		'D',
		'E',
		'F',
		'G',
		'H',
		'I',
		'J',
		'K',
		'L',
		'M',
		'N',
		'O',
		'P',
		'Q',
		'R',
		'S',
		'T',
		'U',
		'V',
		'W',
		'X',
		'Y',
		'Z',
	}
	bytes := make([]byte, 16)
	if err := bus.WriteToReg(addr, 0x00, bytes); err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	for _, r := range runes {
		item := digits[r]
		bytes[0], bytes[1] = byte(item), byte(item>>8)
		if err := bus.WriteToReg(addr, 0x00, bytes); err != nil {
			return err
		}
		fmt.Printf("%s\n", string(r))
		reader.ReadString('\n')
		time.Sleep(time.Second)
	}
	/*
		for row := 0; row < 8; row++ {
			inc := uint16(1)
			i := inc
			for j := 0; j < 8; j++ {
				bytes[row] = byte(i)
				if err := bus.WriteToReg(addr, 0x00, bytes); err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 100)
				inc = inc << 1
				i = i + inc
			}
			fmt.Println("Row:", row)
		}
		return bus.WriteToReg(addr, 0x00, make([]byte, 17))
	*/
	return nil
}
