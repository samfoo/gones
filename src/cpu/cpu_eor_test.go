package cpu

import "testing"

func (p *CPU) eorImmediate(first byte, second byte) {
    p.A = first
    p.Memory.Write(second, 0)
    p.Eor(0)
}

func TestEorSetsAccumulator(t *testing.T) {
    var p *CPU = NewCPU()

    p.eorImmediate(0x01, 0xff)

    if p.A != 0xfe {
        t.Errorf("Binary exclusive or seems not to have worked correctly")
        t.Errorf("Expected 0xfe, got %#02x", p.A)
        t.FailNow()
    }
}

func TestEorZeroFlagSet(t *testing.T) {
    var p *CPU = NewCPU()

    p.eorImmediate(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestEorNegativeFlagSet(t *testing.T) {
    var p *CPU = NewCPU()

    p.eorImmediate(0x80, 0x00)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when high bit is 1 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
