package ppu

import "cpu"

type Tile struct {
    raw []byte
}

func NewTile(raw []byte) *Tile {
    t := new(Tile)
    t.raw = raw

    return t
}

func (t *Tile) Pixel(x uint, y uint) uint8 {
    first := t.raw[y]
    second := t.raw[y+8]

    index := 7 - x
    var pixel = (first >> index) & 0x1
    pixel |= ((second >> index) & 0x1) << 1

    return pixel
}

type Patterntable struct {
    buffer []byte
}

func NewPatterntable(buffer []byte) *Patterntable {
    p := new(Patterntable)

    p.buffer = buffer

    return p
}

func (p *Patterntable) Tile(offset uint) *Tile {
    slice := p.buffer[offset*16:offset*16+16]
    return NewTile(slice)
}

func (p *Patterntable) Read(location cpu.Address) byte {
    return p.buffer[location]
}

func (p *Patterntable) Write(val byte, location cpu.Address) {
    p.buffer[location] = val
}
