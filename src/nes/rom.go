package nes

import (
    "io"
    "io/ioutil"
    "bytes"
    "errors"
)

type ROM struct {
    Header *Header
    Banks [][]byte

    data []byte
}

const (
    PrgBankSize = 0x4000
    ChrBankSize = 0x2000
)

func ReadROM(r io.Reader) (rom *ROM, err error) {
    rom = new(ROM)

    var raw []byte
    raw, err = ioutil.ReadAll(r)
    if err != nil { return }

    rom.data = raw[16:]

    rom.Header, err = ParseHeader(raw)
    if err != nil { return }

    rom.Banks = make([][]byte, rom.Header.PrgRomSize)
    for i := 0; i < rom.Header.PrgRomSize; i++ {
        rom.Banks[i] = rom.data[PrgBankSize*i:PrgBankSize*(i+1)]
    }

    return
}

type Header struct {
    PrgRomSize int
    ChrRomSize int
    Mapper uint8
    Flags6 byte
    Flags7 byte
    PrgRamSize uint8
    Flags9 byte
}

var nes = []byte{0x4E, 0x45, 0x53, 0x1A}

func ParseHeader(raw []byte) (header *Header, err error) {
    header = new(Header)

    if !bytes.Equal(raw[:4], nes) {
        err = errors.New("iNES header invalid. Is this really an NES ROM?")
    }

    header.PrgRomSize = int(raw[4])
    header.ChrRomSize = int(raw[5])
    header.Mapper = (raw[6] >> 4) | (raw[7] & 0xF0)

    return
}
