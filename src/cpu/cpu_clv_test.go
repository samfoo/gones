package cpu

import "testing"

func TestClvClearsOverflowFlag(t *testing.T) {
    var p *CPU = NewCPU()

    p.Flags = 0xff
    p.Clv()

    if p.Overflow() {
        t.Errorf("Overflow flag set when it should have been cleared (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
