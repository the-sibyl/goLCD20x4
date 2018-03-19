package main

import (
	"fmt"
	"time"
	"github.com/the-sibyl/goLCD20x4"
)

func main() {
	fmt.Println("Running example 7")

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

	lcd.WriteLine("Scrolling text demo", 1)
	lcd.ScrollFromLeft(" ------------------ ", 2, time.Millisecond * 500)
	lcd.ScrollFromRight("Test!", 3, time.Millisecond * 70)
	lcd.ScrollFromLeft("Test!", 4, time.Millisecond * 5)
}
