package ppu

import (
    "cpu"
    "testing"
    "github.com/stretchrcom/testify/assert"
    "github.com/stretchrcom/testify/mock"
)

func TestSetBaseNametableAddress(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.BaseNametableAddress, Address(0x2000))

    ctrl.Set(0x01)
    assert.Equal(t, ctrl.BaseNametableAddress, Address(0x2400))

    ctrl.Set(0x02)
    assert.Equal(t, ctrl.BaseNametableAddress, Address(0x2800))

    ctrl.Set(0x03)
    assert.Equal(t, ctrl.BaseNametableAddress, Address(0x2c00))
}

func TestSetVRAMAddressInc(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.VRAMAddressInc, uint8(0))

    ctrl.Set(0x04)
    assert.Equal(t, ctrl.VRAMAddressInc, uint8(1))
}

func TestSetSpriteTableAddress(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.SpriteTableAddress, Address(0x0000))

    ctrl.Set(0x08)
    assert.Equal(t, ctrl.SpriteTableAddress, Address(0x1000))
}

func TestSetBackgroundTableAddress(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.BackgroundTableAddress, Address(0x0000))

    ctrl.Set(0x10)
    assert.Equal(t, ctrl.BackgroundTableAddress, Address(0x1000))
}

func TestSetSpriteSize(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.SpriteSize, uint8(0))

    ctrl.Set(0x20)
    assert.Equal(t, ctrl.SpriteSize, uint8(1))
}

func TestSetGenerateNMIOnVBlank(t *testing.T) {
    ctrl := new(Ctrl)

    ctrl.Set(0x00)
    assert.Equal(t, ctrl.GenerateNMIOnVBlank, false)

    ctrl.Set(0x80)
    assert.Equal(t, ctrl.GenerateNMIOnVBlank, true)
}

func TestSetMasks(t *testing.T) {
    m := new(Masks)

    assert.False(t, m.Grayscale)
    m.Set(0x01)
    assert.True(t, m.Grayscale)


    assert.False(t, m.ShowBackgroundLeft)
    m.Set(0x02)
    assert.True(t, m.ShowBackgroundLeft)

    assert.False(t, m.ShowSpritesLeft)
    m.Set(0x04)
    assert.True(t, m.ShowSpritesLeft)

    assert.False(t, m.ShowBackground)
    m.Set(0x08)
    assert.True(t, m.ShowBackground)

    assert.False(t, m.ShowSprites)
    m.Set(0x10)
    assert.True(t, m.ShowSprites)

    assert.False(t, m.IntenseReds)
    m.Set(0x20)
    assert.True(t, m.IntenseReds)

    assert.False(t, m.IntenseGreens)
    m.Set(0x40)
    assert.True(t, m.IntenseGreens)

    assert.False(t, m.IntenseBlues)
    m.Set(0x80)
    assert.True(t, m.IntenseBlues)
}

func TestSetAddr(t *testing.T) {
    p := NewPPU()

    p.SetAddr(0xbe)
    p.SetAddr(0xef)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0xbeef))
}

type MockReg struct {
    mock.Mock
}

func (m *MockReg) Set(val byte) {
    m.Mock.Called(val)
}

func TestPPUWriteCtrl(t *testing.T) {
    p := NewPPU()
    ctrl := new(MockReg)
    ctrl.On("Set", byte(0xff)).Return(nil)
    p.Ctrl = ctrl

    p.Write(0xff, PPUCTRL)
    ctrl.AssertCalled(t, "Set", byte(0xff))
}

func TestPPUWriteMask(t *testing.T) {
    p := NewPPU()
    masks := new(MockReg)
    masks.On("Set", byte(0xff)).Return(nil)
    p.Masks = masks

    p.Write(0xff, PPUMASK)
    masks.AssertCalled(t, "Set", byte(0xff))
}

func TestPPUWriteOAMAddr(t *testing.T) {
    p := NewPPU()

    p.Write(0xff, OAMADDR)

    assert.Equal(t, p.OAMAddr, byte(0xff))
}

func TestPPUWriteOAMData(t *testing.T) {
    p := NewPPU()

    p.OAMAddr = 0xff
    p.Write(0xbe, OAMDATA)

    assert.Equal(t, p.OAMRAM[0xff], byte(0xbe))
    assert.Equal(t, p.OAMAddr, byte(0x00))

    p.Write(0xbe, OAMDATA)
    assert.Equal(t, p.OAMAddr, byte(0x01))
}

func TestPPUWritePPUAddr(t *testing.T) {
    p := NewPPU()

    p.Write(0xbe, PPUADDR)
    p.Write(0xef, PPUADDR)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0xbeef))
}

func TestPPUWritePPUData(t *testing.T) {
    p := NewPPU()
    p.VRAMAddr = 0x0000

    p.Write(0xbe, PPUDATA)

    assert.Equal(t, p.Memory.Read(0x0000), byte(0xbe))
    assert.Equal(t, p.VRAMAddr, cpu.Address(0x0001))
}
