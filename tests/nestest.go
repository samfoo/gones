package main

import (
    "cpu"
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

    proc := new(cpu.CPU)
    proc.Reset()

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    for i := 0; i < 0x4000; i++ {
        low := 0x8000 + cpu.Address(i)
        high := 0xc000 + cpu.Address(i)

        proc.Memory.Write(rom.Banks[0][i], low)
        proc.Memory.Write(rom.Banks[0][i], high)
    }

    proc.PC = 0xc000

    for i:=0; i < 9000; i++ {
        proc.Step()
    }
}
