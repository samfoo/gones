package nes

import (
    "io"
    "cpu"
    "io/ioutil"
    "bytes"
    "errors"
)

type ROM struct {
    Header *Header
    PrgBanks [][]byte
    ChrBanks [][]byte
    Cartridge

    data []byte
}

type MountableStruct struct {
    read func(cpu.Address)byte
    write func(byte, cpu.Address)
}

func (ms *MountableStruct) Read(location cpu.Address) byte {
    return ms.read(location)
}

func (ms *MountableStruct) Write(val byte, location cpu.Address) {
    ms.write(val, location)
}

type Cartridge interface {
    Graphics() cpu.Mountable
    Program() cpu.Mountable
}

const (
    PrgBankSize = 0x4000
    ChrBankSize = 0x1000
)

func ReadROM(r io.Reader) (rom *ROM, err error) {
    rom = new(ROM)

    var raw []byte
    raw, err = ioutil.ReadAll(r)
    if err != nil { return }

    rom.data = raw[16:]

    rom.Header, err = ParseHeader(raw)
    if err != nil { return }

    rom.PrgBanks = make([][]byte, rom.Header.PrgRomSize)
    for i := 0; i < rom.Header.PrgRomSize; i++ {
        rom.PrgBanks[i] = rom.data[PrgBankSize*i:PrgBankSize*(i+1)]
    }

    rom.ChrBanks = make([][]byte, rom.Header.ChrRomSize*2)
    for i := 0; i < rom.Header.ChrRomSize*2; i++ {
        offset := PrgBankSize * rom.Header.PrgRomSize
        start := offset + ChrBankSize * i
        end := offset + ChrBankSize * (i + 1)
        rom.ChrBanks[i] = rom.data[start:end]
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
