package main

import (
    "nes"
    "os"
    "log"
)

func main() {
    var file *os.File
    var err error
    if file, err = os.Open("assets/nestest.nes"); err != nil {
        log.Fatal(err)
        return
    }

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    var machine = nes.NewMachine()
    machine.Insert(rom)

    machine.CPU.Debug = true
    machine.CPU.Reset()
    machine.CPU.PC = 0xc000

    for i:=0; i < 9000; i++ {
        machine.CPU.Step()
    }
}
