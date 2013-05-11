package ppu

import "cpu"

type Nametable struct {
    buffer []byte
}

func NewNametable() *Nametable {
    n := new(Nametable)

    n.buffer = make([]byte, 0x400)

    return n
}

func (n *Nametable) TileIndex(x int, y int) uint8 {
    offset := cpu.Address(x + y * 32)
    return n.Read(offset)
}

func (n *Nametable) Attribute(x int, y int) uint8 {
    offset := cpu.Address(y / 4 * 8 + x / 4 + 0x3c0)

    attribs := n.Read(offset)

    xt := (x / 2) % 2
    yt := (y / 2) % 2

    bit_index := (yt * 2) + xt
    return (attribs >> uint(bit_index * 2)) & 0x3
}

func (n *Nametable) Read(location cpu.Address) byte {
    return n.buffer[location]
}

func (n *Nametable) Write(val byte, location cpu.Address) {
    n.buffer[location] = val
}
