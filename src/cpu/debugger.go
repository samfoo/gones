package cpu

import "fmt"

func (p *CPU) Debugf(opcode Opcode, op Op) {
    var debug = ""
    debug += fmt.Sprintf("%4.04X  ", p.PC-1)

    debug += fmt.Sprintf("%-2.02X ", opcode)

    switch op.Mode {
        case Immediate, ZeroPage, ZeroPageX, ZeroPageY, IndexedIndirect, IndirectIndexed, Relative:
            debug += fmt.Sprintf("%-5.02X ", p.Memory.ReadDebug(p.PC))
        case Absolute, AbsoluteX, AbsoluteY, Indirect:
            debug += fmt.Sprintf("%02X %-2.02X ", p.Memory.ReadDebug(p.PC), p.Memory.ReadDebug(p.PC+1))
        default:
            debug += fmt.Sprintf("%-6s", " ")
    }

    debug += fmt.Sprintf("%4s ", op.Name)

    switch op.Mode {
        case Immediate:
            debug += fmt.Sprintf("#$%-26.02X", p.Memory.ReadDebug(p.PC))
        case ZeroPage:
            location := p.zeroPage(&p.Memory)
            debug += fmt.Sprintf("$%02X = %-22.02X", p.Memory.ReadDebug(p.PC), p.Memory.ReadDebug(location))
        case ZeroPageX:
            location := p.zeroPageX(&p.Memory)
            debug += fmt.Sprintf("$%02X,X @ %02X = %-15.02X", p.Memory.ReadDebug(p.PC), location & 0xff, p.Memory.ReadDebug(location))
        case ZeroPageY:
            location := p.zeroPageY(&p.Memory)
            debug += fmt.Sprintf("$%02X,Y @ %02X = %-15.02X", p.Memory.ReadDebug(p.PC), location & 0xff, p.Memory.ReadDebug(location))
        case Absolute:
            location := p.absolute(&p.Memory)

            if op.Name == "JMP" || op.Name == "JSR" {
                debug += fmt.Sprintf("$%-27.04X", location)
            } else {
                debug += fmt.Sprintf("$%04X = %-20.02X", location, p.Memory.ReadDebug(location))
            }
        case Indirect:
            location := p.absolute(&p.Memory)
            high := p.Memory.ReadDebug(location+1)
            low := p.Memory.ReadDebug(location)
            indirectLocation := (Address(high) << 8) + Address(low)

            debug += fmt.Sprintf("($%04X) = %-18.04X", location, indirectLocation)
        case AbsoluteX:
            debug += fmt.Sprintf("$%04X,X @ %04X = %-11.02X",
                p.absolute(&p.Memory),
                p.absolute(&p.Memory) + Address(p.X),
                p.Memory.ReadDebug(p.absolute(&p.Memory)+Address(p.X)),
            )
        case AbsoluteY:
            debug += fmt.Sprintf("$%04X,Y @ %04X = %-11.02X",
                p.absolute(&p.Memory),
                p.absolute(&p.Memory) + Address(p.Y),
                p.Memory.ReadDebug(p.absolute(&p.Memory)+Address(p.Y)),
            )
        case IndexedIndirect:
            location := p.indexedIndirect(&p.Memory)
            debug += fmt.Sprintf("($%02X,X) @ %02X = %04X = %-6.02X", p.Memory.ReadDebug(p.PC), p.Memory.ReadDebug(p.PC) + p.X, location, p.Memory.ReadDebug(location))
        case IndirectIndexed:
            location := p.indirectIndexed(&p.Memory)
            debug += fmt.Sprintf("($%02X),Y = %04X @ %04X = %-4.02X", p.Memory.ReadDebug(p.PC), location - Address(p.Y), location, p.Memory.ReadDebug(location))
        case Relative:
            location := p.relative(&p.Memory)
            debug += fmt.Sprintf("$%-27.04X", location)
        case Accumulator:
            debug += fmt.Sprintf("%-28s", "A")
        case Implied:
            debug += fmt.Sprintf("%-28s", " ")
    }

    debug += fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", p.A, p.X, p.Y, p.Flags, p.SP)

    fmt.Printf(debug)
}
