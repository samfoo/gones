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

func NewMemory() *Memory {
    m := new(Memory)

    // Surely we'll never have more than 10 mounts, right?
    m.Mounts = make([]Mount, 0, 10)

    return m
}

type RAM struct {
    buffer []byte
}

func NewRAM(size uint16) *RAM {
    ram := new(RAM)
    ram.buffer = make([]byte, size)

    return ram
}

func (r *RAM) Read(location Address) byte {
    return r.buffer[location]
}

func (r *RAM) Write(val byte, location Address) {
    r.buffer[location] = val
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

func (m *Memory) findMount(location Address) *Mount {
    for i := range m.Mounts {
        mount := m.Mounts[i]

        if mount.From <= location && mount.To >= location {
            return &mount
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
    mount := m.findMount(location)

    if mount != nil {
        normalized := location - mount.From
        return mount.Device.Read(normalized)
    }

    panic(fmt.Sprintf("Read occurred an unmounted memory location %#04x", location))
}

func (m *Memory) Copy(data []byte, location Address) {
    for i := range data {
        m.Write(data[i], location + Address(i))
    }
}

func (m *Memory) Write(value byte, location Address) {
    mount := m.findMount(location)

    if mount != nil {
        normalized := location - mount.From
        mount.Device.Write(value, normalized)
    } else {
        panic(fmt.Sprintf("Write occurred at an unmounted memory location %#02x -> %#04x", value, location))
    }
}
