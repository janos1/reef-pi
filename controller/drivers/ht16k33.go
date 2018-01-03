package drivers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	'3': 207,
	'4': 230,
	'5': 237,
	'6': 253,
	'7': 7,
	'8': 255,
	'9': 239,
	'A': 247,
	'B': 4815,
	'C': 57,
	'D': 4623,
	'E': 249,
	'F': 113,
	'G': 189,
	'H': 246,
	'I': 4617,
	'J': 30,
	'K': 9328,
	'L': 56,
	'M': 1334,
	'N': 8502,
	'O': 63,
	'P': 243,
	'Q': 8255,
	'R': 8435,
	'S': 237,
	'T': 4609,
	'U': 62,
	'V': 3120,
	'W': 10294,
	'X': 11520,
	'Y': 5376,
	'Z': 3081,
	' ': 0,
}

var allRune = []rune{
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
	bytes := make([]byte, 16)
	if err := bus.WriteToReg(addr, 0x00, bytes); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	fmt.Println(s)
	l := len(s)
	for i := 0; 
	for i, c := range s {
		row := i % 4
		item := digits[c]
		bytes[row*2], bytes[row*2+1] = byte(item), byte(item>>8)
		bus.WriteToReg(addr, 0x00, bytes)
		if l == i+1 {
			continue
		}
		row++
		item = digits[rune(s[i+1])]
		bytes[row*2], bytes[row*2+1] = byte(item), byte(item>>8)
		bus.WriteToReg(addr, 0x00, bytes)

		if l == i+2 {
			continue
		}
		row++
		item = digits[rune(s[i+2])]
		bytes[row*2], bytes[row*2+1] = byte(item), byte(item>>8)
		bus.WriteToReg(addr, 0x00, bytes)

		if l == i+3 {
			continue
		}
		row++
		item = digits[rune(s[i+3])]
		bytes[row*2], bytes[row*2+1] = byte(item), byte(item>>8)
		bus.WriteToReg(addr, 0x00, bytes)
		time.Sleep(time.Second)
	}
	return nil
}
