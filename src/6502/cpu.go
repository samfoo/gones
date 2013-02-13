package cpu

type Address uint16

type CPU struct {
    A, X, Y, SP, Flags byte
    PC Address
    Memory [0x10000]byte
}

func (p *CPU) setFlag(mask byte, value bool) {
    if value {
        p.Flags |= mask
    } else {
        p.Flags &= ^mask
    }
}

func (p *CPU) setCarryFlag(value bool) {
    p.setFlag(0x01, value)
}

func (p *CPU) setZeroFlag(value bool) {
    p.setFlag(0x02, value)
}

func (p *CPU) setOverflowFlag(value bool) {
    p.setFlag(0x40, value)
}

func (p *CPU) setNegativeFlag(value bool) {
    p.setFlag(0x80, value)
}

func addOverflowed(first byte, second byte, result byte) bool {
    if (first & 0x80 == 0x00) &&
        (second & 0x80 == 0x00) &&
        (result & 0x80 == 0x80) {
        return true
    } else if (first & 0x80 == 0x80) &&
        (second & 0x80 == 0x80) &&
        (result & 0x80 == 0x00) {
        return true
    }

    return false
}

func (p *CPU) Zero() bool {
    return p.Flags & 0x02 == 0x02
}

func (p *CPU) Carry() bool {
    return p.Flags & 0x01 == 0x01
}

func (p *CPU) Overflow() bool {
    return p.Flags & 0x40 == 0x40
}

func (p *CPU) Negative() bool {
    return p.Flags & 0x80 == 0x80
}

func (p *CPU) Adc(location Address) {
    other := p.Memory[location]
    old := p.A

    p.A += other

    if p.A < old {
        p.setCarryFlag(true)
    }

    if p.A == 0x00 {
        p.setZeroFlag(true)
    }

    if addOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    }

    if p.A & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }
}

func (p *CPU) And(location Address) {
    other := p.Memory[location]
    p.A &= other

    if p.A == 0x00 {
        p.setZeroFlag(true)
    }

    if p.A & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }
}

func (p *CPU) Asl(memory *byte) {
    if *memory & 0x80 == 0x80 {
        p.setCarryFlag(true)
    }

    *memory = *memory << 1

    if *memory & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }

    if *memory == 0x00 {
        p.setZeroFlag(true)
    }
}

func (p *CPU) Bit(location Address) {
    result := p.A & p.Memory[location]

    if result & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }

    if result & 0x40 == 0x40 {
        p.setOverflowFlag(true)
    }

    if result == 0x00 {
        p.setZeroFlag(true)
    }
}
