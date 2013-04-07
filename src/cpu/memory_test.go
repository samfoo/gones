package cpu

import (
    "testing"
    "github.com/stretchrcom/testify/assert"
)

type MockDevice struct {
    buffer [10]byte
}

func (d *MockDevice) Read(location Address) byte {
    return d.buffer[location]
}

func (d *MockDevice) Write(val byte, location Address) {
    d.buffer[location] = val
}

func TestDoubleMountFails(t *testing.T) {
    var m = NewMemory()
    m.Mounts = []Mount{}

    first := new(MockDevice)
    second := new(MockDevice)

    assert.Nil(t, m.Mount(first, 0x0500, 0x0fff))
    assert.NotNil(t, m.Mount(second, 0x0400, 0x600))
    assert.NotNil(t, m.Mount(second, 0x0f00, 0x1000))
    assert.NotNil(t, m.Mount(second, 0x0501, 0x0f00))
}

func TestMountReadsFromDevice(t *testing.T) {
    var m = NewMemory()
    m.Mounts = []Mount{}

    dev := new(MockDevice)
    dev.buffer[0] = 0xdd
    m.Mount(dev, 0x0000, 0x0100)

    assert.Equal(t, m.Read(0x0000), byte(0xdd))
}

func TestMountWritesToDevice(t *testing.T) {
    var m = NewMemory()
    m.Mounts = []Mount{}

    dev := new(MockDevice)
    m.Mount(dev, 0x0000, 0x0100)

    m.Write(0xff, 0x0000)
    assert.Equal(t, dev.buffer[0x0000], byte(0xff))
}
