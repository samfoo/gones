package ppu

import "cpu"

type Address uint16

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

type Scroll struct {
    X uint8
    Y uint8
    flip bool
}

func (s *Scroll) Set(val uint8) {
    if !s.flip {
        s.X = val
    } else {
        s.Y = val
    }

    s.flip = !s.flip
}

type Addr struct {
    Location cpu.Address
    flip bool
}

func (a *Addr) Set(val byte) {
    if !a.flip {
        a.Location = cpu.Address(val) << 8 | (0x00ff & a.Location)
    } else {
        a.Location = cpu.Address(val) | (0xff00 & a.Location)
    }

    a.flip = !a.flip
}

type PPU struct {
    Ctrl
    Masks
    Scroll
    cpu.Memory
}

func NewPPU() *PPU {
    p := new(PPU)

    p.Memory = *cpu.NewMemory()
    p.Memory.Mount(NewInternalRAM(), 0x0000, 0x3fff)

    return p
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
}
