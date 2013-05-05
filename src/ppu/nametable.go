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

func (n *Nametable) Read(location cpu.Address) byte {
    return n.buffer[location]
}

func (n *Nametable) Write(val byte, location cpu.Address) {
    n.buffer[location] = val
}
