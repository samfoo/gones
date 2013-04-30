package cpu

func (p *CPU) HandleNMI() {
    p.push(byte(p.PC >> 8))
    p.push(byte(p.PC & 0x00ff))

    p.push(p.Flags)

    low := p.Memory.Read(0xfffa)
    high := p.Memory.Read(0xfffb)

    p.PC = (Address(high) << 8) | Address(low)
}
