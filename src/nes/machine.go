package nes

import (
    "cpu"
    "ppu"
)

type Machine struct {
    CPU *cpu.CPU
    PPU *ppu.PPU
}

func NewMachine() *Machine {
    m := new(Machine)

    m.CPU = cpu.NewCPU()
    m.PPU = ppu.NewPPU()

    m.CPU.Memory.Mount(m.PPU, 0x2000, 0x3fff)

    // Swallow anything to the APU right now
    // TODO: Mount a real APU here.
    m.CPU.Memory.Mount(cpu.NewRAM(0x0020), 0x4000, 0x401f)

    // At this point all of the memory mapped devices are mounted except the
    // cartridge. The caller should figure out how to do that?

    return m
}
