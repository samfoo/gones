package main

import (
    "cpu"
    "nes"
    "os"
    "log"
    "fmt"
    "strings"
    "path/filepath"
)

func run(filename string) {
    var file *os.File
    var err error
    if file, err = os.Open(filename); err != nil {
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

    machine.CPU.Cycle = func() {
        for i:=0; i < 3; i++ {
            machine.PPU.Step()
        }
    }

    // First step until the tests start.
    for machine.CPU.Memory.Read(0x6000) != 0x80 {
        machine.CPU.Step()
    }

    parts := strings.Split(filename, string(filepath.Separator))
    fmt.Printf("Running %s tests...\n", parts[len(parts)-1])

    // Next step until the tests are finished.
    for machine.CPU.Memory.Read(0x6000) > 0x7f {
        machine.CPU.Step()

        if machine.CPU.Memory.Read(0x6000) == 0x81 {
            panic("Need to press reset, but isn't implemented")
        }
    }

    var output = ""
    for i:=0; machine.CPU.Memory.Read(cpu.Address(0x6004+i)) != 0x00; i++ {
        char := machine.CPU.Memory.Read(cpu.Address(0x6004+i))
        output += fmt.Sprintf("%c", char)
    }

    if strings.Contains(output, "Passed") {
        fmt.Printf("Results: OK\n")
    } else {
        fmt.Printf("Results: FAIL\n")
        fmt.Printf("====\n%s====\n", output)
    }
}

func main() {
    filenames := os.Args[1:]

    for i:=0; i < len(filenames); i++ {
        run(filenames[i])
    }
}

