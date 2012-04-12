package main

import (
    "github.com/nsf/termbox-go"
)

type Video struct {}

func (v *Video) Init() (error) {
	return termbox.Init()
}

func (v *Video) Close() {
    termbox.Close()
}

func (v *Video) DrawScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    i := 0
    for y := 0; y < 16; y++ {
        for x := 0; x < 32; x++ {
            character := rune(Memory[0x8000 + i] & 0xff)
            termbox.SetCell(x, y, character, termbox.ColorDefault, termbox.ColorDefault)
            i++
        }
    }

	termbox.Flush()
}
