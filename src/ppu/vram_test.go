package ppu

import (
    "cpu"
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestNormalizeMirrorsPaletteIndexes(t *testing.T) {
    ram := new(VRAM)

    assert.Equal(t, ram.normalize(0x0020), cpu.Address(0x0000))
    assert.Equal(t, ram.normalize(0x0040), cpu.Address(0x0000))
    assert.Equal(t, ram.normalize(0x0060), cpu.Address(0x0000))
    assert.Equal(t, ram.normalize(0x00ff), cpu.Address(0x001f))
}
