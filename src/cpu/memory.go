package cpu

import (
    "fmt"
    "errors"
)

type Mount struct {
    From Address
    To Address
    Device Mountable
}

type Memory struct {
    Mounts []Mount
}

type Mountable interface {
    Read(Address) byte
    Write(byte, Address)
}

type RAM struct {
    buffer []byte
}

func NewRAM() *RAM {
    r := new(RAM)

    r.buffer = make([]byte, 0x10000)

    return r
}

func normalize(location Address) Address {
    if location <= 0x1fff {
        return location & 0x07ff
    }

    return location
}

func (r *RAM) Write(value byte, location Address) {
    r.buffer[normalize(location)] = value
}

func (r *RAM) Read(location Address) byte {
    return r.buffer[normalize(location)]
}

func NewMemory() *Memory {
    m := new(Memory)

    // Surely we'll never have more than 10 mounts, right?
    m.Mounts = make([]Mount, 0, 10)

    m.Mount(NewRAM(), 0x0000, 0xffff)

    return m
}

func (m *Memory) Mount(dev Mountable, from Address, to Address) error {
    for i := range m.Mounts {
        other := m.Mounts[i]

        if (other.From <= from && other.To >= from) ||
            (other.From <= to && other.To >= to) {
            return errors.New("A device already exists at that mount point")
        }
    }

    m.Mounts = append(m.Mounts, Mount { from, to, dev })

    return nil
}

func (m *Memory) findDevice(location Address) Mountable {
    for i := range m.Mounts {
        mount := m.Mounts[i]

        if mount.From <= location && mount.To >= location {
            return mount.Device
        }
    }

    return nil
}

func (m *Memory) Range(from Address, to Address) []byte {
    results := make([]byte, to - from)

    for i := range results {
        results[i] = m.Read(from + Address(i))
    }

    return results
}

func (m *Memory) Read(location Address) byte {
    dev := m.findDevice(location)

    if dev != nil {
        return dev.Read(location)
    }

    panic(fmt.Sprintf("Read occurred an unmounted memory location %#04x", location))
}

func (m *Memory) Copy(data []byte, location Address) {
    for i := range data {
        m.Write(data[i], location + Address(i))
    }
}

func (m *Memory) Write(value byte, location Address) {
    dev := m.findDevice(location)

    if dev != nil {
        dev.Write(value, location)
    } else {
        panic(fmt.Sprintf("Write occurred at an unmounted memory location %#02x -> %#04x", value, location))
    }
}
