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
    return Address(p.Memory.Read(p.PC))
}

func (p *CPU) ZeroPageX() Address {
    return Address(p.Memory.Read(p.PC) + p.X)
}

func (p *CPU) ZeroPageY() Address {
    return Address(p.Memory.Read(p.PC) + p.Y)
}

func (p *CPU) absolute() Address {
    high := p.Memory.Read(p.PC+1)
    low := p.Memory.Read(p.PC)

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) Absolute() Address {
    return p.absolute()
}

func (p *CPU) AbsoluteX() Address {
    return p.absolute() + Address(p.X)
}

func (p *CPU) AbsoluteY() Address {
    return p.absolute() + Address(p.Y)
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
    pointer := p.Memory.Read(p.PC) + p.X

    high := p.Memory.Read(Address(pointer+1))
    low := p.Memory.Read(Address(pointer))

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) IndirectIndexed() Address {
    indirect := p.Memory.Read(p.PC)

    high := p.Memory.Read(Address(indirect+1))
    low := p.Memory.Read(Address(indirect))

    return ((Address(high) << 8) + Address(low)) + Address(p.Y)
}

