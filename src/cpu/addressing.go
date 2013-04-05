package cpu

type AddressMode func(*CPU) Address

type Reader interface {
    Read(Address) byte
}

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

func (p *CPU) relative(r Reader) Address {
    var offset = Address(r.Read(p.PC))
    if offset < 0x0080 {
        offset += p.PC + 1
    } else {
        offset += (p.PC - 0x100) + 1
    }

    return offset
}

func (p *CPU) Relative() Address {
    return p.relative(p)
}

func (p *CPU) Immediate() Address {
    return p.PC
}

func (p *CPU) zeroPage(r Reader) Address {
    return Address(r.Read(p.PC))
}

func (p *CPU) ZeroPage() Address {
    return p.zeroPage(p)
}

func (p *CPU) zeroPageX(r Reader) Address {
    addr := r.Read(p.PC)

    return Address(addr + p.X)
}

func (p *CPU) ZeroPageX() Address {
    // zpx does an extra read of the pre-x address, which cycles an extra time.
    // See -- http://nemulator.com/files/nes_emu.txt
    p.cycles++

    return p.zeroPageX(p)
}

func (p *CPU) zeroPageY(r Reader) Address {
    addr := r.Read(p.PC)

    return Address(addr + p.Y)
}

func (p *CPU) ZeroPageY() Address {
    // zpy does an extra read of the pre-x address, which cycles an extra time.
    // See -- http://nemulator.com/files/nes_emu.txt
    p.cycles++

    return p.zeroPageY(p)
}

func (p *CPU) absolute(r Reader) Address {
    high := r.Read(p.PC+1)
    low := r.Read(p.PC)

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) Absolute() Address {
    return p.absolute(p)
}

func (p *CPU) AbsoluteX() Address {
    abs := p.absolute(p)
    addr := abs + Address(p.X)

    // absx does an extra read if the address crosses a page boundary.
    if (abs & 0x100) ^ (addr & 0x100) != 0 {
        p.cycles++
    }

    return addr
}

func (p *CPU) AbsoluteY() Address {
    abs := p.absolute(p)
    addr := abs + Address(p.Y)

    // absx does an extra read if the address crosses a page boundary.
    if (abs & 0x100) ^ (addr & 0x100) != 0 {
        p.cycles++
    }

    return addr
}

func (p *CPU) indirect(r Reader) Address {
    location := p.absolute(r)

    low := r.Read(location)

    var high byte
    if location & 0x00ff == 0x00ff {
        high = r.Read(location & 0xff00)
    } else {
        high = r.Read(location+1)
    }

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) Indirect() Address {
    return p.indirect(p)
}

func (p *CPU) indexedIndirect(r Reader) Address {
    pointer := r.Read(p.PC)

    high := r.Read(Address(pointer+p.X+1))
    low := r.Read(Address(pointer+p.X))

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) IndexedIndirect() Address {
    // indx does an extra read pointer, which cycles an extra time.
    // See -- http://nemulator.com/files/nes_emu.txt
    p.cycles++

    return p.indexedIndirect(p)
}

func (p *CPU) indirectIndexed(r Reader) Address {
    indirect := r.Read(p.PC)

    high := r.Read(Address(indirect+1))
    low := r.Read(Address(indirect))

    ind := ((Address(high) << 8) + Address(low))
    addr := ind + Address(p.Y)

    // indirect idx does an extra read if the address crosses a page boundary.
    if (ind & 0x100) ^ (addr & 0x100) != 0 {
        p.cycles++
    }

    return addr
}

func (p *CPU) IndirectIndexed() Address {
    return p.indirectIndexed(p)
}
