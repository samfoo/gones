package main

import (
    "nes"
    "os"
    "fmt"
    "log"
)

func main() {
    path := os.Args[1]

    var file *os.File
    var err error
    if file, err = os.Open(path); err != nil {
        log.Fatal(err)
        return
    }

    var rom *nes.ROM
    rom, err = nes.ReadROM(file)
    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Printf("Prg ROM (* 16KB) : %d\n", rom.Header.PrgRomSize)
    fmt.Printf("Chr ROM (* 8KB)  : %d\n", rom.Header.ChrRomSize)
    fmt.Printf("Mapper           : %s\n", nes.MapperNames[rom.Header.Mapper])
}
