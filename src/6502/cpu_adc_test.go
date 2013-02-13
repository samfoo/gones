package cpu

import "testing"

func (p *CPU) addImmediate(first byte, second byte) {
    p.A = first
    p.Memory[0] = second
    p.Adc(0)
}

func TestZeroFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestAdcSetsNegativeFlagWhenHighBitIs1(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0x00, 0x80)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when result is negative (flags: %08b)", p.Flags)
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

    if p.Carry() {
        t.Errorf("Carry flag set when it shouldn't be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestOverflowSetsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.addImmediate(0xff, 0x01)

    if !p.Carry() {
        t.Errorf("Carry flag not set when it should be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestSignedPositiveOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = new(CPU)

    // First and second numbers' sign bit is 0 and the result's sign bit is 1
    // is positive overflow
    p.addImmediate(0x7f, 0x7f)

    if p.A & 0x80 != 0x80 {
        t.Errorf("Overflow arithmetic didn't work properly")
        t.FailNow()
    }

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when it should be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestSignedNegativeOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = new(CPU)

    // First and second numbers' sign bit is 1 and the result's sign bit is 0
    // is negative overflow
    p.addImmediate(0x80, 0x80)

    if p.A & 0x80 != 0x00 {
        t.Errorf("Overflow arithmetic didn't work properly %#02x", p.A)
        t.FailNow()
    }

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when it should be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
