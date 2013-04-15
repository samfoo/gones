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
    // TODO: Mirroing of 3000-3f00 <-> 2000-2EFF
    return location
}

func (r *InternalRAM) Write(value byte, location cpu.Address) {
    r.buffer[r.normalize(location)] = value
}

func (r *InternalRAM) Read(location cpu.Address) byte {
    return r.buffer[r.normalize(location)]
}

