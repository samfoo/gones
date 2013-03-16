package cpu

import "testing"

func dcp(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Memory.buffer[0] = value
    p.Dcp(0)

    return p
}

func lax(value byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A, p.X = 0, 0
    p.Memory.buffer[0] = value
    p.Lax(0)

    return p
}

func aac(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = first
    p.Memory.buffer[0] = second
    p.Aac(0)

    return p
}

func sax(a byte, x byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = a
    p.X = x
    p.Sax(0)

    return p
}

func TestAacSetsAccummulator(t *testing.T) {
    p := aac(0xff, 0x44)

    if p.A != 0x44 {
        t.Errorf("AAC didn't and with the accumulator")
        t.Errorf("Expected %#02x, got %#02x", 0x44, p.A)
    }
}

func TestAacSetsCarryFlagWhenResultIsNegative(t *testing.T) {
    p := aac(0xff, 0xff)

    if !p.Carry() {
        t.Errorf("Carry flag not set on negative result")
    }
}

func TestAacSetsZeroFlagWhenResultIsZero(t *testing.T) {
    p := aac(0x00, 0x00)

    if !p.Zero() {
        t.Errorf("Zero flag not set on zero result")
    }
}

func TestLaxSetsAAndX(t *testing.T) {
    p := lax(0xbb)

    if p.X != 0xbb || p.A != 0xbb {
        t.Errorf("Lax didn't load both A and X registers")
    }
}

func TestLaxSetsZeroFlag(t *testing.T) {
    p := lax(0x00)

    if !p.Zero() {
        t.Errorf("Lax didn't set the zero flag on zero result")
    }
}

func TestLaxSetsNegativeFlag(t *testing.T) {
    p := lax(0x80)

    if !p.Negative() {
        t.Errorf("Lax didn't set the negative flag on zero result")
    }
}

func TestSaxAndsAAndX(t *testing.T) {
    p := sax(0xff, 0x45)

    if p.Memory.Read(0) != 0x45 {
        t.Errorf("Sax didn't and X and A into memory")
    }
}

func TestSaxDoesntSetZeroFlag(t *testing.T) {
    p := sax(0xff, 0x00)

    if p.Zero() {
        t.Errorf("Sax set the zero flag on zero result")
    }
}

func TestSaxDoesntSetNegativeFlag(t *testing.T) {
    p := sax(0x80, 0x80)

    if p.Negative() {
        t.Errorf("Sax set the negative flag on zero result")
    }
}

func TstDcpSubs1FromMemory(t *testing.T) {
    p := dcp(0x02)

    if p.Memory.Read(0) != 0x01 {
        t.Errorf("Dcp didn't decrement memory")
    }
}
