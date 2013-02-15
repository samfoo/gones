package cpu

import "testing"

func (p *CPU) aslAccumulator(a byte) {
    p.A = a
    p.Asl(&p.A)
}

func (p *CPU) aslMemory(a byte) {
    p.Memory[0] = a
    p.Asl(&p.Memory[0])
}

func TestAslShiftsAccumulator(t *testing.T) {
    var p *CPU = new(CPU)

    p.aslAccumulator(0x01)

    if p.A != 0x02 {
        t.Errorf("Shift left didn't work on the accumulator")
        t.Errorf("Expected 0x02, got %#02x", p.A)
        t.FailNow()
    }
}

func TestAslShiftsMemoryLocation(t *testing.T) {
    var p *CPU = new(CPU)

    p.aslMemory(0x01)

    if p.Memory[0x00] != 0x02 {
        t.Errorf("Shift left didn't work on a memory location")
        t.Errorf("Expected 0x02, got %#02x", p.A)
        t.FailNow()
    }
}

func TestAslZeroFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.aslAccumulator(0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestAslNegativeFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.aslAccumulator(0x40)

    if !p.Negative() {
        t.Errorf("Negative flag should be set when high bit is 1 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestAslCarryFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.aslAccumulator(0x80)

    if !p.Carry() {
        t.Errorf("Carry flag should be set when shift won't fit (flags: %08b)", p.Flags)
        t.FailNow()
    }
}
