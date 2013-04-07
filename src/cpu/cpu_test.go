package cpu

import "testing"

func (p *CPU) address(mode AddressMode) Address {
    var location = mode(p)

    return location
}

func TestIndirectReference(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0x02, 0)
    p.Memory.Write(0x00, 1)
    p.Memory.Write(0xef, 2)
    p.Memory.Write(0xbe, 3)

    p.PC = 0x00

    location := p.Indirect()

    if location != 0xbeef {
        t.Errorf("Indirect reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, location)
    }
}

func TestIndirectReferenceOnPageBoundary(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0xbe, 0)
    p.Memory.Write(0xff, 1)
    p.Memory.Write(0x00, 2)
    p.Memory.Write(0xef, 0x00ff)

    p.PC = 0x01

    location := p.Indirect()

    if location != 0xbeef {
        t.Errorf("Indirect reference failed to emulate the notorious page boundary indirect bug")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, location)
    }
}

func TestImmediateReference(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0xbe, 0)

    p.PC = 0x00

    if p.address((*CPU).Immediate) != 0x0000 {
        t.Errorf("Immediate memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0x0000, p.address((*CPU).Immediate))
    }
}

func TestZeroPageReference(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0xbe, 0)

    p.PC = 0x00

    if p.address((*CPU).ZeroPage) != 0x00be {
        t.Errorf("Zero page memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0x00be, p.address((*CPU).ZeroPage))
    }
}

func TestZeroPageXReference(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0xbe, 0)

    p.X = 0x01
    p.PC = 0x00

    if p.address((*CPU).ZeroPageX) != 0x00bf {
        t.Errorf("Zero page X memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0x00bf, p.address((*CPU).ZeroPageX))
    }
}

func TestAbsoluteReference(t *testing.T) {
    var p *CPU = NewCPU()

    p.Memory.Write(0xef, 0)
    p.Memory.Write(0xbe, 1)

    p.PC = 0x00

    if p.address((*CPU).Absolute) != 0xbeef {
        t.Errorf("Absolute memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, p.address((*CPU).Absolute))
    }

    p.X = 0x01

    if p.address((*CPU).AbsoluteX) != 0xbef0 {
        t.Errorf("Absolute X memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbef0, p.address((*CPU).AbsoluteX))
    }

    p.Y = 0x02

    if p.address((*CPU).AbsoluteY) != 0xbef1 {
        t.Errorf("Absolute Y memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbef1, p.address((*CPU).AbsoluteY))
    }
}

func TestIndexedIndirect(t *testing.T) {
    var p *CPU = NewCPU()

    p.PC = 0x00
    p.X = 0x01

    p.Memory.Write(0x00, 0)
    p.Memory.Write(0xef, 1)
    p.Memory.Write(0xbe, 2)

    if p.address((*CPU).IndexedIndirect) != 0xbeef {
        t.Errorf("Indexed indirect memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, p.address((*CPU).IndexedIndirect))
    }

    p.X = 0x00
    p.Memory.Write(0xfe, 0)
    p.Memory.Write(0xfa, 0xff)
    p.Memory.Write(0xca, 0xfe)

    if p.address((*CPU).IndexedIndirect) != 0xfaca {
        t.Errorf("Indexed indirect memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xfaca, p.address((*CPU).IndexedIndirect))
    }
}

func TestIndirectIndexed(t *testing.T) {
    var p *CPU = NewCPU()

    p.PC = 0x00
    p.Y = 0x00

    p.Memory.Write(0x01, 0)
    p.Memory.Write(0xef, 1)
    p.Memory.Write(0xbe, 2)

    var location = p.IndirectIndexed()

    if location != 0xbeef {
        t.Errorf("Indirect indexed memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbeef, location)
    }

    p.Y = 0x01

    location = p.IndirectIndexed()

    if location != 0xbef0 {
        t.Errorf("Indirect indexed memory reference pointed to the wrong location")
        t.Errorf("Expected %#04x, got %#04x", 0xbef0, location)
    }
}

func TestReset(t *testing.T) {
    var p *CPU = NewCPU()

    p.Reset()

    if p.Flags != 0x24 {
        t.Errorf("Flags not set to 0x24")
    }

    if p.A != 0x00 {
        t.Errorf("A not set to 0x00")
    }

    if p.X != 0x00 {
        t.Errorf("X not set to 0x00")
    }

    if p.Y != 0x00 {
        t.Errorf("Y not set to 0x00")
    }

    if p.SP != 0xfd {
        t.Errorf("SP not set to 0xfd")
    }
}

func TestNewCPUShouldHaveRegistersZerod(t *testing.T) {
    var p *CPU = NewCPU()

    if p.A != 0x00 {
        t.Errorf("Accumulator is not zero'd")
        t.FailNow()
    }
}

