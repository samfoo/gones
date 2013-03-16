package cpu

type RAM struct {
    buffer [0x10000]byte
}

func normalize(location Address) Address {
    switch {
        case location <= 0x1fff:
            return location & 0x07ff
    }

    return location
}

func (m *RAM) Read(location Address) byte {
    return m.buffer[normalize(location)]
}

func (m *RAM) Write(value byte, location Address) {
    m.buffer[normalize(location)] = value
}
