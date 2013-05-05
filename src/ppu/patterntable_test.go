package ppu

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

/*
    Bit Planes            Pixel Pattern
    $0xx0=$41  01000001
    $0xx1=$C2  11000010
    $0xx2=$44  01000100
    $0xx3=$48  01001000
    $0xx4=$10  00010000
    $0xx5=$20  00100000         .1.....3
    $0xx6=$40  01000000         11....3.
    $0xx7=$80  10000000  =====  .1...3..
                                .1..3...
    $0xx8=$01  00000001  =====  ...3.22.
    $0xx9=$02  00000010         ..3....2
    $0xxA=$04  00000100         .3....2.
    $0xxB=$08  00001000         3....222
    $0xxC=$16  00010110
    $0xxD=$21  00100001
    $0xxE=$42  01000010
    $0xxF=$87  10000111
*/

var tile = []byte {
    0x41, 0xc2, 0x44, 0x48, 0x10, 0x20, 0x40, 0x80,
    0x01, 0x02, 0x04, 0x08, 0x16, 0x21, 0x42, 0x87,
}

func TestTileReadsPixelsCorrectly(t *testing.T) {
    pt := NewTile(tile)

    expected := []uint8 {
        0, 1, 0, 0, 0, 0, 0, 3,
        1, 1, 0, 0, 0, 0, 3, 0,
        0, 1, 0, 0, 0, 3, 0, 0,
        0, 1, 0, 0, 3, 0, 0, 0,
        0, 0, 0, 3, 0, 2, 2, 0,
        0, 0, 3, 0, 0, 0, 0, 2,
        0, 3, 0, 0, 0, 0, 2, 0,
        3, 0, 0, 0, 0, 2, 2, 2,
    }

    for i:=0; i < len(expected); i++ {
        y := uint(i / 8)
        x := uint(i % 8)

        assert.Equal(t, pt.Pixel(x, y), expected[i])
    }
}

func TestPatternTableSelectsRightTile(t *testing.T) {
    p := NewPatterntable()

    p.Write(0xff, 0x0000)
    p.Write(0xef, 0x0010)
    p.Write(0xdf, 0x0140)

    assert.Equal(t, p.Tile(0).raw[0], byte(0xff))
    assert.Equal(t, p.Tile(1).raw[0], byte(0xef))
    assert.Equal(t, p.Tile(20).raw[0], byte(0xdf))
}
