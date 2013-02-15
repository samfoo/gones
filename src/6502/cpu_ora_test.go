package cpu

import "testing"

func (p *CPU) oraImmediate(first byte, second byte) {
    p.A = first
    p.Memory[0] = second
    p.Ora(0)
}

func TestOraSetsAccumulator(t *testing.T) {
    var p *CPU = new(CPU)

    p.oraImmediate(0x01, 0xff)

    if p.A != 0xff {
        t.Errorf("Binary ora seems not to have worked correctly")
        t.Errorf("Expected 0xff, got %#02x", p.A)
        t.FailNow()
    }
}

func TestOraZeroFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.oraImmediate(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestOraNegativeFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.oraImmediate(0x80, 0x80)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when high bit is 1 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
