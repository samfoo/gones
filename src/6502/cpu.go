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

func (p *CPU) setInterruptDisable(value bool) {
    p.setFlag(0x04, value)
}

func (p *CPU) setDecimalFlag(value bool) {
    p.setFlag(0x08, value)
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

func addOverflowed(first byte, second byte, result byte) bool {
    if (first & 0x80 == 0x00) &&
        (second & 0x80 == 0x00) &&
        (result & 0x80 == 0x80) {
        // Adding two positives should not result in a negative
        return true
    } else if (first & 0x80 == 0x80) &&
        (second & 0x80 == 0x80) &&
        (result & 0x80 == 0x00) {
        // Adding two negatives should not result in a positive
        return true
    }

    return false
}

func subtractOverflowed(first byte, second byte, result byte) bool {
    if (first & 0x80 == 0x00) &&
        (second & 0x80 == 0x80) &&
        (result & 0x80 == 0x80) {
        // Subtracting a negative from a positive shouldn't result in a
        // negative
        return true
    } else if (first & 0x80 == 0x80) &&
        (second & 0x80 == 0x00) &&
        (result & 0x80 == 0x00) {
        // Subtracting a positive from a negative shoudn't result in a
        // positive
        return true
    }

    return false
}

func (p *CPU) Carry() bool {
    return p.Flags & 0x01 == 0x01
}

func (p *CPU) Zero() bool {
    return p.Flags & 0x02 == 0x02
}

func (p *CPU) InterruptDisable() bool {
    return p.Flags & 0x04 == 0x04
}

func (p *CPU) Decimal() bool {
    return p.Flags & 0x08 == 0x08
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

    if p.Carry() {
        p.A += 0x01
    }

    if p.A < old {
        p.setCarryFlag(true)
    }

    if addOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    }

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Sbc(location Address) {
    other := p.Memory[location]
    old := p.A

    p.A -= other

    if p.Carry() {
        p.A -= 0x01
    }

    if p.A > old {
        p.setCarryFlag(true)
    }

    if subtractOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    }

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) And(location Address) {
    other := p.Memory[location]
    p.A &= other

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Ora(location Address) {
    other := p.Memory[location]
    p.A |= other

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

func (p *CPU) Nop() {
}

func (p *CPU) push(value byte) {
    p.Memory[0x0100 + Address(p.SP)] = value

    p.SP -= 1
}

func (p *CPU) pull(memory *byte) {
    p.SP += 1

    *memory = p.Memory[0x0100 + Address(p.SP)]
}

func (p *CPU) Pha() {
    p.push(p.A)
}

func (p *CPU) Php() {
    p.push(p.Flags | 0x10)
}

func (p *CPU) Pla() {
    p.pull(&p.A)

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Plp() {
    p.pull(&p.Flags)
}

func (p *CPU) Rol(memory *byte) {
    p.setCarryFlag((*memory & 0x80) == 0x80)

    *memory = *memory << 1

    if p.Carry() {
        *memory |= 0x01
    }

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Ror(memory *byte) {
    p.setCarryFlag((*memory & 0x01) == 0x01)

    *memory = *memory >> 1

    if p.Carry() {
        *memory |= 0x80
    }

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Sec() {
    p.setCarryFlag(true)
}

func (p *CPU) Sed() {
    p.setDecimalFlag(true)
}

func (p *CPU) Sei() {
    p.setInterruptDisable(true)
}

func (p *CPU) Sta(location Address) {
    p.Memory[location] = p.A
}

func (p *CPU) Stx(location Address) {
    p.Memory[location] = p.X
}

func (p *CPU) Sty(location Address) {
    p.Memory[location] = p.Y
}

func (p *CPU) Tax() {
    p.X = p.A

    p.setNegativeAndZeroFlags(p.X)
}

func (p *CPU) Tay() {
    p.Y = p.A

    p.setNegativeAndZeroFlags(p.Y)
}

func (p *CPU) Tsx() {
    p.X = p.SP

    p.setNegativeAndZeroFlags(p.X)
}

func (p *CPU) Txa() {
    p.A = p.X

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Txs() {
    p.SP = p.X
}

func (p *CPU) Tya() {
    p.A = p.Y

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Bcc(location Address) {
    if p.Carry() {
        p.PC += 1
    } else {
        p.PC = location
    }
}

func (p *CPU) Bcs(location Address) {
    if p.Carry() {
        p.PC = location
    } else {
        p.PC += 1
    }
}

func (p *CPU) Beq(location Address) {
    if p.Zero() {
        p.PC = location
    } else {
        p.PC += 1
    }
}

func (p *CPU) Bmi(location Address) {
    if p.Negative() {
        p.PC = location
    } else {
        p.PC += 1
    }
}

func (p *CPU) Bne(location Address) {
    if p.Zero() {
        p.PC += 1
    } else {
        p.PC = location
    }
}

func (p *CPU) Bpl(location Address) {
    if p.Negative() {
        p.PC += 1
    } else {
        p.PC = location
    }
}

func (p *CPU) Bvc(location Address) {
    if p.Overflow() {
        p.PC += 1
    } else {
        p.PC = location
    }
}

func (p *CPU) Bvs(location Address) {
    if p.Overflow() {
        p.PC = location
    } else {
        p.PC += 1
    }
}

func (p *CPU) Brk() {
    p.PC += 1

    p.push(byte(p.PC >> 8))
    p.push(byte(p.PC & 0x00ff))

    p.push(p.Flags | 0x10)
    p.setInterruptDisable(true)

    p.PC = Address(p.Memory[0xffff]) << 8
    p.PC |= Address(p.Memory[0xfffe])
}
