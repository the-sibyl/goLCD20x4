package main

import (
	"fmt"
	"github.com/the-sibyl/goLCD20x4"
	"time"
)

func main() {
	fmt.Println("Running example 4")

	// This code was tested on an RPi Zero 1.3 with a NewHaven (Chinese) NHD-0420DZ-FL-YBW-3V3.
	// The product on the site appears to end in "33V3" which is probably a mistake.
	// The BCM GPIO numbers are as follows.
	// 2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5
	// RS, RW, E, DB0, DB1, DB2, DB3, DB4, DB5, DB6, DB7

	lcd := goLCD20x4.Open(2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5)
	defer lcd.Close()

	lcd.FunctionSet(1, 1, 0)
	lcd.DisplayOnOffControl(1, 0, 0)
	lcd.EntryModeSet(1, 0)

	lcd.ClearDisplay()

	lcd.WriteLine("System date and time", 2)
	time.Sleep(time.Millisecond * 500)
	lcd.WriteLine(" [Fetching time...]", 3)
	time.Sleep(time.Millisecond * 2000)
	for {
		lcd.WriteLine(time.Now().String(), 3)
		time.Sleep(time.Millisecond * 100)
	}
}
