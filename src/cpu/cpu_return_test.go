package cpu

import "testing"

func TestRtiStack(t *testing.T) {
    var p *CPU = NewCPU()
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)

    p.SP = 0xff
    p.Flags = 0x00
    p.PC = 0xbeef

    p.Brk()

    p.Flags = 0x00
    p.PC = 0x00

    p.Rti()

    if p.Flags != 0x20 {
        t.Errorf("Flags were not pulled properly from the stack")
        t.Errorf("Expected %#02x, got %#02x", 0x20, p.Flags)
    }

    if p.PC != 0xbef0 {
        t.Errorf("PC was not pulled properly from the stack")
        t.Errorf("Expected %#04x, got %#04x", 0xbef0, p.PC)
    }
}

func TestRtsStack(t *testing.T) {
    var p *CPU = NewCPU()
    p.SP = 0xff
    p.Flags = 0x00
    p.PC = 0xbeef
    p.Memory.Mount(NewRAM(0xe000), 0x2000, 0xffff)

    p.Jsr(0xdead)
    p.Rts()

    if p.PC != 0xbeef {
        t.Errorf("PC was not pulled properly from the stack")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, p.PC)
    }
}
