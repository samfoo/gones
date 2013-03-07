package cpu

import "testing"

func TestJmpSetsPC(t *testing.T) {
    var p *CPU = new(CPU)

    p.Jmp(0xbeef)

    if p.PC != 0xbeef {
        t.Errorf("Jmp didn't set PC correctly")
    }
}

func TestJsrStack(t *testing.T) {
    var p *CPU = new(CPU)

    p.SP = 0xff
    p.PC = 0xbeef
    p.Jsr(0xdead)

    stack := p.Memory[0x01fe:0x0200]

    if stack[0] != 0xee {
        t.Errorf("PC low byte not pushed properly to the stack")
        t.Errorf("Expected %#02x, got %#02x", 0xee, stack[0])
    }

    if stack[1] != 0xbe {
        t.Errorf("PC high byte not pushed properly to the stack")
        t.Errorf("Expected %#02x, got %#02x", 0xbe, stack[1])
    }
}

func TestJsrPC(t *testing.T) {
    var p *CPU = new(CPU)

    p.Jsr(0xbeef)

    if p.PC != 0xbeef {
        t.Errorf("Jsr didn't set PC correctly")
    }
}
