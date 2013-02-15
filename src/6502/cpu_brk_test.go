package cpu

import "testing"

func TestBrkStack(t *testing.T) {
    var p *CPU = new(CPU)
    p.SP = 0xff
    p.Flags = 0x00
    p.PC = 0xbeef

    p.Brk()

    stack := p.Memory[0x01fd:0x0200]

    if stack[0] != 0x10 {
        t.Errorf("Flags were not pushed properly to the stack")
        t.Errorf("Expected %#02x, got %#02x", 0x10, stack[0])
    }

    if stack[1] != 0xef + 1 {
        t.Errorf("PC low byte not pushed properly to the stack")
        t.Errorf("Expected %#02x, got %#02x", 0xef + 1, stack[1])
    }

    if stack[2] != 0xbe {
        t.Errorf("PC high byte not pushed properly to the stack")
        t.Errorf("Expected %#02x, got %#02x", 0xbe, stack[2])
    }
}

func TestBrkLoadsPC(t *testing.T) {
    var p *CPU = new(CPU)
    p.SP = 0xff
    p.Flags = 0x00
    p.PC = 0x0000

    // Set the interrupt program pointer
    p.Memory[0xffff] = 0xbe
    p.Memory[0xfffe] = 0xef

    p.Brk()

    if p.PC != 0xbeef {
        t.Errorf("PC not loaded with IRQ interrupt vector")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, p.PC)
    }
}

func TestBrkSetsInterruptDiable(t *testing.T) {
    var p *CPU = new(CPU)
    p.SP = 0xff
    p.Flags = 0x00
    p.PC = 0x0000

    p.Brk()

    if !p.InterruptDisable() {
        t.Errorf("Interrupts not disabled on breaking")
    }
}
