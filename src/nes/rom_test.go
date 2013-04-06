package nes

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

var example = []byte{
    0x4e, 0x45, 0x53, 0x1a,
    0x10, 0x00, 0x10, 0x00,
    0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,
}

func TestParseHeaderMarker(t *testing.T) {
    var header = make([]byte, 16)
    copy(header[:4], []byte{0x4E, 0x45, 0x53, 0x1A})

    var _, err = ParseHeader(header)
    assert.Nil(t, err)

    copy(header[:4], []byte{0xde, 0xad, 0xbe, 0xef})

    _, err = ParseHeader(header)
    assert.NotNil(t, err)
}

func TestParseHeaderSizeofPrg(t *testing.T) {
    header, _ := ParseHeader(example)

    assert.Equal(t, header.PrgRomSize, 0x10)
}

func TestParseHeaderSizeOfChr(t *testing.T) {
    header, _ := ParseHeader(example)

    assert.Equal(t, header.ChrRomSize, 0)
}

func TestParseHeaderMapper(t *testing.T) {
    header, _ := ParseHeader(example)

    assert.Equal(t, header.Mapper, uint8(1))
}
