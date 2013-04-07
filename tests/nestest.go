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

    proc := cpu.NewCPU()

    // Swallow anything to the APU or PPU
    proc.Memory.Mount(cpu.NewRAM(0x6000), 0x2000, 0x7fff)

    proc.Debug = true
    proc.Reset()

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    var mapper = new(nes.NROM)
    mapper.Rom = rom
    proc.Memory.Mount(mapper, 0x8000, 0xffff)

    proc.PC = 0xc000

    for i:=0; i < 9000; i++ {
        proc.Step()
    }
}
