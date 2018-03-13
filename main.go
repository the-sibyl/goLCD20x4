package main

import (
	"fmt"
	"time"

	"github.com/the-sibyl/sysfsGPIO"
)

// Note: I am copy-pasting a lot of code in case I need to change timing for a particular instruction. It looks as 
// though most of these instructions can be collapsed into functions that call a generalized register write function.

func main() {
	fmt.Println("Running")
	//	font := make(map[rune]int)

	// 2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5
	// RS, RW, E, DB0, DB1, DB2, DB3, DB4, DB5, DB6, DB7

	lcd := Open(2, 3, 4, 17, 27, 22, 10, 9, 11, 0, 5)
	defer lcd.Close()

	lcd.FunctionSet(1, 1, 1)
	lcd.DisplayOnOffControl(1, 0, 0)
	lcd.EntryModeSet(0, 0)

	for {
		lcd.ReturnHome()

		for k := 0; k < 80; k++ {
			lcd.WriteCharacter((byte) (k + 0x30))
		}

		time.Sleep(time.Second)

		lcd.ClearDisplay()
		lcd.SetDDRAMAddress(0)
		lcd.WriteCharacter(0x30)
		lcd.SetDDRAMAddress(40)
		lcd.WriteCharacter(0x31)
		lcd.SetDDRAMAddress(20)
		lcd.WriteCharacter(0x32)
		lcd.SetDDRAMAddress(60)
		lcd.WriteCharacter(0x33)

		time.Sleep(time.Second)
	}
}



// Write a raw character supported in the CGROM
func (lcd *LCD20x4) WriteCharacter(rawCharacter byte) {
	lcd.PinRS.SetHigh()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.writeDBBus(rawCharacter)

	time.Sleep(time.Microsecond)

	lcd.PinE.SetLow()

	time.Sleep(time.Microsecond)
}

func (lcd *LCD20x4) ClearDisplay() {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.writeDBBus(1)

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}

func (lcd *LCD20x4) ReturnHome() {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.writeDBBus(0)

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}


func (lcd *LCD20x4) EntryModeSet(incrementOrDecrement int, shiftEntireDisplay int) {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.PinDB7.SetLow()
	lcd.PinDB6.SetLow()
	lcd.PinDB5.SetLow()
	lcd.PinDB4.SetLow()
	lcd.PinDB3.SetLow()
	lcd.PinDB2.SetLow()

	// Cursor moves to the right if set true or left if set false
	if incrementOrDecrement > 0 {
		lcd.PinDB1.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}

	// Shift the entire display to the right if set true or left if set false
	if shiftEntireDisplay > 0 {
		lcd.PinDB0.SetHigh()
	} else {
		lcd.PinDB0.SetLow()
	}

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}

func (lcd *LCD20x4) DisplayOnOffControl(displayOnOff int, cursorOnOff int, cursorBlinking int) {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.PinDB7.SetLow()
	lcd.PinDB6.SetLow()
	lcd.PinDB5.SetLow()
	lcd.PinDB4.SetLow()
	lcd.PinDB3.SetHigh()


	// Set the entire display on
	if displayOnOff > 0 {
		lcd.PinDB2.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}

	// Set the cursor on or off
	if cursorOnOff > 0 {
		lcd.PinDB1.SetHigh()
	} else {
		lcd.PinDB1.SetLow()
	}

	// Set cursor blinking on or off
	if cursorBlinking > 0 {
		lcd.PinDB0.SetHigh()
	} else {
		lcd.PinDB0.SetLow()
	}

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}

func (lcd *LCD20x4) CursorOrDisplayShift(displayShiftOrCursorMove int, rightOrLeft int) {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.PinDB7.SetLow()
	lcd.PinDB6.SetLow()
	lcd.PinDB5.SetLow()
	lcd.PinDB4.SetHigh()

	// Display shift if true. Cursor move if false.
	if displayShiftOrCursorMove > 0 {
		lcd.PinDB3.SetHigh()
	} else {
		lcd.PinDB3.SetLow()
	}

	// Shift to the right if true. Shift to the left if false.
	if rightOrLeft > 0 {
		lcd.PinDB2.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}

	// Note: DB1 and DB0 are "don't care" values

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}

func (lcd *LCD20x4) FunctionSet(dataLength int, numDisplayLines int, characterFont int) {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)

	lcd.PinDB7.SetLow()
	lcd.PinDB6.SetLow()
	lcd.PinDB5.SetHigh()

	// Select between 8-bit and 4-bit bus width
	if dataLength > 0 {
		lcd.PinDB4.SetHigh()
	} else {
		lcd.PinDB4.SetLow()
	}

	// Select between 2-line display mode or 1-line display mode
	if numDisplayLines > 0 {
		lcd.PinDB3.SetHigh()
	} else {
		lcd.PinDB3.SetLow()
	}

	// Choose the display font type, 5x8 or 5x11
	if characterFont > 0 {
		lcd.PinDB2.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}

	// Note: DB1 and DB0 are "don't care" values

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Millisecond)
}

