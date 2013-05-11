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

func TestReturnsCorrectAttribute(t *testing.T) {
    n := NewNametable()

    /*
    Given an attribute table entry arranged in the following manner:

    +--+--+--+
    |AA|BB|...
    +--+--+
    |CC|DD|
    +--+--+
    |...

    The attribute table entry bits should be spread like:

    bit index  = 01 23 45 67
    group      = DD CC BB AA

    And each entry represents 4 nametable tiles in a particular region. So the
    nametables in the above diagram are colored like:

        0  1  2  3
      +--+--+--+--+--+
     0| A| A| B| B|...
      +--+--+--+--+
     1| A| A| B| B|
      +--+--+--+--+
     2| C| C| D| D|
      +--+--+--+--+
     3| C| C| D| D|
      +--+--+--+--+
      |...

    */

    // In this test...
    a := uint8(3)
    b := uint8(2)
    c := uint8(1)
    d := uint8(0)

    n.Write(0x1b, 0x03c0)

    // In the above example, this tests all tiles in group A
    assert.Equal(t, n.Attribute(0, 0), a)
    assert.Equal(t, n.Attribute(0, 1), a)
    assert.Equal(t, n.Attribute(1, 0), a)
    assert.Equal(t, n.Attribute(1, 1), a)

    // Test all tiles in group B
    assert.Equal(t, n.Attribute(2, 0), b)
    assert.Equal(t, n.Attribute(3, 0), b)
    assert.Equal(t, n.Attribute(2, 1), b)
    assert.Equal(t, n.Attribute(3, 1), b)

    // Test all tiles in group C
    assert.Equal(t, n.Attribute(0, 2), c)
    assert.Equal(t, n.Attribute(1, 2), c)
    assert.Equal(t, n.Attribute(0, 3), c)
    assert.Equal(t, n.Attribute(1, 3), c)

    // Test all tiles in group D
    assert.Equal(t, n.Attribute(2, 2), d)
    assert.Equal(t, n.Attribute(3, 2), d)
    assert.Equal(t, n.Attribute(2, 3), d)
    assert.Equal(t, n.Attribute(3, 3), d)
}
