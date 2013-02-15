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

func (p *CPU) setNegativeAndZeroFlags(value byte) {
    if value & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }

    if value == 0x00 {
        p.setZeroFlag(true)
    }
}

func (p *CPU) setInterruptDisable(value bool) {
    p.setFlag(0x04, value)
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

func (p *CPU) InterruptDisable() bool {
    return p.Flags & 0x04 == 0x04
}

func (p *CPU) Adc(location Address) {
    other := p.Memory[location]
    old := p.A

    p.A += other

    if p.A < old {
        p.setCarryFlag(true)
    }

    if addOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    }

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) And(location Address) {
    other := p.Memory[location]
    p.A &= other

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Asl(memory *byte) {
    if *memory & 0x80 == 0x80 {
        p.setCarryFlag(true)
    }

    *memory = *memory << 1

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Bit(location Address) {
    result := p.A & p.Memory[location]

    if result & 0x40 == 0x40 {
        p.setOverflowFlag(true)
    }

    p.setNegativeAndZeroFlags(result)
}

func (p *CPU) Clc() {
    p.setCarryFlag(false)
}

func (p *CPU) Cli() {
    p.setInterruptDisable(false)
}

func (p *CPU) Clv() {
    p.setOverflowFlag(false)
}

func (p *CPU) compare(register byte, value byte) {
    result := register - value

    if register >= value {
        p.setCarryFlag(true)
    }

    p.setNegativeAndZeroFlags(result)
}

func (p *CPU) Cmp(location Address) {
    p.compare(p.A, p.Memory[location])
}

func (p *CPU) Cpx(location Address) {
    p.compare(p.X, p.Memory[location])
}

func (p *CPU) Cpy(location Address) {
    p.compare(p.Y, p.Memory[location])
}

func (p *CPU) decrement(memory *byte) {
    *memory -= 1

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Dec(location Address) {
    p.decrement(&p.Memory[location])
}

func (p *CPU) Dex() {
    p.decrement(&p.X)
}

func (p *CPU) Dey() {
    p.decrement(&p.Y)
}

func (p *CPU) Eor(location Address) {
    p.A ^= p.Memory[location]

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) increment(memory *byte) {
    *memory += 1

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Inc(location Address) {
    p.increment(&p.Memory[location])
}

func (p *CPU) Inx() {
    p.increment(&p.X)
}

func (p *CPU) Iny() {
    p.increment(&p.Y)
}

func (p *CPU) load(memory *byte, value byte) {
    *memory = value

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Lda(location Address) {
    p.load(&p.A, p.Memory[location])
}

func (p *CPU) Ldx(location Address) {
    p.load(&p.X, p.Memory[location])
}

func (p *CPU) Ldy(location Address) {
    p.load(&p.Y, p.Memory[location])
}

func (p *CPU) Lsr(memory *byte) {
    if *memory & 0x01 == 0x01 {
        p.setCarryFlag(true)
    }

    *memory = *memory >> 1

    p.setNegativeAndZeroFlags(*memory)
}
