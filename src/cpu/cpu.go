package cpu

import "fmt"

type Address uint16

type CPU struct {
    A, X, Y, SP, Flags byte
    PC Address
    Memory [0x10000]byte
    operations map[Opcode]Op
}

type Opcode byte
type AddressMode func(*CPU) (Address, int)

type Op struct {
    Name string
    Method interface{}
    Mode interface{}
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

func (p *CPU) Execute(op Op) {
    switch m := op.Method.(type) {
        case func(*CPU, Address):
            m(p, p.Address(addressing[op.Mode.(int)]))
        case func(*CPU, *byte):
            if op.Mode == nil {
                m(p, &p.A)
            } else {
                location := p.Address(addressing[op.Mode.(int)])
                m(p, &p.Memory[location])
            }
        case func(*CPU):
            m(p)
    }
}

func (p *CPU) Step() {
    opcode := Opcode(p.Memory[p.PC])

    fmt.Printf("%4.04X  ", p.PC)

    p.PC++

    op := p.Operations()[opcode]

    fmt.Printf("%-2.02X ", opcode)

    switch op.Mode {
        case Immediate, ZeroPage, ZeroPageX, IndexedIndirect, IndirectIndexed, Relative:
            fmt.Printf("%-6.02X ", p.Memory[p.PC])
        case Absolute, AbsoluteX, AbsoluteY, Indirect:
            fmt.Printf("%02X %-3.02X ", p.Memory[p.PC], p.Memory[p.PC+1])
        case nil:
            fmt.Printf("%-7s", " ")
    }

    fmt.Printf("%s ", op.Name)

    switch op.Mode {
        case Immediate:
            fmt.Printf("#$%-26.02X", p.Memory[p.PC])
        case ZeroPage:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X = %-22.02X", p.Memory[p.PC], p.Memory[location])
        case ZeroPageX:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X,X @ %02X = %-16.02X", p.Memory[p.PC], location & 0xff, p.Memory[location])
        case ZeroPageY:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X,Y @ %02X = %-16.02X", p.Memory[p.PC], location & 0xff, p.Memory[location])
        case Absolute:
            if op.Name == "JMP" || op.Name == "JSR" {
                fmt.Printf("$%-27.04X", p.absolute())
            } else {
                fmt.Printf("$%04X = %-20.02X", p.absolute(), p.Memory[p.absolute()])
            }
        case Indirect:
            location := p.Memory[p.absolute()]
            high := p.Memory[location+1]
            low := p.Memory[location]
            indirectLocation := (Address(high) << 8) + Address(low)

            fmt.Printf("($%04X) = %-19.04X", p.absolute(), indirectLocation)
        case AbsoluteX:
            fmt.Printf("$%04X,X @ %04X = %-11.02X", p.absolute(), p.absolute() + Address(p.X), p.Memory[p.absolute()+Address(p.X)])
        case AbsoluteY:
            fmt.Printf("$%04Y,Y @ %04Y = %-11.02Y", p.absolute(), p.absolute() + Address(p.Y), p.Memory[p.absolute()+Address(p.Y)])
        case IndexedIndirect:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("($%02X,X) @ %02X = %04X = %-6.02X", p.Memory[p.PC], p.Memory[p.PC] + p.X, location, p.Memory[location])
        case IndirectIndexed:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("($%02X),Y = %04X @ %04X = %02X", p.Memory[p.PC], location - Address(p.Y), location, p.Memory[location])
        case Relative:
            fmt.Printf("$%-27.04X", p.PC+Address(p.Memory[p.PC])+1)
        case nil:
            switch op.Method.(type) {
                case func(*CPU, *byte):
                    // Accumulator addressing.
                    fmt.Printf("%-28s", "A")
                default:
                    // No addressing
                    fmt.Printf("%-28s", " ")
            }
    }

    fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", p.A, p.X, p.Y, p.Flags, p.SP)

    p.Execute(op)
}

func (p *CPU) Operations() map[Opcode]Op {
    if p.operations == nil {
        p.operations = map[Opcode]Op {
            0x69: Op{"ADC", (*CPU).Adc, Immediate},
            0x65: Op{"ADC", (*CPU).Adc, ZeroPage},
            0x75: Op{"ADC", (*CPU).Adc, ZeroPageX},
            0x6d: Op{"ADC", (*CPU).Adc, Absolute},
            0x7d: Op{"ADC", (*CPU).Adc, AbsoluteX},
            0x79: Op{"ADC", (*CPU).Adc, AbsoluteY},
            0x61: Op{"ADC", (*CPU).Adc, IndexedIndirect},
            0x71: Op{"ADC", (*CPU).Adc, IndirectIndexed},
            0x29: Op{"AND", (*CPU).And, Immediate},
            0x25: Op{"AND", (*CPU).And, ZeroPage},
            0x35: Op{"AND", (*CPU).And, ZeroPageX},
            0x2d: Op{"AND", (*CPU).And, Absolute},
            0x3d: Op{"AND", (*CPU).And, AbsoluteX},
            0x39: Op{"AND", (*CPU).And, AbsoluteY},
            0x21: Op{"AND", (*CPU).And, IndexedIndirect},
            0x31: Op{"AND", (*CPU).And, IndirectIndexed},
            0x0a: Op{"ASL", (*CPU).Asl, nil},
            0x06: Op{"ASL", (*CPU).Asl, ZeroPage},
            0x16: Op{"ASL", (*CPU).Asl, ZeroPageX},
            0x0e: Op{"ASL", (*CPU).Asl, Absolute},
            0x1e: Op{"ASL", (*CPU).Asl, AbsoluteX},
            0x90: Op{"BCC", (*CPU).Bcc, Relative},
            0xb0: Op{"BCS", (*CPU).Bcs, Relative},
            0xf0: Op{"BEQ", (*CPU).Beq, Relative},
            0xd0: Op{"BNE", (*CPU).Bne, Relative},
            0x24: Op{"BIT", (*CPU).Bit, ZeroPage},
            0x2c: Op{"BIT", (*CPU).Bit, Absolute},
            0x30: Op{"BMI", (*CPU).Bmi, Relative},
            0x10: Op{"BPL", (*CPU).Bpl, Relative},
            0x00: Op{"BRK", (*CPU).Brk, nil},
            0x50: Op{"BVC", (*CPU).Bvc, Relative},
            0x70: Op{"BVS", (*CPU).Bvs, Relative},
            0x18: Op{"CLC", (*CPU).Clc, nil},
            0xd8: Op{"CLD", (*CPU).Cld, nil},
            0x58: Op{"CLI", (*CPU).Cli, nil},
            0xb8: Op{"CLV", (*CPU).Clv, nil},
            0xc9: Op{"CMP", (*CPU).Cmp, Immediate},
            0xc5: Op{"CMP", (*CPU).Cmp, ZeroPage},
            0xd5: Op{"CMP", (*CPU).Cmp, ZeroPageX},
            0xcd: Op{"CMP", (*CPU).Cmp, Absolute},
            0xdd: Op{"CMP", (*CPU).Cmp, AbsoluteX},
            0xd9: Op{"CMP", (*CPU).Cmp, AbsoluteY},
            0xc1: Op{"CMP", (*CPU).Cmp, IndexedIndirect},
            0xd1: Op{"CMP", (*CPU).Cmp, IndirectIndexed},
            0xe0: Op{"CPX", (*CPU).Cpx, Immediate},
            0xe4: Op{"CPX", (*CPU).Cpx, ZeroPage},
            0xec: Op{"CPX", (*CPU).Cpx, Absolute},
            0xc0: Op{"CPY", (*CPU).Cpy, Immediate},
            0xc4: Op{"CPY", (*CPU).Cpy, ZeroPage},
            0xcc: Op{"CPY", (*CPU).Cpy, Absolute},
            0xc6: Op{"DEC", (*CPU).Dec, ZeroPage},
            0xd6: Op{"DEC", (*CPU).Dec, ZeroPageX},
            0xce: Op{"DEC", (*CPU).Dec, Absolute},
            0xde: Op{"DEC", (*CPU).Dec, AbsoluteX},
            0xca: Op{"DEX", (*CPU).Dex, nil},
            0x88: Op{"DEY", (*CPU).Dey, nil},
            0xe6: Op{"INC", (*CPU).Inc, ZeroPage},
            0xf6: Op{"INC", (*CPU).Inc, ZeroPageX},
            0xee: Op{"INC", (*CPU).Inc, Absolute},
            0xfe: Op{"INC", (*CPU).Inc, AbsoluteX},
            0xe8: Op{"INX", (*CPU).Inx, nil},
            0xc8: Op{"INY", (*CPU).Iny, nil},
            0x49: Op{"EOR", (*CPU).Eor, Immediate},
            0x45: Op{"EOR", (*CPU).Eor, ZeroPage},
            0x55: Op{"EOR", (*CPU).Eor, ZeroPageX},
            0x4d: Op{"EOR", (*CPU).Eor, Absolute},
            0x5d: Op{"EOR", (*CPU).Eor, AbsoluteX},
            0x59: Op{"EOR", (*CPU).Eor, AbsoluteY},
            0x41: Op{"EOR", (*CPU).Eor, IndexedIndirect},
            0x51: Op{"EOR", (*CPU).Eor, IndirectIndexed},
            0x4c: Op{"JMP", (*CPU).Jmp, Absolute},
            0x6c: Op{"JMP", (*CPU).Jmp, Indirect},
            0x20: Op{"JSR", (*CPU).Jsr, Absolute},
            0xa9: Op{"LDA", (*CPU).Lda, Immediate},
            0xa5: Op{"LDA", (*CPU).Lda, ZeroPage},
            0xb5: Op{"LDA", (*CPU).Lda, ZeroPageX},
            0xad: Op{"LDA", (*CPU).Lda, Absolute},
            0xbd: Op{"LDA", (*CPU).Lda, AbsoluteX},
            0xb9: Op{"LDA", (*CPU).Lda, AbsoluteY},
            0xa1: Op{"LDA", (*CPU).Lda, IndexedIndirect},
            0xb1: Op{"LDA", (*CPU).Lda, IndirectIndexed},
            0xa2: Op{"LDX", (*CPU).Ldx, Immediate},
            0xa6: Op{"LDX", (*CPU).Ldx, ZeroPage},
            0xb6: Op{"LDX", (*CPU).Ldx, ZeroPageY},
            0xae: Op{"LDX", (*CPU).Ldx, Absolute},
            0xbe: Op{"LDX", (*CPU).Ldx, AbsoluteY},
            0xa0: Op{"LDY", (*CPU).Ldy, Immediate},
            0xa4: Op{"LDY", (*CPU).Ldy, ZeroPage},
            0xb4: Op{"LDY", (*CPU).Ldy, ZeroPageX},
            0xac: Op{"LDY", (*CPU).Ldy, Absolute},
            0xbc: Op{"LDY", (*CPU).Ldy, AbsoluteX},
            0x4a: Op{"LSR", (*CPU).Lsr, nil},
            0x46: Op{"LSR", (*CPU).Lsr, ZeroPage},
            0x56: Op{"LSR", (*CPU).Lsr, ZeroPageX},
            0x4e: Op{"LSR", (*CPU).Lsr, Absolute},
            0x5e: Op{"LSR", (*CPU).Lsr, AbsoluteX},
            0xea: Op{"NOP", (*CPU).Nop, nil},
            0x09: Op{"ORA", (*CPU).Ora, Immediate},
            0x05: Op{"ORA", (*CPU).Ora, ZeroPage},
            0x15: Op{"ORA", (*CPU).Ora, ZeroPageX},
            0x0d: Op{"ORA", (*CPU).Ora, Absolute},
            0x1d: Op{"ORA", (*CPU).Ora, AbsoluteX},
            0x19: Op{"ORA", (*CPU).Ora, AbsoluteY},
            0x01: Op{"ORA", (*CPU).Ora, IndexedIndirect},
            0x11: Op{"ORA", (*CPU).Ora, IndirectIndexed},
            0x48: Op{"PHA", (*CPU).Pha, nil},
            0x08: Op{"PHP", (*CPU).Php, nil},
            0x68: Op{"PLA", (*CPU).Pla, nil},
            0x28: Op{"PLP", (*CPU).Plp, nil},
            0x2a: Op{"ROL", (*CPU).Rol, nil},
            0x26: Op{"ROL", (*CPU).Rol, ZeroPage},
            0x36: Op{"ROL", (*CPU).Rol, ZeroPageX},
            0x2e: Op{"ROL", (*CPU).Rol, Absolute},
            0x3e: Op{"ROL", (*CPU).Rol, AbsoluteX},
            0x6a: Op{"ROR", (*CPU).Ror, nil},
            0x66: Op{"ROR", (*CPU).Ror, ZeroPage},
            0x76: Op{"ROR", (*CPU).Ror, ZeroPageX},
            0x6e: Op{"ROR", (*CPU).Ror, Absolute},
            0x7e: Op{"ROR", (*CPU).Ror, AbsoluteX},
            0x40: Op{"RTI", (*CPU).Rti, nil},
            0x60: Op{"RTS", (*CPU).Rts, nil},
            0xe9: Op{"SBC", (*CPU).Sbc, Immediate},
            0xe5: Op{"SBC", (*CPU).Sbc, ZeroPage},
            0xf5: Op{"SBC", (*CPU).Sbc, ZeroPageX},
            0xed: Op{"SBC", (*CPU).Sbc, Absolute},
            0xfd: Op{"SBC", (*CPU).Sbc, AbsoluteX},
            0xf9: Op{"SBC", (*CPU).Sbc, AbsoluteY},
            0xe1: Op{"SBC", (*CPU).Sbc, IndexedIndirect},
            0xf1: Op{"SBC", (*CPU).Sbc, IndirectIndexed},
            0x38: Op{"SEC", (*CPU).Sec, nil},
            0xf8: Op{"SED", (*CPU).Sed, nil},
            0x78: Op{"SEI", (*CPU).Sei, nil},
            0x85: Op{"STA", (*CPU).Sta, ZeroPage},
            0x95: Op{"STA", (*CPU).Sta, ZeroPageX},
            0x8d: Op{"STA", (*CPU).Sta, Absolute},
            0x9d: Op{"STA", (*CPU).Sta, AbsoluteX},
            0x99: Op{"STA", (*CPU).Sta, AbsoluteY},
            0x81: Op{"STA", (*CPU).Sta, IndexedIndirect},
            0x91: Op{"STA", (*CPU).Sta, IndirectIndexed},
            0x84: Op{"STY", (*CPU).Sty, ZeroPage},
            0x94: Op{"STY", (*CPU).Sty, ZeroPageX},
            0x8c: Op{"STY", (*CPU).Sty, Absolute},
            0x86: Op{"STX", (*CPU).Stx, ZeroPage},
            0x96: Op{"STX", (*CPU).Stx, ZeroPageY},
            0x8e: Op{"STX", (*CPU).Stx, Absolute},
            0xaa: Op{"TAX", (*CPU).Tax, nil},
            0xa8: Op{"TAY", (*CPU).Tay, nil},
            0xba: Op{"TSX", (*CPU).Tsx, nil},
            0x8a: Op{"TXA", (*CPU).Txa, nil},
            0x9a: Op{"TXS", (*CPU).Txs, nil},
            0x98: Op{"TYA", (*CPU).Tya, nil},
        }
    }

    return p.operations;
}

func (p *CPU) Reset() {
    p.Flags = 0x24
    p.A, p.X, p.Y = 0x00, 0x00, 0x00
    p.SP = 0xfd
}

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

func (p *CPU) ZeroPageY() (Address, int) {
    return Address(p.Memory[p.PC] + p.Y), 1
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

func (p *CPU) Indirect() (Address, int) {
    location := p.absolute()

    high := p.Memory[location+1]
    low := p.Memory[location]

    return (Address(high) << 8) + Address(low), 2
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
    } else {
        p.setNegativeFlag(false)
    }

    if value == 0x00 {
        p.setZeroFlag(true)
    } else {
        p.setZeroFlag(false)
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
    } else {
        p.setCarryFlag(false)
    }

    if addOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    } else {
        p.setOverflowFlag(false)
    }

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Sbc(location Address) {
    other := p.Memory[location]
    old := p.A

    p.A -= other

    if !p.Carry() {
        p.A -= 0x01
    }

    if p.A > old {
        p.setCarryFlag(false)
    } else {
        p.setCarryFlag(true)
    }

    if subtractOverflowed(old, other, p.A) {
        p.setOverflowFlag(true)
    } else {
        p.setOverflowFlag(false)
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
    } else {
        p.setCarryFlag(false)
    }

    *memory = *memory << 1

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Bit(location Address) {
    result := p.A & p.Memory[location]

    if result == 0x00 {
        p.setZeroFlag(true)
    } else {
        p.setZeroFlag(false)
    }

    if p.Memory[location] & 0x40 == 0x40 {
        p.setOverflowFlag(true)
    } else {
        p.setOverflowFlag(false)
    }

    if p.Memory[location] & 0x80 == 0x80 {
        p.setNegativeFlag(true)
    } else {
        p.setNegativeFlag(false)
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
    } else {
        p.setCarryFlag(false)
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
    } else {
        p.setCarryFlag(false)
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
    p.Flags = (p.Flags | 0x30) - 0x10
}

func (p *CPU) Rol(memory *byte) {
    carried := (*memory & 0x80) == 0x80

    *memory = *memory << 1

    if p.Carry() {
        *memory |= 0x01
    }

    if carried {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
    }

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Ror(memory *byte) {
    carried := (*memory & 0x01) == 0x01

    *memory = *memory >> 1

    if p.Carry() {
        *memory |= 0x80
    }

    if carried {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
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
    p.Flags = (p.Flags | 0x30) - 0x10

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
