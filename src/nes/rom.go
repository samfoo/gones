package nes

import (
    "io"
    "fmt"
    "cpu"
    "ppu"
    "io/ioutil"
    "bytes"
    "errors"
)

var MapperNames = map[uint8]string {
    0:   "NROM",
    1:   "SxROM, MMC1",
    2:   "UxROM",
    3:   "CNROM",
    4:   "TxROM, MMC3, MMC6",
    5:   "ExROM, MMC5",
    7:   "AxROM",
    9:   "PxROM, MMC2",
    10:  "FxROM, MMC4",
    11:  "Color Dreams",
    13:  "CPROM",
    15:  "100-in-1 Contra Function 16",
    16:  "Bandai EPROM (24C02)",
    18:  "Jaleco SS8806",
    19:  "Namco 163",
    21:  "VRC4a, VRC4c",
    22:  "VRC2a",
    23:  "VRC2b, VRC4e",
    24:  "VRC6a",
    25:  "VRC4b, VRC4d",
    26:  "VRC6b",
    34:  "BNROM, NINA-001",
    64:  "RAMBO-1",
    66:  "GxROM, MxROM",
    68:  "After Burner",
    69:  "FME-7, Sunsoft 5B",
    71:  "Camerica/Codemasters",
    73:  "VRC3",
    74:  "Pirate MMC3 derivative",
    75:  "VRC1",
    76:  "Namco 109 variant",
    79:  "NINA-03/NINA-06",
    85:  "VRC7",
    86:  "JALECO-JF-13",
    94:  "Senjou no Ookami",
    105: "NES-EVENT Similar to MMC1",
    113: "NINA-03/NINA-06??",
    118: "TxSROM, MMC3",
    119: "TQROM, MMC3",
    159: "Bandai EPROM (24C01)",
    166: "SUBOR",
    167: "SUBOR",
    180: "Crazy Climber",
    185: "CNROM with protection diodes",
    192: "Pirate MMC3 derivative",
    206: "DxROM, Namco 118 / MIMIC-1",
    210: "Namco 175 and 340",
    228: "Action 52",
    232: "Camerica/Codemasters Quattro",
}

type ROM struct {
    Header *Header
    PrgBanks [][]byte
    ChrBanks [][]byte
    Mapper

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

type Mapper interface {
    Patterntable(int) *ppu.Patterntable
    Graphics() cpu.Mountable
    Program() cpu.Mountable
}

const (
    PrgBankSize = 0x4000
    ChrBankSize = 0x1000
)

func NewMapper(r *ROM) Mapper {
    switch r.Header.Mapper {
        case 0x00: return &NROM { r }
        case 0x01: return &MMC1 { r }
        default:
            panic(fmt.Sprintf("Unsupported mapper %s", MapperNames[r.Header.Mapper]))
    }
}

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

    rom.Mapper = NewMapper(rom)

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
