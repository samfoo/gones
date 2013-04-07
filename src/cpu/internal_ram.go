package cpu

type InternalRAM struct {
    buffer []byte
}

func NewInternalRAM() *InternalRAM {
    r := new(InternalRAM)

    r.buffer = make([]byte, 0x2000)

    return r
}

func (r *InternalRAM) normalize(location Address) Address {
    return location & 0x07ff
}

func (r *InternalRAM) Write(value byte, location Address) {
    r.buffer[r.normalize(location)] = value
}

func (r *InternalRAM) Read(location Address) byte {
    return r.buffer[r.normalize(location)]
}