func (lcd *LCD20x4) SetDDRAMAddress(address byte) {
	lcd.PinRS.SetLow()
	lcd.PinRW.SetLow()

	time.Sleep(time.Microsecond)

	lcd.PinE.SetHigh()

	time.Sleep(time.Microsecond)


	shiftedValue := address

	if shiftedValue & 1 != 0 {
		lcd.PinDB0.SetHigh()
	} else {
		lcd.PinDB0.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB1.SetHigh()
	} else {
		lcd.PinDB1.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB2.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB3.SetHigh()
	} else {
		lcd.PinDB3.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB4.SetHigh()
	} else {
		lcd.PinDB4.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB5.SetHigh()
	} else {
		lcd.PinDB5.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB6.SetHigh()
	} else {
		lcd.PinDB6.SetLow()
	}

	lcd.PinDB7.SetHigh()

	// Write operations appear to need a >= 1000ns settling time
	time.Sleep(time.Microsecond)

	// Latch data into the device
	lcd.PinE.SetLow()

	time.Sleep(time.Microsecond)
}

// Helper function to write the 8-bit data bus
func (lcd *LCD20x4) writeDBBus(value byte) {
	shiftedValue := value

	if shiftedValue & 1 != 0 {
		lcd.PinDB0.SetHigh()
	} else {
		lcd.PinDB0.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB1.SetHigh()
	} else {
		lcd.PinDB1.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB2.SetHigh()
	} else {
		lcd.PinDB2.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB3.SetHigh()
	} else {
		lcd.PinDB3.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB4.SetHigh()
	} else {
		lcd.PinDB4.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB5.SetHigh()
	} else {
		lcd.PinDB5.SetLow()
	}
	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB6.SetHigh()
	} else {
		lcd.PinDB6.SetLow()
	}

	shiftedValue = shiftedValue >> 1

	if shiftedValue & 1 != 0 {
		lcd.PinDB7.SetHigh()
	} else {
		lcd.PinDB7.SetLow()
	}
	shiftedValue = shiftedValue >> 1
}
type LCD20x4 struct {
	PinRS  *sysfsGPIO.IOPin
	PinRW  *sysfsGPIO.IOPin
	PinE   *sysfsGPIO.IOPin
	PinDB0 *sysfsGPIO.IOPin
	PinDB1 *sysfsGPIO.IOPin
	PinDB2 *sysfsGPIO.IOPin
	PinDB3 *sysfsGPIO.IOPin
	PinDB4 *sysfsGPIO.IOPin
	PinDB5 *sysfsGPIO.IOPin
	PinDB6 *sysfsGPIO.IOPin
	PinDB7 *sysfsGPIO.IOPin
}

func Open(rsGPIONum int, rwGPIONum int, eGPIONum int, db0GPIONum int, db1GPIONum int, db2GPIONum int, db3GPIONum int,
	db4GPIONum int, db5GPIONum int, db6GPIONum int, db7GPIONum int) *LCD20x4 {
	var lcd LCD20x4
	var err error

	lcd.PinRS, err = sysfsGPIO.InitPin(rsGPIONum, "out")
	pinErrHandler(err)

	lcd.PinRW, err = sysfsGPIO.InitPin(rwGPIONum, "out")
	pinErrHandler(err)

	lcd.PinE, err = sysfsGPIO.InitPin(eGPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB0, err = sysfsGPIO.InitPin(db0GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB1, err = sysfsGPIO.InitPin(db1GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB2, err = sysfsGPIO.InitPin(db2GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB3, err = sysfsGPIO.InitPin(db3GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB4, err = sysfsGPIO.InitPin(db4GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB5, err = sysfsGPIO.InitPin(db5GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB6, err = sysfsGPIO.InitPin(db6GPIONum, "out")
	pinErrHandler(err)

	lcd.PinDB7, err = sysfsGPIO.InitPin(db7GPIONum, "out")
	pinErrHandler(err)

	return &lcd
}

func pinErrHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func (lcd *LCD20x4) Close() {
	lcd.PinRS.ReleasePin()
	lcd.PinRW.ReleasePin()
	lcd.PinE.ReleasePin()
	lcd.PinDB0.ReleasePin()
	lcd.PinDB1.ReleasePin()
	lcd.PinDB2.ReleasePin()
	lcd.PinDB3.ReleasePin()
	lcd.PinDB4.ReleasePin()
	lcd.PinDB5.ReleasePin()
	lcd.PinDB6.ReleasePin()
	lcd.PinDB7.ReleasePin()
}
