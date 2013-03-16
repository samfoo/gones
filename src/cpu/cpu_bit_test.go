package cpu

import "testing"

func (p *CPU) bitTest(first byte, second byte) {
    p.A = first
    p.Memory.buffer[0] = second
    p.Bit(0)
}

func TestBitSetsZeroFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.bitTest(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
    }

    p = new(CPU)
    p.setZeroFlag(true)

    p.bitTest(0x40, 0xff)
    if p.Zero() {
        t.Errorf("Zero flag should be unset when result is not 0x00")
    }
}

func TestBitSetsOverflowFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.bitTest(0x40, 0xff)

    if !p.Overflow() {
        t.Errorf("Overflow flag not set when it should be (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestBitSetsNegativeFlag(t *testing.T) {
    var p *CPU = new(CPU)

    p.bitTest(0x80, 0xff)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when result is negative (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

