package cpu

import "testing"

func (p *CPU) addImmediate(first byte, second byte) {
    p.A = first
    p.Memory[0x00] = second
    p.Adc(0x00)
}

func TestZeroFlagSets(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0x00, 0x00)

    if p.Flags & 0x02 != 0x02 {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestNoOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0x00, 0x99)

    if p.A != 0x99 {
        t.Errorf("Accumulator should be 0x99, but was %#02x", p.A)
        t.FailNow()
    }
}

func TestOverflow(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0xff, 0x01)

    if p.A != 0x00 {
        t.Errorf("Accumulator should overflow to 0x00, but was %#02x", p.A)
        t.FailNow()
    }
}

func TestNoOverflowUnsetsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0x00, 0x01)

    if p.Flags & 0x01 != 0x00 {
        t.Errorf("Carry flag set when it shouldn't be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestOverflowSetsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0xff, 0x01)

    if p.Flags & 0x01 != 0x01 {
        t.Errorf("Carry flag not set when it should be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
