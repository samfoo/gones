package ppu

import "cpu"

type InternalRAM struct {
    buffer []byte
}

func NewInternalRAM() *InternalRAM {
    r := new(InternalRAM)

    r.buffer = make([]byte, 0x4000)

    return r
}

func (r *InternalRAM) normalize(location cpu.Address) cpu.Address {
    if location >= 0x3000 && location < 0x3f00 {
        return location & 0x2fff
    } else if location >= 0x3f20 && location < 0x4000 {
        return location & 0x3f1f
    } else {
        return location
    }
}

func (r *InternalRAM) Write(value byte, location cpu.Address) {
    r.buffer[r.normalize(location)] = value
}

func (r *InternalRAM) Read(location cpu.Address) byte {
    return r.buffer[r.normalize(location)]
}

