package main

import (
    "github.com/nsf/termbox-go"
    // "fmt"
)

type Video struct {}

func (v *Video) Init() (error) {
	return termbox.Init()
}

func (v *Video) Close() {
    termbox.Close()
}

func GetColor(color Word) (termbox.Attribute) {
    switch color {
    case 0x0:
        return termbox.ColorBlack
    case 0x1:
        return termbox.ColorRed
    case 0x2:
        return termbox.ColorGreen
    case 0x3:
        return termbox.ColorYellow
    case 0x4:
        return termbox.ColorBlue
    case 0x5:
        return termbox.ColorMagenta
    case 0x6:
        return termbox.ColorCyan
    case 0x7:
        return termbox.ColorWhite
    }

    return termbox.ColorDefault
}

func (v *Video) DrawScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    i := 0
    for y := 0; y < 16; y++ {
        for x := 0; x < 32; x++ {
            data := Memory[0x8000 + i]

            if data > 0 {
                // fmt.Printf("Value: %X\n", data)
            }

            /*
            colors := data >> 8

            fg := GetColor(colors & 0xf)
            bg := GetColor(colors >> 4)
            */

            character := rune(data & 0x7f)
            termbox.SetCell(x, y, character, termbox.ColorDefault, termbox.ColorDefault)

//            character := rune(data & 0x7f)
//            termbox.SetCell(x, y, character, termbox.ColorDefault, termbox.ColorDefault)
            i++
        }
    }

	termbox.Flush()
}
