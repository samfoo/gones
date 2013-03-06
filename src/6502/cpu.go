package cpu

type Address uint16

type CPU struct {
    A, X, Y, SP, Flags byte
    PC Address
    Memory [0x10000]byte
}

type Opcode byte
type Operation func(*CPU)

var opcodes = map[Opcode]Operation {
    0x69: func(p *CPU) {
        p.Adc(p.Address((*CPU).Immediate))
    },
    0x65: func(p *CPU) {
        p.Adc(p.Address((*CPU).ZeroPage))
    },
    0x75: func(p *CPU) {
        p.Adc(p.Address((*CPU).ZeroPageX))
    },
    0x6d: func(p *CPU) {
        p.Adc(p.Address((*CPU).Absolute))
    },
    0x7d: func(p *CPU) {
        p.Adc(p.Address((*CPU).AbsoluteX))
    },
    0x79: func(p *CPU) {
        p.Adc(p.Address((*CPU).AbsoluteY))
    },
    0x61: func(p *CPU) {
        p.Adc(p.Address((*CPU).IndexedIndirect))
    },
    0x71: func(p *CPU) {
        p.Adc(p.Address((*CPU).IndirectIndexed))
    },
    0x29: func(p *CPU) {
        p.And(p.Address((*CPU).Immediate))
    },
    0x25: func(p *CPU) {
        p.And(p.Address((*CPU).ZeroPage))
    },
    0x35: func(p *CPU) {
        p.And(p.Address((*CPU).ZeroPageX))
    },
    0x2d: func(p *CPU) {
        p.And(p.Address((*CPU).Absolute))
    },
    0x3d: func(p *CPU) {
        p.And(p.Address((*CPU).AbsoluteX))
    },
    0x39: func(p *CPU) {
        p.And(p.Address((*CPU).AbsoluteY))
    },
    0x21: func(p *CPU) {
        p.And(p.Address((*CPU).IndexedIndirect))
    },
    0x31: func(p *CPU) {
        p.And(p.Address((*CPU).IndirectIndexed))
    },
    0x0a: func(p *CPU) {
        p.Asl(&p.A)
    },
    0x06: func(p *CPU) {
        address := p.Address((*CPU).ZeroPage)
        p.Asl(&p.Memory[address])
    },
    0x16: func(p *CPU) {
        address := p.Address((*CPU).ZeroPageX)
        p.Asl(&p.Memory[address])
    },
    0x0e: func(p *CPU) {
        address := p.Address((*CPU).Absolute)
        p.Asl(&p.Memory[address])
    },
    0x1e: func(p *CPU) {
        address := p.Address((*CPU).AbsoluteX)
        p.Asl(&p.Memory[address])
    },
    0x90: func(p *CPU) {
        p.Bcc(p.Address((*CPU).Relative))
    },
    0xb0: func(p *CPU) {
        p.Bcs(p.Address((*CPU).Relative))
    },
    0xf0: func(p *CPU) {
        p.Beq(p.Address((*CPU).Relative))
    },
    0xd0: func(p *CPU) {
        p.Bne(p.Address((*CPU).Relative))
    },
    0x24: func(p *CPU) {
        p.Bit(p.Address((*CPU).ZeroPage))
    },
    0x2c: func(p *CPU) {
        p.Bit(p.Address((*CPU).Absolute))
    },
    0x30: func(p *CPU) {
        p.Bmi(p.Address((*CPU).Relative))
    },
    0x10: func(p *CPU) {
        p.Bpl(p.Address((*CPU).Relative))
    },
    0x00: func(p *CPU) {
        p.Brk()
    },
    0x50: func(p *CPU) {
        p.Bvc(p.Address((*CPU).Relative))
    },
    0x70: func(p *CPU) {
        p.Bvs(p.Address((*CPU).Relative))
    },
    0x18: func(p *CPU) {
        p.Clc()
    },
    0xd8: func(p *CPU) {
        p.Cld()
    },
    0x58: func(p *CPU) {
        p.Cli()
    },
    0xb8: func(p *CPU) {
        p.Clv()
    },
    0xc9: func(p *CPU) {
        p.Cmp(p.Address((*CPU).Immediate))
    },
    0xc5: func(p *CPU) {
        p.Cmp(p.Address((*CPU).ZeroPage))
    },
    0xd5: func(p *CPU) {
        p.Cmp(p.Address((*CPU).ZeroPageX))
    },
    0xcd: func(p *CPU) {
        p.Cmp(p.Address((*CPU).Absolute))
    },
    0xdd: func(p *CPU) {
        p.Cmp(p.Address((*CPU).AbsoluteX))
    },
    0xd9: func(p *CPU) {
        p.Cmp(p.Address((*CPU).AbsoluteY))
    },
    0xc1: func(p *CPU) {
        p.Cmp(p.Address((*CPU).IndexedIndirect))
    },
    0xd1: func(p *CPU) {
        p.Cmp(p.Address((*CPU).IndirectIndexed))
    },
}

