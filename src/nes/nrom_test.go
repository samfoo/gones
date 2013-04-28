package nes

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestNROMSelectRightChrBank(t *testing.T) {
    rom := new(ROM)

    rom.ChrBanks = make([][]byte, 2)
    rom.ChrBanks[0] = make([]byte, 0x1000)
    rom.ChrBanks[1] = make([]byte, 0x1000)

    rom.ChrBanks[0][0] = 0xef
    rom.ChrBanks[1][0] = 0xbe

    mapper := NROM { rom }

    assert.Equal(t, mapper.Graphics().Read(0x1000), byte(0xbe))
    assert.Equal(t, mapper.Graphics().Read(0x0000), byte(0xef))
}

func TestNROM128MirrorsMemory(t *testing.T) {
    rom := new(ROM)

    rom.PrgBanks = make([][]byte, 1)
    rom.PrgBanks[0] = make([]byte, 16 * 0x4000)

    rom.PrgBanks[0][0] = 0xbe

    mapper := NROM { rom }

    assert.Equal(t, mapper.Program().Read(0x0000), byte(0xbe))
    assert.Equal(t, mapper.Program().Read(0x4000), byte(0xbe))
}

func TestNROM256DoesntMirror(t *testing.T) {
    rom := new(ROM)

    rom.PrgBanks = make([][]byte, 2)
    rom.PrgBanks[0] = make([]byte, 16 * 0x4000)
    rom.PrgBanks[1] = make([]byte, 16 * 0x4000)

    rom.PrgBanks[0][0] = 0xbe
    rom.PrgBanks[1][0] = 0xef

    mapper := NROM { rom }

    assert.Equal(t, mapper.Program().Read(0x0000), byte(0xbe))
    assert.Equal(t, mapper.Program().Read(0x4000), byte(0xef))
}
