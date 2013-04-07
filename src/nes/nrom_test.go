package nes

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestNROM128MirrorsMemory(t *testing.T) {
    rom := new(ROM)

    rom.Banks = make([][]byte, 1)
    rom.Banks[0] = make([]byte, 16 * 0x4000)

    rom.Banks[0][0] = 0xbe

    mapper := NROM { rom }

    assert.Equal(t, mapper.Read(0x0000), byte(0xbe))
    assert.Equal(t, mapper.Read(0x4000), byte(0xbe))
}

func TestNROM256DoesntMirror(t *testing.T) {
    rom := new(ROM)

    rom.Banks = make([][]byte, 2)
    rom.Banks[0] = make([]byte, 16 * 0x4000)
    rom.Banks[1] = make([]byte, 16 * 0x4000)

    rom.Banks[0][0] = 0xbe
    rom.Banks[1][0] = 0xef

    mapper := NROM { rom }

    assert.Equal(t, mapper.Read(0x0000), byte(0xbe))
    assert.Equal(t, mapper.Read(0x4000), byte(0xef))
}
