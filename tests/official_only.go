package main

import (
    "cpu"
    "nes"
    "os"
    "log"
    "fmt"
)

func main() {
    var file *os.File
    var err error
    if file, err = os.Open("assets/official_only.nes"); err != nil {
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
        proc.Memory.Write(rom.Banks[len(rom.Banks)-1][i], high)
    }

    proc.PC = cpu.Address(proc.Memory.Read(0xFFFC)) |
        (cpu.Address(proc.Memory.Read(0xFFFD))<<8)

    // First step until the tests start.
    for proc.Memory.Read(0x6000) != 0x80 {
        proc.Step()
    }

    fmt.Printf("Running tests...\n")

    // Next step until the tests are finished.
    for proc.Memory.Read(0x6000) != 0x00 {
        proc.Step()
    }

    fmt.Printf("Results:\n")

    for i:=0; proc.Memory.Read(cpu.Address(0x6004+i)) != 0x00; i++ {
        char := proc.Memory.Read(cpu.Address(0x6004+i))
        fmt.Printf("%c", char)
    }
    fmt.Printf("\n")
}

