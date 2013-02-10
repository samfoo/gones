package cpu

type Address uint16

type CPU struct {
    A, X, Y, SP, Flags byte
    PC Address
    Memory [0x10000]byte
}

func (p *CPU) Adc(location Address) {
    val := p.Memory[location]
    old := p.A

    p.A += val

    if p.A < old {
        // There was a carry in unsigned arithmetic
        p.Flags |= 0x01
    }
}
