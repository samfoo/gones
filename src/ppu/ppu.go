package ppu

import "cpu"

type Address uint16

type Register interface {
    Set(byte)
}

type Ctrl struct {
    BaseNametableAddress Address
    VRAMAddressInc uint8
    SpriteTableAddress Address
    BackgroundTableAddress Address
    SpriteSize uint8
    GenerateNMIOnVBlank bool
}

func (c *Ctrl) Set(val byte) {
    c.BaseNametableAddress = 0x2000 + 0x400 * Address(val & 0x03)
    c.VRAMAddressInc = (val & 0x04) >> 2
    c.SpriteTableAddress = 0x1000 * Address((val & 0x08) >> 3)
    c.BackgroundTableAddress = 0x1000 * Address((val & 0x10) >> 4)
    c.SpriteSize = (val & 0x20) >> 5
    c.GenerateNMIOnVBlank = (val & 0x80) == 0x80
}

type Masks struct {
    Grayscale bool
    ShowBackgroundLeft bool
    ShowSpritesLeft bool
    ShowBackground bool
    ShowSprites bool
    IntenseReds bool
    IntenseGreens bool
    IntenseBlues bool
}

func (m *Masks) Set(val byte) {
    m.Grayscale = (val & 0x01) == 0x01
    m.ShowBackgroundLeft = (val & 0x02) == 0x02
    m.ShowSpritesLeft = (val & 0x04) == 0x04
    m.ShowBackground = (val & 0x08) == 0x08
    m.ShowSprites = (val & 0x10) == 0x10
    m.IntenseReds = (val & 0x20) == 0x20
    m.IntenseGreens = (val & 0x40) == 0x40
    m.IntenseBlues = (val & 0x80) == 0x80
}

type PPU struct {
    Ctrl Register
    Masks Register

    VRAMAddr cpu.Address
    cpu.Memory

    flip bool
}

func NewPPU() *PPU {
    p := new(PPU)

    p.Memory = *cpu.NewMemory()
    p.Memory.Mount(NewInternalRAM(), 0x0000, 0x3fff)

    return p
}

func (p *PPU) SetAddr(val byte) {
    if !p.flip {
        p.VRAMAddr = cpu.Address(val) << 8 | (0x00ff & p.VRAMAddr)
    } else {
        p.VRAMAddr = cpu.Address(val) | (0xff00 & p.VRAMAddr)
    }

    p.flip = !p.flip
}

const (
    PPUCTRL = 0x0000
    PPUMASK = 0x0001
    OAMADDR = 0x0003
    OAMDATA = 0x0004
    PPUSCROLL = 0x0005
    PPUADDR = 0x0006
    PPUDATA = 0x0007
)

func (p *PPU) Write(val byte, location cpu.Address) {
    switch location {
        case PPUCTRL:
            p.Ctrl.Set(val)
        case PPUMASK:
            p.Masks.Set(val)
    }
}
