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

    // Mount Battery Backed Save or Work RAM
    // TODO: Do some mappers do something with this?
    m.CPU.Memory.Mount(cpu.NewRAM(0x2000), 0x6000, 0x7fff)

    // Setup the interrupt bus to call methods on the CPU
    m.PPU.Bus = m.CPU

    // At this point all of the memory mapped devices are mounted except the
    // cartridge. The caller should figure out how to do that?

    return m
}

func (m *Machine) Insert(rom *ROM) {
    m.PPU.Memory.Mount(rom.Mapper.Graphics(), 0x0000, 0x1fff)
    m.CPU.Memory.Mount(rom.Mapper.Program(), 0x8000, 0xffff)

    m.CPU.PC = cpu.Address(m.CPU.Memory.Read(0xFFFC)) |
        (cpu.Address(m.CPU.Memory.Read(0xFFFD))<<8)
}
