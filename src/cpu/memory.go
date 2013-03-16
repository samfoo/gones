package cpu

type RAM struct {
    buffer [0x10000]byte
}

func (m *RAM) Read(location Address) byte {
    return m.buffer[location]
}

func (m *RAM) Write(value byte, location Address) {
    m.buffer[location] = value
}
