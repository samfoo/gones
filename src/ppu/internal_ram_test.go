package ppu

import (
    "cpu"
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestNormalizeMirrorsNametables(t *testing.T) {
    ram := new(InternalRAM)

    assert.Equal(t, ram.normalize(0x3000), cpu.Address(0x2000))
    assert.Equal(t, ram.normalize(0x31f0), cpu.Address(0x21f0))
    assert.Equal(t, ram.normalize(0x3eff), cpu.Address(0x2eff))
}

func TestNormalizeMirrorsPaletteIndexes(t *testing.T) {
    ram := new(InternalRAM)

    assert.Equal(t, ram.normalize(0x3f20), cpu.Address(0x3f00))
    assert.Equal(t, ram.normalize(0x3f40), cpu.Address(0x3f00))
    assert.Equal(t, ram.normalize(0x3f60), cpu.Address(0x3f00))
    assert.Equal(t, ram.normalize(0x3fff), cpu.Address(0x3f1f))
}
