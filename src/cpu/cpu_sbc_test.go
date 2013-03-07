package cpu

import "testing"

func (p *CPU) sbc(first byte, second byte) {
    p.A = first
    p.Memory[0] = second
    p.Sbc(0)
}

func TestSbcNoOverflow(t *testing.T) {
    var p (*CPU) = new(CPU)

    p.sbc(0x01, 0x01)

    if p.A != 0x00 {
        t.Errorf("Sbc didn't subtract correctly")
        t.Errorf("Expected %#02x, got %#02x", 0x00, p.A)
    }
}

func TestSbcPositiveFromNegativeOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.sbc(0x80, 0x0f)

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when subtraction overflows")
    }
}

func TestSbcNegativeFromPositiveOverflowSetsOverflowFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.sbc(0x00, 0x80)

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when subtraction overflows")
    }
}

func TestSbcUnsignedOverflowSetsCarryFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.sbc(0x00, 0x01)

    if !p.Carry() {
        t.Errorf("Carry flag not set on unsigned subtraction overflow")
    }
}

func TestSbcWithCarryAlreadySetSubtracts1(t *testing.T) {
    var p *CPU = new(CPU)
    p.Flags |= 0x01
    p.sbc(0x01, 0x00)

    if p.A != 0x00 {
        t.Errorf("Carry flag didn't affect subtraction properly")
        t.Errorf("Expected %#02x, got %#02x", 0x00, p.A)
    }
}

func TestSbcSetsNegative(t *testing.T) {
    var p *CPU = new(CPU)

    p.sbc(0x00, 0x01)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when result is negative")
    }
}

func TestSbcSetsZero(t *testing.T) {
    var p *CPU = new(CPU)

    p.sbc(0x01, 0x01)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is negative")
    }
}
