package nes

import (
    "io"
    "bytes"
    "errors"
)

type ROM struct {
    Data []byte
}

func ReadROM(r io.Reader) (rom *ROM, err error) {
    rom = new(ROM)

    var header *ines
    header, err = ReadHeader(r)
    if err != nil { return }

    rom.Data = make([]byte, header.PrgRomSize)
    _, err = r.Read(rom.Data)
    if err != nil { return }

    return
}

type ines struct {
    PrgRomSize int
    ChrRomSize int
    Mapper uint8
    flags6 byte
    flags7 byte
    prgRamSize uint8
    flags9 byte
}

var nes = []byte{0x4E, 0x45, 0x53, 0x1A}

func ReadHeader(r io.Reader) (header *ines, err error) {
    header = new(ines)
    h := make([]byte, 16)

    _, err = r.Read(h)
    if err != nil { return }

    if !bytes.Equal(h[:4], nes) {
        err = errors.New("iNES header invalid. Is this really an NES ROM?")
    }

    header.PrgRomSize = int(h[4]) * 16384
    header.ChrRomSize = int(h[5]) * 8192

    header.Mapper = (h[6] >> 4) | (h[7] & 0xF0)

    return
}
