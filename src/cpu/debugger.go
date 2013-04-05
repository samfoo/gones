package cpu

import "fmt"

func (p *CPU) Debugf(opcode Opcode, op Op) {
    fmt.Printf("%4.04X  ", p.PC-1)

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
            location := p.zeroPage(&p.Memory)
            fmt.Printf("$%02X = %-22.02X", p.Memory.Read(p.PC), p.Memory.Read(location))
        case ZeroPageX:
            location := p.zeroPageX(&p.Memory)
            fmt.Printf("$%02X,X @ %02X = %-15.02X", p.Memory.Read(p.PC), location & 0xff, p.Memory.Read(location))
        case ZeroPageY:
            location := p.zeroPageY(&p.Memory)
            fmt.Printf("$%02X,Y @ %02X = %-15.02X", p.Memory.Read(p.PC), location & 0xff, p.Memory.Read(location))
        case Absolute:
            location := p.absolute(&p.Memory)

            if op.Name == "JMP" || op.Name == "JSR" {
                fmt.Printf("$%-27.04X", location)
            } else {
                fmt.Printf("$%04X = %-20.02X", location, p.Memory.Read(location))
            }
        case Indirect:
            location := p.absolute(&p.Memory)
            high := p.Memory.Read(location+1)
            low := p.Memory.Read(location)
            indirectLocation := (Address(high) << 8) + Address(low)

            fmt.Printf("($%04X) = %-18.04X", location, indirectLocation)
        case AbsoluteX:
            fmt.Printf("$%04X,X @ %04X = %-11.02X",
                p.absolute(&p.Memory),
                p.absolute(&p.Memory) + Address(p.X),
                p.Memory.Read(p.absolute(&p.Memory)+Address(p.X)),
            )
        case AbsoluteY:
            fmt.Printf("$%04X,Y @ %04X = %-11.02X",
                p.absolute(&p.Memory),
                p.absolute(&p.Memory) + Address(p.Y),
                p.Memory.Read(p.absolute(&p.Memory)+Address(p.Y)),
            )
        case IndexedIndirect:
            location := p.indexedIndirect(&p.Memory)
            fmt.Printf("($%02X,X) @ %02X = %04X = %-6.02X", p.Memory.Read(p.PC), p.Memory.Read(p.PC) + p.X, location, p.Memory.Read(location))
        case IndirectIndexed:
            location := p.indirectIndexed(&p.Memory)
            fmt.Printf("($%02X),Y = %04X @ %04X = %-4.02X", p.Memory.Read(p.PC), location - Address(p.Y), location, p.Memory.Read(location))
        case Relative:
            location := p.relative(&p.Memory)
            fmt.Printf("$%-27.04X", location)
        case Accumulator:
            fmt.Printf("%-28s", "A")
        case Implied:
            fmt.Printf("%-28s", " ")
    }

    fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", p.A, p.X, p.Y, p.Flags, p.SP)
}
