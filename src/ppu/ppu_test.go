package ppu

import (
    "cpu"
    "testing"
    "github.com/stretchrcom/testify/assert"
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

    assert.Equal(t, p.VRAM, cpu.Address(0xbeef))
}
