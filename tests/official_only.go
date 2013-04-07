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

    proc := cpu.NewCPU()

    // Swallow anything to the APU or PPU
    proc.Memory.Mount(cpu.NewRAM(0x6000), 0x2000, 0x7fff)

    proc.Debug = false
    proc.Reset()

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    var mapper = new(nes.MMC1)
    mapper.Rom = rom
    proc.Memory.Mount(mapper, 0x8000, 0xffff)

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

