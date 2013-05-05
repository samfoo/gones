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

const (
    VRAM_INC_ACROSS = uint8(0x00)
    VRAM_INC_DOWN = uint8(0x01)
)

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

type Status struct {
    SpriteOverflow bool
    Sprite0Hit bool
    VBlankStarted bool
}

func (s *Status) Value() byte {
    var value = byte(0x00)

    if s.SpriteOverflow { value |= 0x20 }
    if s.Sprite0Hit { value |= 0x40 }
    if s.VBlankStarted { value |= 0x80 }

    return value
}

type PPU struct {
    Ctrl
    Masks

    Status

    VRAMAddr cpu.Address

    OAMAddr uint8
    OAMRAM [0x100]byte

    Memory *cpu.Memory
    Display *Display

    Cycle int
    Frame int
    Scanline int
    AddressLatch bool

    Bus cpu.Bus

    vram *VRAM
}

func NewPPU() *PPU {
    p := new(PPU)

    p.Memory = cpu.NewMemory()
    p.Memory.Mount(NewVRAM(), 0x0000, 0x3fff)

    p.Frame = 0
    p.Scanline = PRERENDER_SCANLINE

    return p
}

func (p *PPU) WriteVRAMAddr(val byte) {
    if p.AddressLatch {
        p.VRAMAddr = cpu.Address(val) << 8 | (0x00ff & p.VRAMAddr)
    } else {
        p.VRAMAddr = cpu.Address(val) | (0xff00 & p.VRAMAddr)
    }

    p.AddressLatch = !p.AddressLatch
}

func (p *PPU) ReadData() byte {
    value := p.Memory.Read(p.VRAMAddr)

    if p.Ctrl.VRAMAddressInc == VRAM_INC_ACROSS {
        p.VRAMAddr += 1
    } else {
        p.VRAMAddr += 32
    }

    return value
}

const (
    PPUCTRL = 0x0000
    PPUMASK = 0x0001
    PPUSTATUS = 0x0002
    OAMADDR = 0x0003
    OAMDATA = 0x0004
    PPUSCROLL = 0x0005
    PPUADDR = 0x0006
    PPUDATA = 0x0007
)

const (
    PRERENDER_SCANLINE = -1
    VISIBLE_SCANLINES = 240
    POSTRENDER_SCANLINE = 240
    VBLANK_SCANLINE = 260

    LAST_CYCLE = 341
)

func (p *PPU) GenerateNMI() {
    if p.Bus != nil && p.Ctrl.GenerateNMIOnVBlank && p.Status.VBlankStarted {
        p.Bus.Interrupt(cpu.NMI)
    }
}

func (p *PPU) Step() {
    if p.Scanline == VBLANK_SCANLINE && p.Cycle == LAST_CYCLE {
        p.Frame++
        p.Cycle = 1
        p.Scanline = PRERENDER_SCANLINE
    } else {
        switch {
            case p.Scanline == PRERENDER_SCANLINE && p.Cycle == 1:
                p.Status.SpriteOverflow = false
                p.Status.Sprite0Hit = false
                p.Status.VBlankStarted = false
            case p.Scanline < VISIBLE_SCANLINES:
                // TODO Do rendering
            case p.Scanline == POSTRENDER_SCANLINE + 1 && p.Cycle == 1:
                p.Status.VBlankStarted = true
                p.GenerateNMI()
        }

        if p.Cycle == LAST_CYCLE {
            p.Cycle = 1
            p.Scanline++
        } else if p.Scanline == PRERENDER_SCANLINE &&
            p.Cycle == LAST_CYCLE - 1 &&
            p.shortPrerender() {
            p.Cycle = 1
            p.Scanline++
        } else {
            p.Cycle++
        }
    }
}

func (p *PPU) shortPrerender() bool {
    return p.Frame % 2 == 1 &&
        (p.Masks.ShowBackground || p.Masks.ShowSprites)
}

func (p *PPU) normalize(location cpu.Address) cpu.Address {
    return location & 0x7
}

func (p *PPU) Write(val byte, location cpu.Address) {
    switch p.normalize(location) {
        case PPUCTRL:
            generateAlreadySet := p.Ctrl.GenerateNMIOnVBlank
            p.Ctrl.Set(val)

            if !generateAlreadySet {
                p.GenerateNMI()
            }
        case PPUMASK:
            p.Masks.Set(val)
        case OAMADDR:
            p.OAMAddr = val
        case OAMDATA:
            p.OAMRAM[p.OAMAddr] = val
            p.OAMAddr++
        case PPUSCROLL:
            // TODO
        case PPUADDR:
            p.WriteVRAMAddr(val)
        case PPUDATA:
            p.Memory.Write(val, p.VRAMAddr)
            p.VRAMAddr++
    }
}

func (p *PPU) ReadDebug(location cpu.Address) byte {
    switch p.normalize(location) {
        case PPUSTATUS:
            return p.Status.Value()
        case OAMDATA:
            return p.OAMRAM[p.OAMAddr]
        case PPUDATA:
            return p.Memory.ReadDebug(p.VRAMAddr)
        default:
            return 0
    }
}

func (p *PPU) Read(location cpu.Address) byte {
    switch p.normalize(location) {
        case PPUSTATUS:
            p.AddressLatch = true
            serialized := p.Status.Value()
            p.Status.VBlankStarted = false

            return serialized
        case OAMDATA:
            return p.OAMRAM[p.OAMAddr]
        case PPUDATA:
            return p.ReadData()
        default:
            return 0
    }
}
