package main

import (
    "os"
    "cpu"
    "log"
)

func main() {
    var file *os.File
    var err error
    if file, err = os.Open("assets/nestest"); err != nil {
        log.Fatal(err)
        return
    }

    proc := new(cpu.CPU)
    proc.Reset()

    rom := make([]byte, 0x10000)

    if _, err := file.Read(rom); err != nil {
        log.Fatal(err)
        return
    }

    for i := 0; i < 0x10000; i++ {
        proc.Memory.Write(0x00, cpu.Address(i))
    }

    for i := 0; i < 0x4000; i++ {
        location := 0xc000 + cpu.Address(i)

        proc.Memory.Write(rom[i], location)
    }

    proc.PC = 0xc000

    for i:=0; i < 9000; i++ {
        proc.Step()
    }
}
