package main

import (
	"fmt"
	"github.com/the-sibyl/goLCD20x4"
)

func main() {
	fmt.Println("Running example 5")

	// This code was tested on an RPi Zero 1.3 with a NewHaven (Chinese) NHD-0420DZ-FL-YBW-3V3.
	// The product on the site appears to end in "33V3" which is probably a mistake.
	// The BCM GPIO numbers are as follows.
	// 2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5
	// RS, RW, E, DB0, DB1, DB2, DB3, DB4, DB5, DB6, DB7

	lcd := goLCD20x4.Open(2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5)
	defer lcd.Close()

	lcd.FunctionSet(1, 1, 1)
	lcd.DisplayOnOffControl(1, 0, 0)
	lcd.EntryModeSet(1, 0)

	lcd.ClearDisplay()

	for k := 0; k < 4; k++ {
		currentLine := ""
		for j := 0; j < 20; j++ {
			// Change the 0x20 to 0xA0 to see the second character set available on some displays. I cannot
			// get this to work with the NewHaven display. I see a few Japanese characters but not what I 
			// am expecting. Changing between 5x11 and 5x8 font size with the last parameter of 
			// FunctionSet() had no effect.
			currentChar := 0x20 + j + k * 20
			fmt.Println("currentChar:", currentChar)
			currentLine += string(currentChar)
		}

		lcd.WriteLine(currentLine, k + 1)
		fmt.Println(currentLine)
	}
}
