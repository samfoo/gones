package ppu

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

func TestTileIndexReturnsCorrectByte(t *testing.T) {
    n := NewNametable()

    n.Write(0xff, 0x0000)     // Row 0, column 0
    n.Write(0xef, 0x20)       // Row 1, column 0
    n.Write(0xdf, 1 + 2 * 32) // Row 2, column 1

    assert.Equal(t, n.TileIndex(0, 0), uint8(0xff))
    assert.Equal(t, n.TileIndex(0, 1), uint8(0xef))
    assert.Equal(t, n.TileIndex(1, 2), uint8(0xdf))
}
