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
}
