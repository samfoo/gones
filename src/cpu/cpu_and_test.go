package cpu

import "testing"

func (p *CPU) andImmediate(first byte, second byte) {
    p.A = first
    p.Memory[0] = second
    p.And(0)
}

func TestAndSetsAccumulator(t *testing.T) {
    var p *CPU = new(CPU)

    p.andImmediate(0x01, 0xff)

    if p.A != 0x01 {
        t.Errorf("Binary and seems not to have worked correctly")
        t.Errorf("Expected 0x01, got %#02x", p.A)
        t.FailNow()
    }
}

func TestAndZeroFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.andImmediate(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestAndNegativeFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.andImmediate(0x80, 0x80)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when high bit is 1 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}