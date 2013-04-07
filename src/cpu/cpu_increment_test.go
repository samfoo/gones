package cpu

import "testing"

func inc(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.Memory.Write(value, 0)
    p.Inc(0)

    return p
}

func inx(value byte) (*CPU) {
    var p *CPU = NewCPU()
    p.X = value
    p.Inx()
    return p
}

func iny(value byte) (*CPU) {
    var p *CPU = NewCPU()
    p.Y = value
    p.Iny()
    return p
}

func TestIncrementMutatesItem(t *testing.T) {
    var p = inc(0x01)
    if p.Memory.Read(0) != 0x02 {
        t.Errorf("Increment memory failed")
        t.Errorf("Expected 0x02, got %#02x", p.Memory.Read(0))
    }

    p = inx(0x01)
    if p.X != 0x02 {
        t.Errorf("Increment register X failed")
        t.Errorf("Expected 0x02, got %#02x", p.X)
    }

    p = iny(0x01)
    if p.Y != 0x02 {
        t.Errorf("Increment register Y failed")
        t.Errorf("Expected 0x02, got %#02x", p.Y)
    }
}

func TestIncrementNonZero(t *testing.T) {
    var validate = func(p *CPU) {
        if p.Zero() {
            t.Errorf("Zero flag shouldn't be set when increment isn't zero (flags: %08b)", p.Flags)
        }
    }

    validate(inc(0x01))
    validate(inx(0x01))
    validate(iny(0x01))
}

func TestIncrementZero(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Zero() {
            t.Errorf("Zero flag should be set when increment is zero (flags: %08b)", p.Flags)
        }
    }

    validate(inc(0xff))
    validate(inx(0xff))
    validate(iny(0xff))
}

func TestIncrementNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Negative flag should be set when increment is negative (flags: %08b)", p.Flags)
        }
    }

    validate(inc(0x80))
    validate(inx(0x80))
    validate(iny(0x80))
}
