package cpu

import "testing"

func sta(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.A = value
    p.Sta(0x0000)

    return p
}

func stx(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.X = value
    p.Stx(0x0000)

    return p
}

func sty(value byte) (*CPU) {
    var p *CPU = NewCPU()

    p.Y = value
    p.Sty(0x0000)

    return p
}

func TestStoring(t *testing.T) {
    var p = sta(0xbe)
    if p.Memory.Read(0) != 0xbe {
        t.Errorf("Storing accumulator didn't work")
    }

    p = stx(0xbe)
    if p.Memory.Read(0) != 0xbe {
        t.Errorf("Storing X didn't work")
    }

    p = sty(0xbe)
    if p.Memory.Read(0) != 0xbe {
        t.Errorf("Storing Y didn't work")
    }
}
