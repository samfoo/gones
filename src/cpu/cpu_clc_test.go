package cpu

import "testing"

func TestClcClearsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.Flags = 0xff
    p.Clc()

    if p.Carry() {
        t.Errorf("Carry flag set when it should have been cleared (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
