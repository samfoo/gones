package nes

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

type RomData struct {
    buffer []byte
}

func min(first int, second int) int {
    if second < first {
        return second
    }

    return first
}

func (r *RomData) Read(p []byte) (n int, err error) {
    expected := min(len(p), len(r.buffer))

    copy(p[:], r.buffer[:expected])
    n = expected

    return
}

func TestReadHeaderMarker(t *testing.T) {
    var rom = new(RomData)
                     // N     E     S     EOF
    rom.buffer = []byte{0x4E, 0x45, 0x53, 0x1A}

    var _, err = ReadHeader(rom)
    assert.Nil(t, err)

    rom.buffer = []byte{0xde, 0xad, 0xbe, 0xef}

    _, err = ReadHeader(rom)
    assert.NotNil(t, err)
}

var example = []byte{0x4e, 0x45, 0x53, 0x1a, 0x10, 0xfe, 0x10}

func TestReadHeaderSizeofPrg(t *testing.T) {
    var rom = new(RomData)
    rom.buffer = example
    header, _ := ReadHeader(rom)

    assert.Equal(t, header.PrgRomSize, 0x10 * 16384)
}

func TestReadHeaderSizeOfChr(t *testing.T) {
    var rom = new(RomData)
    rom.buffer = example
    header, _ := ReadHeader(rom)

    assert.Equal(t, header.ChrRomSize, 0xfe * 8192)
}

func TestReadHeaderMapper(t *testing.T) {
    var rom = new(RomData)
    rom.buffer = example
    header, _ := ReadHeader(rom)

    assert.Equal(t, header.Mapper, uint8(1))
}