func (p *CPU) Op(opcode Opcode) func() {
    return func() {
        opcodes[opcode](p)
    }
}

func (p *CPU) Reset() {
    p.Flags = 0x34
    p.A, p.X, p.Y = 0x00, 0x00, 0x00
    p.SP = 0xfd
}

type AddressMode func(*CPU) (Address, int)

func (p *CPU) Address(mode AddressMode) Address {
    var location, i = mode(p)

    p.PC += Address(i)

    return location
}

func (p *CPU) Relative() (Address, int) {
    var offset = Address(p.Memory[p.PC])
    if offset < 0x80 {
        offset += p.PC + 1
    } else {
        offset += (p.PC - 0x100) + 1
    }

    return offset, 1
}

func (p *CPU) Immediate() (Address, int) {
    return p.PC, 1
}

func (p *CPU) ZeroPage() (Address, int) {
    return Address(p.Memory[p.PC]), 1
}

func (p *CPU) ZeroPageX() (Address, int) {
    return Address(p.Memory[p.PC] + p.X), 1
}

func (p *CPU) absolute() Address {
    high := p.Memory[p.PC+1]
    low := p.Memory[p.PC]

    return (Address(high) << 8) + Address(low)
}

func (p *CPU) Absolute() (Address, int) {
    return p.absolute(), 2
}

func (p *CPU) AbsoluteX() (Address, int) {
    return p.absolute() + Address(p.X), 2
}

func (p *CPU) AbsoluteY() (Address, int) {
    return p.absolute() + Address(p.Y), 2
}

func (p *CPU) IndexedIndirect() (Address, int) {
    pointer := p.Memory[p.PC] + p.X

    high := p.Memory[pointer+1]
    low := p.Memory[pointer]

    return (Address(high) << 8) + Address(low), 1
}

func (p *CPU) IndirectIndexed() (Address, int) {
    high := p.Memory[p.PC+1]
    low := p.Memory[p.PC]

    return (Address(high) << 8) + Address(low) + Address(p.Y), 1
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

    if result == 0x00 {
        p.setZeroFlag(true)
    }

    if p.Memory[location] & 0x40 == 0x40 {
        p.setOverflowFlag(true)
    }

    if p.Memory[location] & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    }
}

func (p *CPU) Clc() {
    p.setCarryFlag(false)
}

func (p *CPU) Cld() {
    p.setDecimalFlag(false)
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
    if !p.Carry() {
        p.PC = location
    }
}

func (p *CPU) Bcs(location Address) {
    if p.Carry() {
        p.PC = location
    }
}

func (p *CPU) Beq(location Address) {
    if p.Zero() {
        p.PC = location
    }
}

func (p *CPU) Bmi(location Address) {
    if p.Negative() {
        p.PC = location
    }
}

func (p *CPU) Bne(location Address) {
    if !p.Zero() {
        p.PC = location
    }
}

func (p *CPU) Bpl(location Address) {
    if !p.Negative() {
        p.PC = location
    }
}

func (p *CPU) Bvc(location Address) {
    if !p.Overflow() {
        p.PC = location
    }
}

func (p *CPU) Bvs(location Address) {
    if p.Overflow() {
        p.PC = location
    }
}

func (p *CPU) Brk() {
    p.push(byte((p.PC+1) >> 8))
    p.push(byte((p.PC+1) & 0x00ff))

    p.push(p.Flags | 0x10)
    p.setInterruptDisable(true)

    p.PC = Address(p.Memory[0xffff]) << 8
    p.PC |= Address(p.Memory[0xfffe])
}

func (p *CPU) Jmp(location Address) {
    p.PC = location
}

func (p *CPU) Jsr(location Address) {
    p.push(byte((p.PC-1) >> 8))
    p.push(byte((p.PC-1) & 0x00ff))

    p.PC = location
}

func (p *CPU) Rti() {
    p.pull(&p.Flags)

    var low byte = 0x00
    p.pull(&low)
    var high byte = 0x00
    p.pull(&high)

    p.PC = (Address(high) << 8) + Address(low)
}

func (p *CPU) Rts() {
    var low byte = 0x00
    p.pull(&low)
    var high byte = 0x00
    p.pull(&high)

    p.PC = (Address(high) << 8) + Address(low) + 1
}
