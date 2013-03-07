package cpu

import "testing"

func cmp(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.A = first
    p.Memory[0] = second
    p.Cmp(0)

    return p
}

func cpx(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.X = first
    p.Memory[0] = second
    p.Cpx(0)

    return p
}

func cpy(first byte, second byte) (*CPU) {
    var p *CPU = new(CPU)

    p.Y = first
    p.Memory[0] = second
    p.Cpy(0)

    return p
}

func TestCompareGreaterThan(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Carry() {
            t.Errorf("Carry flag should be set when {register} > {value} but isn't (flags: %08b)", p.Flags)
        }
    }
    validate(cmp(0x01, 0x00))
    validate(cpy(0x01, 0x00))
    validate(cpx(0x01, 0x00))
}

func TestCompareEqual(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Carry() {
            t.Errorf("Carry flag should be set when {register} == {value} but isn't (flags: %08b)", p.Flags)
        }

        if !p.Zero() {
            t.Errorf("Zero flag should be set when {register} == {value} but isn't (flags: %08b)", p.Flags)
        }
    }

    validate(cmp(0x01, 0x01))
    validate(cpy(0x01, 0x01))
    validate(cpx(0x01, 0x01))
}

func TestCompareNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Negative flag should be set when {register} < {value} but isn't (flags: %08b)", p.Flags)
        }
    }

    validate(cmp(0x00, 0x01))
    validate(cpy(0x00, 0x01))
    validate(cpx(0x00, 0x01))
}
