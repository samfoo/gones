package cpu

import "testing"

func (p *CPU) adc(first byte, second byte) {
    p.A = first
    p.Memory.Write(second, 0)
    p.Adc(0)
}

func TestZeroFlagSet(t *testing.T) {
    var p *CPU = NewCPU()

    p.adc(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00")
    }

    p.adc(0x01, 0x00)
    if p.Zero() {
        t.Errorf("Zero flag should be unset when result isn't zero")
    }
}

func TestAdcSetsNegative(t *testing.T) {
    var p *CPU = NewCPU()

    p.adc(0x00, 0x80)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when result is negative (flags: %08b)", p.Flags)
    }
}

func TestAdcNoOverflow(t *testing.T) {
    var p *CPU = NewCPU()

    p.adc(0x00, 0x99)

    if p.A != 0x99 {
        t.Errorf("Accumulator should be 0x99, but was %#02x", p.A)
    }
}

func TestAdcOverflow(t *testing.T) {
    var p *CPU = NewCPU()

    p.adc(0xff, 0x01)

    if p.A != 0x00 {
        t.Errorf("Accumulator should overflow to 0x00, but was %#02x", p.A)
    }
}

func TestNoUnsignedOverflowUnsetsCarryFlag(t *testing.T) {
    var p *CPU = NewCPU()

    p.setCarryFlag(true)
    p.adc(0x00, 0x01)

    if p.Carry() {
        t.Errorf("Carry flag set when it shouldn't be (flags: %08b)", p.Flags)
    }
}

func TestAdcUnsignedOverflowSetsCarryFlag(t *testing.T) {
    var p *CPU = NewCPU()

    p.adc(0xff, 0x01)

    if !p.Carry() {
        t.Errorf("Carry flag not set when it should be (flags: %08b)", p.Flags)
    }
}

func TestAdcWithCarryAlreadySetAdds1(t *testing.T) {
    var p *CPU = NewCPU()
    p.Flags |= 0x01
    p.adc(0x00, 0x00)

    if p.A != 0x01 {
        t.Errorf("Carry flag didn't affect addition properly")
        t.Errorf("Expected %#02x, got %#02x", 0x01, p.A)
    }
}

func TestSignedPositiveOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = NewCPU()

    // First and second numbers' sign bit is 0 and the result's sign bit is 1
    // is positive overflow
    p.adc(0x7f, 0x7f)

    if p.A & 0x80 != 0x80 {
        t.Errorf("Overflow arithmetic didn't work properly")
    }

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when it should be (flags: %08b)", p.Flags)
    }
}

func TestSignedNegativeOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = NewCPU()

    // First and second numbers' sign bit is 1 and the result's sign bit is 0
    // is negative overflow
    p.adc(0x80, 0x80)

    if p.A & 0x80 != 0x00 {
        t.Errorf("Overflow arithmetic didn't work properly %#02x", p.A)
    }

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when it should be (flags: %08b)", p.Flags)
    }

    p.adc(0x00, 0x80)
    if p.Overflow() {
        t.Errorf("Overflow flag set when it should have been unset")
    }
}

