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

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    var machine = nes.NewMachine()
    machine.Insert(rom)

    machine.CPU.Debug = false
    machine.CPU.Reset()

    // First step until the tests start.
    for machine.CPU.Memory.Read(0x6000) != 0x80 {
        machine.CPU.Step()
    }

    fmt.Printf("Running tests...\n")

    // Next step until the tests are finished.
    for machine.CPU.Memory.Read(0x6000) != 0x00 {
        machine.CPU.Step()
    }

    fmt.Printf("Results:\n")

    for i:=0; machine.CPU.Memory.Read(cpu.Address(0x6004+i)) != 0x00; i++ {
        char := machine.CPU.Memory.Read(cpu.Address(0x6004+i))
        fmt.Printf("%c", char)
    }
    fmt.Printf("\n")
}

