package cpu

import "testing"

func (p *CPU) lsrAccumulator(a byte) {
    p.A = a
    p.LsrAcc()
}

func (p *CPU) lsrMemory(a byte) {
    p.Memory.buffer[0] = a
    p.Lsr(0)
}

func TestLsrShiftsAccumulator(t *testing.T) {
    var p *CPU = new(CPU)

    p.lsrAccumulator(0x01)

    if p.A != 0x00 {
        t.Errorf("Shift left didn't work on the accumulator")
        t.Errorf("Expected 0x00, got %#02x", p.A)
        t.FailNow()
    }
}

func TestLsrShiftsMemoryLocation(t *testing.T) {
    var p *CPU = new(CPU)

    p.lsrMemory(0x01)

    if p.Memory.Read(0) != 0x00 {
        t.Errorf("Shift left didn't work on a memory location")
        t.Errorf("Expected 0x00, got %#02x", p.Memory.Read(0))
        t.FailNow()
    }
}

func TestLsrZeroFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.lsrAccumulator(0x00)

    if !p.Zero() {
        t.Errorf("Zero flag should be set when result is 0x00 (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestLsrCarryFlagSet(t *testing.T) {
    var p *CPU = new(CPU)

    p.lsrAccumulator(0x01)

    if !p.Carry() {
        t.Errorf("Carry flag should be set when shift won't fit (flags: %08b)", p.Flags)
        t.FailNow()
    }
}

func TestLsrCarryFlagUnset(t *testing.T) {
    var p *CPU = new(CPU)

    p.setCarryFlag(true)
    p.lsrAccumulator(0x02)

    if p.Carry() {
        t.Errorf("Carry flag should be unset when shift fits")
    }
}
