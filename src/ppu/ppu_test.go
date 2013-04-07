package ppu

import (
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
