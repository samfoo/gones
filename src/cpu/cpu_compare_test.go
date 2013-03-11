package cpu

import "testing"

func cmp(p *CPU, first byte, second byte) (*CPU) {
    p.A = first
    p.Memory[0] = second
    p.Cmp(0)

    return p
}

func cpx(p *CPU, first byte, second byte) (*CPU) {
    p.X = first
    p.Memory[0] = second
    p.Cpx(0)

    return p
}

func cpy(p *CPU, first byte, second byte) (*CPU) {
    p.Y = first
    p.Memory[0] = second
    p.Cpy(0)

    return p
}

func TestCompareNegativeUnsetsCarryFlag(t *testing.T) {
    var validate = func(p *CPU) {
        if p.Carry() {
            t.Errorf("Carry flag should be unset when {register} < {value} but was set")
        }
    }
    var p *CPU = new(CPU)
    p.setCarryFlag(true)
    validate(cmp(p, 0x00, 0x01))
    p.setCarryFlag(true)
    validate(cpx(p, 0x00, 0x01))
    p.setCarryFlag(true)
    validate(cpy(p, 0x00, 0x01))
}

func TestCompareGreaterThan(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Carry() {
            t.Errorf("Carry flag should be set when {register} > {value} but isn't (flags: %08b)", p.Flags)
        }
    }
    validate(cmp(new(CPU), 0x01, 0x00))
    validate(cpy(new(CPU), 0x01, 0x00))
    validate(cpx(new(CPU), 0x01, 0x00))
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

    validate(cmp(new(CPU), 0x01, 0x01))
    validate(cpy(new(CPU), 0x01, 0x01))
    validate(cpx(new(CPU), 0x01, 0x01))
}

func TestCompareNegative(t *testing.T) {
    var validate = func(p *CPU) {
        if !p.Negative() {
            t.Errorf("Negative flag should be set when {register} < {value} but isn't (flags: %08b)", p.Flags)
        }
    }

    validate(cmp(new(CPU), 0x00, 0x01))
    validate(cpy(new(CPU), 0x00, 0x01))
    validate(cpx(new(CPU), 0x00, 0x01))
}
