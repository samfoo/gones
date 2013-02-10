package cpu

import "testing"

func TestImmediateNoOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.Memory[0x00] = 0x99
    p.Adc(0x00)

    if p.A != 0x99 {
        t.Errorf("Accumulator should be 0x99, but was 0x%x")
        t.FailNow()
    }
}

func TestImmediateOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.A = 0xff
    p.Memory[0x00] = 0x01

    p.Adc(0x00)

    if p.A != 0x00 {
        t.Errorf("Accumulator should overflow to 0x00, but was 0x%x", p.A)
        t.FailNow()
    }
}

func TestImmediateOverflowSetsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.A = 0xff
    p.Memory[0x00] = 0x01

    p.Adc(0x00)

    if p.Flags & 0x01 != 0x01 {
        t.Errorf("Carry flag not set")
        t.FailNow()
    }
}
