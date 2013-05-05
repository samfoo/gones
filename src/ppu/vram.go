package ppu

import "cpu"

type VRAM struct {
    buffer []byte
}

func NewVRAM() *VRAM {
    r := new(VRAM)

    r.buffer = make([]byte, 0x4000)

    return r
}

func (r *VRAM) normalize(location cpu.Address) cpu.Address {
    if location >= 0x0020 && location < 0x0100 {
        return location & 0x1f
    } else {
        return location
    }
}

func (r *VRAM) Write(value byte, location cpu.Address) {
    r.buffer[r.normalize(location)] = value
}

func (r *VRAM) Read(location cpu.Address) byte {
    return r.buffer[r.normalize(location)]
}

