DCPU16 software stack written in Go
=====

# Contains
* Emulator
* Assembler
* Disassembler

## To run the emulator
    $ go run dcpu16.go video.go main.go <inputfile.dcpx>

## To build the emulator
    $ go build dcpu16.go video.go main.go

## To run the assembler
    $ go run dcpu16.go assembler.go <inputfile.dasm>

## To build the assembler
    $ go build dcpu16.go assembler.go

## To run the disassembler
    $ go run dcpu16.go disassembler.go <inputfile.dcpx>

## To build the disassembler
    $ go run dcpu16.go disassembler.go
