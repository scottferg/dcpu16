package main

import (
    "github.com/nsf/termbox-go"
    "io/ioutil"
    "os"
    "fmt"
)

func run(program []Word) {
    cpu := new(Dcpu)

    Memory = make([]Word, 0xffff)

    cpu.PC = 0
    cpu.SP = 0xfffe

    video := new(Video)
    video.Init()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			ev := termbox.PollEvent()
			event_queue <- ev
		}
	}()

    for index, value := range program {
        Memory[index] = value
    }

runloop:
    for int(cpu.PC) < len(program) {
        cpu.Step()

		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break runloop
			}
        default:
            cpu.Step()
            video.DrawScreen()
        }
    }

    video.Close()
}

func main() {
    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        fileLength := len(contents)

        var program []Word

        for i := 0; i < fileLength; i += 2 {
            fmt.Printf("%d | %d\n", i, i + 1)
            program = append(program, (Word(contents[i]) << 8) + Word(contents[i + 1]))
        }

        run(program)
    }
}
