package cpu

type AddressMode func(*CPU) Address

const (
    Immediate = iota
    ZeroPage
    ZeroPageX
    ZeroPageY
    IndexedIndirect
    IndirectIndexed
    Absolute
    AbsoluteX
    AbsoluteY
    Indirect
    Relative
    Accumulator
    Implied
)


var addressing = map[int]AddressMode {
    Immediate: (*CPU).Immediate,
    ZeroPage: (*CPU).ZeroPage,
    ZeroPageX: (*CPU).ZeroPageX,
    ZeroPageY: (*CPU).ZeroPageY,
    IndexedIndirect: (*CPU).IndexedIndirect,
    IndirectIndexed: (*CPU).IndirectIndexed,
    Absolute: (*CPU).Absolute,
    AbsoluteX: (*CPU).AbsoluteX,
    AbsoluteY: (*CPU).AbsoluteY,
    Indirect: (*CPU).Indirect,
    Relative: (*CPU).Relative,
}

func addressSize(mode int) Address {
    switch mode {
        case Absolute, AbsoluteX, AbsoluteY, Indirect:
            return 2
    }

    return 1
}

func (p *CPU) Relative() Address {
    var offset = Address(p.Memory.Read(p.PC))
    if offset < 0x0080 {
        offset += p.PC + 1
    } else {
        offset += (p.PC - 0x100) + 1
    }

    return offset
}

func (p *CPU) Immediate() Address {
    return p.PC
}

func (p *CPU) ZeroPage() Address {
    return Address(p.Read(p.PC))
}

func (p *CPU) ZeroPageX() Address {
    addr := p.Read(p.PC)

    // zpx does an extra read of the pre-x address, which cycles an extra time.
    // See -- http://nemulator.com/files/nes_emu.txt
    p.Read(Address(addr))

    return Address(addr + p.X)
}

func (p *CPU) ZeroPageY() Address {
    return Address(p.Memory.Read(p.PC) + p.Y)
}

func (p *CPU) absolute() Address {
    high := p.Read(p.PC+1)
    low := p.Read(p.PC)

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) Absolute() Address {
    return p.absolute()
}

func (p *CPU) AbsoluteX() Address {
    abs := p.absolute()
    addr := abs + Address(p.X)

    // absx does an extra read if the address crosses a page boundary.
    if (abs & 0x100) ^ (addr & 0x100) != 0 {
        p.Read(Address(addr))
    }

    return addr
}

func (p *CPU) AbsoluteY() Address {
    abs := p.absolute()
    addr := abs + Address(p.Y)

    // absx does an extra read if the address crosses a page boundary.
    if (abs & 0x100) ^ (addr & 0x100) != 0 {
        p.Read(Address(addr))
    }

    return addr
}

func (p *CPU) Indirect() Address {
    location := p.absolute()

    low := p.Memory.Read(location)

    var high byte
    if location & 0x00ff == 0x00ff {
        high = p.Memory.Read(location & 0xff00)
    } else {
        high = p.Memory.Read(location+1)
    }

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) IndexedIndirect() Address {
    pointer := p.Read(p.PC)

    // indx does an extra read pointer, which cycles an extra time.
    // See -- http://nemulator.com/files/nes_emu.txt
    p.Read(Address(pointer))

    high := p.Read(Address(pointer+p.X+1))
    low := p.Read(Address(pointer+p.X))

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) IndirectIndexed() Address {
    indirect := p.Read(p.PC)

    high := p.Read(Address(indirect+1))
    low := p.Read(Address(indirect))

    ind := ((Address(high) << 8) + Address(low))
    addr := ind + Address(p.Y)

    // indirect idx does an extra read if the address crosses a page boundary.
    if (ind & 0x100) ^ (addr & 0x100) != 0 {
        p.Read(Address(addr))
    }

    return addr
}

