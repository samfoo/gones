package cpu

import "fmt"

type Address uint16

type CPU struct {
    A, X, Y, SP, Flags byte
    PC Address
    Memory RAM
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

func (p *CPU) Execute(op Op) {
    switch m := op.Method.(type) {
        case func(*CPU, Address):
            m(p, p.Address(addressing[op.Mode.(int)]))
        case func(*CPU):
            m(p)
    }
}

func (p *CPU) Step() {
    opcode := Opcode(p.Memory.Read(p.PC))

    fmt.Printf("%4.04X  ", p.PC)

    p.PC++

    op := p.Operations()[opcode]

    fmt.Printf("%-2.02X ", opcode)

    switch op.Mode {
        case Immediate, ZeroPage, ZeroPageX, ZeroPageY, IndexedIndirect, IndirectIndexed, Relative:
            fmt.Printf("%-5.02X ", p.Memory.Read(p.PC))
        case Absolute, AbsoluteX, AbsoluteY, Indirect:
            fmt.Printf("%02X %-2.02X ", p.Memory.Read(p.PC), p.Memory.Read(p.PC+1))
        default:
            fmt.Printf("%-6s", " ")
    }

    fmt.Printf("%4s ", op.Name)

    switch op.Mode {
        case Immediate:
            fmt.Printf("#$%-26.02X", p.Memory.Read(p.PC))
        case ZeroPage:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X = %-22.02X", p.Memory.Read(p.PC), p.Memory.Read(location))
        case ZeroPageX:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X,X @ %02X = %-15.02X", p.Memory.Read(p.PC), location & 0xff, p.Memory.Read(location))
        case ZeroPageY:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%02X,Y @ %02X = %-15.02X", p.Memory.Read(p.PC), location & 0xff, p.Memory.Read(location))
        case Absolute:
            if op.Name == "JMP" || op.Name == "JSR" {
                fmt.Printf("$%-27.04X", p.absolute())
            } else {
                fmt.Printf("$%04X = %-20.02X", p.absolute(), p.Memory.Read(p.absolute()))
            }
        case Indirect:
            location := p.absolute()
            high := p.Memory.Read(location+1)
            low := p.Memory.Read(location)
            indirectLocation := (Address(high) << 8) + Address(low)

            fmt.Printf("($%04X) = %-18.04X", p.absolute(), indirectLocation)
        case AbsoluteX:
            fmt.Printf("$%04X,X @ %04X = %-11.02X", p.absolute(), p.absolute() + Address(p.X), p.Memory.Read(p.absolute()+Address(p.X)))
        case AbsoluteY:
            fmt.Printf("$%04X,Y @ %04X = %-11.02X", p.absolute(), p.absolute() + Address(p.Y), p.Memory.Read(p.absolute()+Address(p.Y)))
        case IndexedIndirect:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("($%02X,X) @ %02X = %04X = %-6.02X", p.Memory.Read(p.PC), p.Memory.Read(p.PC) + p.X, location, p.Memory.Read(location))
        case IndirectIndexed:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("($%02X),Y = %04X @ %04X = %-4.02X", p.Memory.Read(p.PC), location - Address(p.Y), location, p.Memory.Read(location))
        case Relative:
            location, _ := addressing[op.Mode.(int)](p)
            fmt.Printf("$%-27.04X", location)
        case Accumulator:
            fmt.Printf("%-28s", "A")
        case Implied:
            fmt.Printf("%-28s", " ")
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
            0x87: Op{"*SAX", (*CPU).Sax, ZeroPage},
            0x97: Op{"*SAX", (*CPU).Sax, ZeroPageY},
            0x83: Op{"*SAX", (*CPU).Sax, IndexedIndirect},
            0x8f: Op{"*SAX", (*CPU).Sax, Absolute},
            0x0a: Op{"ASL", (*CPU).AslAcc, Accumulator},
            0x06: Op{"ASL", (*CPU).Asl, ZeroPage},
            0x16: Op{"ASL", (*CPU).Asl, ZeroPageX},
            0x0e: Op{"ASL", (*CPU).Asl, Absolute},
            0x1e: Op{"ASL", (*CPU).Asl, AbsoluteX},
            0x07: Op{"*SLO", (*CPU).Slo, ZeroPage},
            0x17: Op{"*SLO", (*CPU).Slo, ZeroPageX},
            0x0f: Op{"*SLO", (*CPU).Slo, Absolute},
            0x1f: Op{"*SLO", (*CPU).Slo, AbsoluteX},
            0x1b: Op{"*SLO", (*CPU).Slo, AbsoluteY},
            0x03: Op{"*SLO", (*CPU).Slo, IndexedIndirect},
            0x13: Op{"*SLO", (*CPU).Slo, IndirectIndexed},
            0x90: Op{"BCC", (*CPU).Bcc, Relative},
            0xb0: Op{"BCS", (*CPU).Bcs, Relative},
            0xf0: Op{"BEQ", (*CPU).Beq, Relative},
            0xd0: Op{"BNE", (*CPU).Bne, Relative},
            0x24: Op{"BIT", (*CPU).Bit, ZeroPage},
            0x2c: Op{"BIT", (*CPU).Bit, Absolute},
            0x30: Op{"BMI", (*CPU).Bmi, Relative},
            0x10: Op{"BPL", (*CPU).Bpl, Relative},
            0x00: Op{"BRK", (*CPU).Brk, Implied},
            0x50: Op{"BVC", (*CPU).Bvc, Relative},
            0x70: Op{"BVS", (*CPU).Bvs, Relative},
            0x18: Op{"CLC", (*CPU).Clc, Implied},
            0xd8: Op{"CLD", (*CPU).Cld, Implied},
            0x58: Op{"CLI", (*CPU).Cli, Implied},
            0xb8: Op{"CLV", (*CPU).Clv, Implied},
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
            0xc7: Op{"*DCP", (*CPU).Dcp, ZeroPage},
            0xd7: Op{"*DCP", (*CPU).Dcp, ZeroPageX},
            0xcf: Op{"*DCP", (*CPU).Dcp, Absolute},
            0xdf: Op{"*DCP", (*CPU).Dcp, AbsoluteX},
            0xdb: Op{"*DCP", (*CPU).Dcp, AbsoluteY},
            0xc3: Op{"*DCP", (*CPU).Dcp, IndexedIndirect},
            0xd3: Op{"*DCP", (*CPU).Dcp, IndirectIndexed},
            0xca: Op{"DEX", (*CPU).Dex, Implied},
            0x88: Op{"DEY", (*CPU).Dey, Implied},
            0xe6: Op{"INC", (*CPU).Inc, ZeroPage},
            0xf6: Op{"INC", (*CPU).Inc, ZeroPageX},
            0xee: Op{"INC", (*CPU).Inc, Absolute},
            0xfe: Op{"INC", (*CPU).Inc, AbsoluteX},
            0xe7: Op{"*ISB", (*CPU).Isb, ZeroPage},
            0xf7: Op{"*ISB", (*CPU).Isb, ZeroPageX},
            0xef: Op{"*ISB", (*CPU).Isb, Absolute},
            0xff: Op{"*ISB", (*CPU).Isb, AbsoluteX},
            0xfb: Op{"*ISB", (*CPU).Isb, AbsoluteY},
            0xe3: Op{"*ISB", (*CPU).Isb, IndexedIndirect},
            0xf3: Op{"*ISB", (*CPU).Isb, IndirectIndexed},
            0xe8: Op{"INX", (*CPU).Inx, Implied},
            0xc8: Op{"INY", (*CPU).Iny, Implied},
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
            0xa7: Op{"*LAX", (*CPU).Lax, ZeroPage},
            0xb7: Op{"*LAX", (*CPU).Lax, ZeroPageY},
            0xaf: Op{"*LAX", (*CPU).Lax, Absolute},
            0xbf: Op{"*LAX", (*CPU).Lax, AbsoluteY},
            0xa3: Op{"*LAX", (*CPU).Lax, IndexedIndirect},
            0xb3: Op{"*LAX", (*CPU).Lax, IndirectIndexed},
            0x4a: Op{"LSR", (*CPU).LsrAcc, Accumulator},
            0x46: Op{"LSR", (*CPU).Lsr, ZeroPage},
            0x56: Op{"LSR", (*CPU).Lsr, ZeroPageX},
            0x4e: Op{"LSR", (*CPU).Lsr, Absolute},
            0x5e: Op{"LSR", (*CPU).Lsr, AbsoluteX},
            0x47: Op{"*SRE", (*CPU).Sre, ZeroPage},
            0x57: Op{"*SRE", (*CPU).Sre, ZeroPageX},
            0x4f: Op{"*SRE", (*CPU).Sre, Absolute},
            0x5f: Op{"*SRE", (*CPU).Sre, AbsoluteX},
            0x5b: Op{"*SRE", (*CPU).Sre, AbsoluteY},
            0x43: Op{"*SRE", (*CPU).Sre, IndexedIndirect},
            0x53: Op{"*SRE", (*CPU).Sre, IndirectIndexed},
            0xea: Op{"NOP", (*CPU).Nop, Implied},
            0x04: Op{"*NOP", (*CPU)._Nop, ZeroPage},
            0x14: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0x34: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0x44: Op{"*NOP", (*CPU)._Nop, ZeroPage},
            0x54: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0x64: Op{"*NOP", (*CPU)._Nop, ZeroPage},
            0x74: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0x80: Op{"*NOP", (*CPU)._Nop, Immediate},
            0x82: Op{"*NOP", (*CPU)._Nop, Immediate},
            0x89: Op{"*NOP", (*CPU)._Nop, Immediate},
            0xc2: Op{"*NOP", (*CPU)._Nop, Immediate},
            0xd4: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0xe2: Op{"*NOP", (*CPU)._Nop, Immediate},
            0xf4: Op{"*NOP", (*CPU)._Nop, ZeroPageX},
            0x0c: Op{"*NOP", (*CPU)._Nop, Absolute},
            0x1c: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0x3c: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0x5c: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0x7c: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0xdc: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0xfc: Op{"*NOP", (*CPU)._Nop, AbsoluteX},
            0x1a: Op{"*NOP", (*CPU).Nop, Implied},
            0x3a: Op{"*NOP", (*CPU).Nop, Implied},
            0x5a: Op{"*NOP", (*CPU).Nop, Implied},
            0x7a: Op{"*NOP", (*CPU).Nop, Implied},
            0xda: Op{"*NOP", (*CPU).Nop, Implied},
            0xfa: Op{"*NOP", (*CPU).Nop, Implied},
            0x09: Op{"ORA", (*CPU).Ora, Immediate},
            0x05: Op{"ORA", (*CPU).Ora, ZeroPage},
            0x15: Op{"ORA", (*CPU).Ora, ZeroPageX},
            0x0d: Op{"ORA", (*CPU).Ora, Absolute},
            0x1d: Op{"ORA", (*CPU).Ora, AbsoluteX},
            0x19: Op{"ORA", (*CPU).Ora, AbsoluteY},
            0x01: Op{"ORA", (*CPU).Ora, IndexedIndirect},
            0x11: Op{"ORA", (*CPU).Ora, IndirectIndexed},
            0x48: Op{"PHA", (*CPU).Pha, Implied},
            0x08: Op{"PHP", (*CPU).Php, Implied},
            0x68: Op{"PLA", (*CPU).Pla, Implied},
            0x28: Op{"PLP", (*CPU).Plp, Implied},
            0x2a: Op{"ROL", (*CPU).RolAcc, Accumulator},
            0x26: Op{"ROL", (*CPU).Rol, ZeroPage},
            0x36: Op{"ROL", (*CPU).Rol, ZeroPageX},
            0x2e: Op{"ROL", (*CPU).Rol, Absolute},
            0x3e: Op{"ROL", (*CPU).Rol, AbsoluteX},
            0x27: Op{"*RLA", (*CPU).Rla, ZeroPage},
            0x37: Op{"*RLA", (*CPU).Rla, ZeroPageX},
            0x2f: Op{"*RLA", (*CPU).Rla, Absolute},
            0x3f: Op{"*RLA", (*CPU).Rla, AbsoluteX},
            0x3b: Op{"*RLA", (*CPU).Rla, AbsoluteY},
            0x23: Op{"*RLA", (*CPU).Rla, IndexedIndirect},
            0x33: Op{"*RLA", (*CPU).Rla, IndirectIndexed},
            0x6a: Op{"ROR", (*CPU).RorAcc, Accumulator},
            0x66: Op{"ROR", (*CPU).Ror, ZeroPage},
            0x76: Op{"ROR", (*CPU).Ror, ZeroPageX},
            0x6e: Op{"ROR", (*CPU).Ror, Absolute},
            0x7e: Op{"ROR", (*CPU).Ror, AbsoluteX},
            0x67: Op{"*RRA", (*CPU).Rra, ZeroPage},
            0x77: Op{"*RRA", (*CPU).Rra, ZeroPageX},
            0x6f: Op{"*RRA", (*CPU).Rra, Absolute},
            0x7f: Op{"*RRA", (*CPU).Rra, AbsoluteX},
            0x7b: Op{"*RRA", (*CPU).Rra, AbsoluteY},
            0x63: Op{"*RRA", (*CPU).Rra, IndexedIndirect},
            0x73: Op{"*RRA", (*CPU).Rra, IndirectIndexed},
            0x40: Op{"RTI", (*CPU).Rti, Implied},
            0x60: Op{"RTS", (*CPU).Rts, Implied},
            0xe9: Op{"SBC", (*CPU).Sbc, Immediate},
            0xe5: Op{"SBC", (*CPU).Sbc, ZeroPage},
            0xf5: Op{"SBC", (*CPU).Sbc, ZeroPageX},
            0xed: Op{"SBC", (*CPU).Sbc, Absolute},
            0xfd: Op{"SBC", (*CPU).Sbc, AbsoluteX},
            0xf9: Op{"SBC", (*CPU).Sbc, AbsoluteY},
            0xe1: Op{"SBC", (*CPU).Sbc, IndexedIndirect},
            0xf1: Op{"SBC", (*CPU).Sbc, IndirectIndexed},
            0xeb: Op{"*SBC", (*CPU).Sbc, Immediate},
            0x38: Op{"SEC", (*CPU).Sec, Implied},
            0xf8: Op{"SED", (*CPU).Sed, Implied},
            0x78: Op{"SEI", (*CPU).Sei, Implied},
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
            0xaa: Op{"TAX", (*CPU).Tax, Implied},
            0xa8: Op{"TAY", (*CPU).Tay, Implied},
            0xba: Op{"TSX", (*CPU).Tsx, Implied},
            0x8a: Op{"TXA", (*CPU).Txa, Implied},
            0x9a: Op{"TXS", (*CPU).Txs, Implied},
            0x98: Op{"TYA", (*CPU).Tya, Implied},
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
    var offset = Address(p.Memory.Read(p.PC))
    if offset < 0x0080 {
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
    return Address(p.Memory.Read(p.PC)), 1
}

func (p *CPU) ZeroPageX() (Address, int) {
    return Address(p.Memory.Read(p.PC) + p.X), 1
}

func (p *CPU) ZeroPageY() (Address, int) {
    return Address(p.Memory.Read(p.PC) + p.Y), 1
}

func (p *CPU) absolute() Address {
    high := p.Memory.Read(p.PC+1)
    low := p.Memory.Read(p.PC)

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

    low := p.Memory.Read(location)

    var high byte
    if location & 0x00ff == 0x00ff {
        high = p.Memory.Read(location & 0xff00)
    } else {
        high = p.Memory.Read(location+1)
    }

    return (Address(high) << 8) + Address(low), 2
}

func (p *CPU) IndexedIndirect() (Address, int) {
    pointer := p.Memory.Read(p.PC) + p.X

    high := p.Memory.Read(Address(pointer+1))
    low := p.Memory.Read(Address(pointer))

    return (Address(high) << 8) + Address(low), 1
}

func (p *CPU) indirectIndexed() (Address, Address) {
    indirect := Address(p.Memory.Read(p.PC))

    high := p.Memory.Read(indirect+1)
    low := p.Memory.Read(indirect)

    return indirect, (Address(high) << 8) + Address(low) + Address(p.Y)
}

func (p *CPU) IndirectIndexed() (Address, int) {
    indirect := p.Memory.Read(p.PC)

    high := p.Memory.Read(Address(indirect+1))
    low := p.Memory.Read(Address(indirect))

    return ((Address(high) << 8) + Address(low)) + Address(p.Y), 1
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
    other := p.Memory.Read(location)
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
    other := p.Memory.Read(location)
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

func (p *CPU) Sax(location Address) {
    p.Memory.Write(p.A & p.X, location)
}

func (p *CPU) And(location Address) {
    other := p.Memory.Read(location)
    p.A &= other

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Ora(location Address) {
    other := p.Memory.Read(location)
    p.A |= other

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Slo(location Address) {
    p.Asl(location)
    p.Ora(location)
}

func (p *CPU) asl(val byte) byte {
    if val & 0x80 == 0x80 {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
    }

    val = val << 1

    p.setNegativeAndZeroFlags(val)

    return val
}

func (p *CPU) AslAcc() {
    p.A = p.asl(p.A)
}

func (p *CPU) Asl(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(p.asl(val), location)
}

func (p *CPU) Aac(location Address) {
    p.A &= p.Memory.Read(location)

    if p.A == 0x00 {
        p.setZeroFlag(true)
    } else {
        p.setZeroFlag(false)
    }

    if p.A & 0x80 == 0x80 {
        p.setNegativeFlag(true)
        p.setCarryFlag(true)
    } else {
        p.setNegativeFlag(false)
        p.setCarryFlag(false)
    }
}

func (p *CPU) Bit(location Address) {
    result := p.A & p.Memory.Read(location)

    if result == 0x00 {
        p.setZeroFlag(true)
    } else {
        p.setZeroFlag(false)
    }

    if p.Memory.Read(location) & 0x40 == 0x40 {
        p.setOverflowFlag(true)
    } else {
        p.setOverflowFlag(false)
    }

    if p.Memory.Read(location) & 0x80 == 0x80 {
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
    p.compare(p.A, p.Memory.Read(location))
}

func (p *CPU) Cpx(location Address) {
    p.compare(p.X, p.Memory.Read(location))
}

func (p *CPU) Cpy(location Address) {
    p.compare(p.Y, p.Memory.Read(location))
}

func (p *CPU) Dcp(location Address) {
    p.Memory.Write(p.Memory.Read(location)-1, location)
    p.compare(p.A, p.Memory.Read(location))
}

func (p *CPU) Dec(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(val - 1, location)
    p.setNegativeAndZeroFlags(val - 1)
}

func (p *CPU) Dex() {
    p.X -= 1
    p.setNegativeAndZeroFlags(p.X)
}

func (p *CPU) Dey() {
    p.Y -= 1
    p.setNegativeAndZeroFlags(p.Y)
}

func (p *CPU) Eor(location Address) {
    p.A ^= p.Memory.Read(location)

    p.setNegativeAndZeroFlags(p.A)
}

func (p *CPU) Isb(location Address) {
    p.Inc(location)
    p.Sbc(location)
}

func (p *CPU) Inc(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(val + 1, location)
    p.setNegativeAndZeroFlags(val + 1)
}

func (p *CPU) Inx() {
    p.X += 1
    p.setNegativeAndZeroFlags(p.X)
}

func (p *CPU) Iny() {
    p.Y += 1
    p.setNegativeAndZeroFlags(p.Y)
}

func (p *CPU) load(memory *byte, value byte) {
    *memory = value

    p.setNegativeAndZeroFlags(*memory)
}

func (p *CPU) Lax(location Address) {
    p.load(&p.A, p.Memory.Read(location))
    p.load(&p.X, p.Memory.Read(location))
}

func (p *CPU) Lda(location Address) {
    p.load(&p.A, p.Memory.Read(location))
}

func (p *CPU) Ldx(location Address) {
    p.load(&p.X, p.Memory.Read(location))
}

func (p *CPU) Ldy(location Address) {
    p.load(&p.Y, p.Memory.Read(location))
}

func (p *CPU) Sre(location Address) {
    p.Lsr(location)
    p.Eor(location)
}

func (p *CPU) lsr(val byte) byte {
    if val & 0x01 == 0x01 {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
    }

    val = val >> 1
    p.setNegativeAndZeroFlags(val)

    return val
}

func (p *CPU) LsrAcc() {
    p.A = p.lsr(p.A)
}

func (p *CPU) Lsr(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(p.lsr(val), location)
}

func (p *CPU) _Nop(location Address) {}

func (p *CPU) Nop() {}

func (p *CPU) push(value byte) {
    p.Memory.Write(value, 0x0100 + Address(p.SP))

    p.SP -= 1
}

func (p *CPU) pull(memory *byte) {
    p.SP += 1

    *memory = p.Memory.Read(0x0100 + Address(p.SP))
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

func (p *CPU) Rla(location Address) {
    p.Rol(location)
    p.And(location)
}

func (p *CPU) rol(val byte) byte {
    carried := (val & 0x80) == 0x80

    val = val << 1

    if p.Carry() {
        val |= 0x01
    }

    if carried {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
    }

    p.setNegativeAndZeroFlags(val)

    return val
}

func (p *CPU) RolAcc() {
    p.A = p.rol(p.A)
}

func (p *CPU) Rol(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(p.rol(val), location)
}

func (p *CPU) Rra(location Address) {
    p.Ror(location)
    p.Adc(location)
}

func (p *CPU) ror(val byte) byte {
    carried := (val & 0x01) == 0x01

    val = val >> 1

    if p.Carry() {
        val |= 0x80
    }

    if carried {
        p.setCarryFlag(true)
    } else {
        p.setCarryFlag(false)
    }

    p.setNegativeAndZeroFlags(val)

    return val
}

func (p *CPU) RorAcc() {
    p.A = p.ror(p.A)
}

func (p *CPU) Ror(location Address) {
    var val = p.Memory.Read(location)
    p.Memory.Write(p.ror(val), location)
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
    p.Memory.Write(p.A, location)
}

func (p *CPU) Stx(location Address) {
    p.Memory.Write(p.X, location)
}

func (p *CPU) Sty(location Address) {
    p.Memory.Write(p.Y, location)
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

    p.PC = Address(p.Memory.Read(0xffff)) << 8
    p.PC |= Address(p.Memory.Read(0xfffe))
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
