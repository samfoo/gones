package main

import (
    "nes"
    "os"
    "fmt"
    "log"
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
    fmt.Printf("Mapper           : %s\n", MapperNames[rom.Header.Mapper])
}
