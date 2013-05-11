package ppu

import (
    "cpu"
    "testing"
    "github.com/stretchrcom/testify/assert"
    "github.com/stretchrcom/testify/mock"
)

type MockBus struct {
    mock.Mock
}

func (m *MockBus) Interrupt(kind int) {
    m.Mock.Called(kind)
}

func (m *MockBus) Cancel(kind int) {
    m.Mock.Called(kind)
}

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

func TestWriteVRAMAddr(t *testing.T) {
    p := NewPPU()

    p.AddressLatch = true
    p.WriteVRAMAddr(0xbe)
    p.WriteVRAMAddr(0xef)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0xbeef))
}

func TestPPUWriteCtrl(t *testing.T) {
    p := NewPPU()

    assert.Equal(t, p.Ctrl.VRAMAddressInc, uint8(VRAM_INC_ACROSS))
    p.Write(0xff, PPUCTRL)
    assert.Equal(t, p.Ctrl.VRAMAddressInc, uint8(0x01))
}

func TestPPUWriteMask(t *testing.T) {
    p := NewPPU()

    assert.Equal(t, p.Masks.Grayscale, false)
    p.Write(0xff, PPUMASK)
    assert.Equal(t, p.Masks.Grayscale, true)
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

    p.AddressLatch = true
    p.Write(0xbe, PPUADDR)
    p.Write(0xef, PPUADDR)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0xbeef))
}

func TestPPUWritePPUData(t *testing.T) {
    p := NewPPU()
    p.VRAMAddr = 0x0000
    p.Memory.Mount(NewPatterntable(make([]byte, 0x1000)), 0x0000, 0x0fff)

    p.Write(0xbe, PPUDATA)

    assert.Equal(t, p.Memory.Read(0x0000), byte(0xbe))
    assert.Equal(t, p.VRAMAddr, cpu.Address(0x0001))
}

func TestStatusSerializesFlags(t *testing.T) {
    var s = new(Status)
    assert.Equal(t, s.Value(), byte(0x00))

    s.SpriteOverflow = true
    assert.Equal(t, s.Value(), byte(0x20))

    s = new(Status)
    s.Sprite0Hit = true
    assert.Equal(t, s.Value(), byte(0x40))

    s = new(Status)
    s.VBlankStarted = true
    assert.Equal(t, s.Value(), byte(0x80))
}

func TestReadingStatusSetsTheAddressLatch(t *testing.T) {
    p := NewPPU()

    assert.Equal(t, p.AddressLatch, false)
    p.Read(PPUSTATUS)
    assert.Equal(t, p.AddressLatch, true)
}

func TestPPUReadAt2002ReadsStatus(t *testing.T) {
    p := NewPPU()

    p.Status.SpriteOverflow = true
    assert.Equal(t, p.Read(PPUSTATUS), byte(0x20))
}

func TestPPUReadAt2004ReadsOAMData(t *testing.T) {
    p := NewPPU()

    p.OAMAddr = 0x00
    p.OAMRAM[p.OAMAddr] = 0xbe

    assert.Equal(t, p.Read(OAMDATA), byte(0xbe))
    assert.Equal(t, p.OAMAddr, byte(0x00))
}

func TestPPUReadAt2007ReadsData(t *testing.T) {
    p := NewPPU()
    p.Memory.Mount(NewPatterntable(make([]byte, 0x1000)), 0x0000, 0x0fff)

    p.VRAMAddr = cpu.Address(0x0000)
    p.Memory.Write(0xbe, p.VRAMAddr)

    assert.Equal(t, p.Read(PPUDATA), byte(0xbe))
}

func TestPPUDATAReadIncrementsVRAMAddrCorrectly(t *testing.T) {
    p := NewPPU()
    p.Memory.Mount(NewPatterntable(make([]byte, 0x1000)), 0x0000, 0x0fff)

    p.VRAMAddr = cpu.Address(0x0000)
    p.Ctrl.VRAMAddressInc = VRAM_INC_ACROSS
    p.Read(PPUDATA)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0x0001))

    p.VRAMAddr = cpu.Address(0x0000)
    p.Ctrl.VRAMAddressInc = VRAM_INC_DOWN
    p.Read(PPUDATA)

    assert.Equal(t, p.VRAMAddr, cpu.Address(0x0020))
}

func TestPPUStepOnScanlineNegFirstCycleClearsSpriteOverflow(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Status.SpriteOverflow = true
    p.Cycle = 1

    p.Step()

    assert.False(t, p.Status.SpriteOverflow)
}

func TestPPUStepOnScanlineNegFirstCycleClearsSprite0Hit(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Status.Sprite0Hit = true
    p.Cycle = 1

    p.Step()

    assert.False(t, p.Status.Sprite0Hit)
}

func TestPPUStepOnScanlineNegNotFirstCycleDoesntClearSpriteOverflow(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Status.SpriteOverflow = true
    p.Cycle = 100

    p.Step()

    assert.True(t, p.Status.SpriteOverflow)
}

func TestPPUStepOnScanlineNegNotFirstCycleDoesntClearSprite0Hit(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Status.Sprite0Hit = true
    p.Cycle = 100

    p.Step()

    assert.True(t, p.Status.Sprite0Hit)
}

func TestPPUStepOnScanlineNegAndFirstCycleClearVBlankStarted(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Cycle = 1
    p.VBlankStarted = true

    p.Step()

    assert.False(t, p.Status.VBlankStarted)
}

func TestPPUStepAfterPostrenderScanlineFirstCycleSetsVBlankStarted(t *testing.T) {
    p := NewPPU()
    p.Scanline = POSTRENDER_SCANLINE + 1
    p.Cycle = 1
    p.VBlankStarted = false

    p.Step()

    assert.True(t, p.Status.VBlankStarted)
}

func TestPPUStepIncrementsCycle(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Cycle = 1

    p.Step()

    assert.Equal(t, p.Cycle, 2)
}

func TestPPUIfLastCycleThenIncrementScanline(t *testing.T) {
    p := NewPPU()
    p.Scanline = -1
    p.Cycle = LAST_CYCLE

    p.Step()

    assert.Equal(t, p.Scanline, 0)
}

func TestPPUIfLastScanlineAndLastCycleScanlineShouldReset(t *testing.T) {
    p := NewPPU()
    p.Scanline = VBLANK_SCANLINE
    p.Cycle = LAST_CYCLE

    p.Step()

    assert.Equal(t, p.Scanline, PRERENDER_SCANLINE)
}

func TestPPUIfLastScanlineAndLastCycleIncrementFrameCount(t *testing.T) {
    p := NewPPU()
    p.Scanline = VBLANK_SCANLINE
    p.Cycle = LAST_CYCLE

    assert.Equal(t, p.Frame, 0)
    p.Step()
    assert.Equal(t, p.Frame, 1)
}

func TestPPUShortensPrerenderScanlineOnOddFrames(t *testing.T) {
    p := NewPPU()
    p.Scanline = FIRST_VISIBLE_SCANLINE
    p.Cycle = 0
    p.Frame = 1
    p.Masks.ShowBackground = true

    p.Step()
    assert.Equal(t, p.Cycle, 2)
}

func TestPPUDoesntShortenPrerenderOnEvenFrames(t *testing.T) {
    p := NewPPU()
    p.Scanline = PRERENDER_SCANLINE
    p.Cycle = LAST_CYCLE - 1
    p.Frame = 0
    p.Masks.ShowBackground = true

    p.Step()
    assert.Equal(t, p.Cycle, LAST_CYCLE)
    assert.Equal(t, p.Scanline, PRERENDER_SCANLINE)
}

func TestPPUDoesntShortenPrerenderWhenRenderingIsDisabled(t *testing.T) {
    p := NewPPU()
    p.Scanline = PRERENDER_SCANLINE
    p.Cycle = LAST_CYCLE - 1
    p.Frame = 1
    p.Masks.ShowBackground = false
    p.Masks.ShowSprites = false

    p.Step()
    assert.Equal(t, p.Cycle, LAST_CYCLE)
    assert.Equal(t, p.Scanline, PRERENDER_SCANLINE)
}

func TestPPUGeneratesNMIOnVBlankWhenEnabled(t *testing.T) {
    p := NewPPU()
    p.Ctrl.GenerateNMIOnVBlank = true
    p.Status.VBlankStarted = true
    bus := new(MockBus)
    bus.On("Interrupt", cpu.NMI).Return(nil)
    p.Bus = bus

    // Set the state to where vblank would occur
    p.Cycle = 1
    p.Scanline = POSTRENDER_SCANLINE + 1

    p.Step()
    bus.Mock.AssertCalled(t, "Interrupt", cpu.NMI)
}

func TestPPUDoesntGenerateNMIOnVBlankWhenNotEnabled(t *testing.T) {
    p := NewPPU()
    p.Ctrl.GenerateNMIOnVBlank = false
    bus := new(MockBus)
    p.Bus = bus

    // Set the state to where vblank would occur
    p.Cycle = 1
    p.Scanline = POSTRENDER_SCANLINE + 1

    p.Step()
    bus.Mock.AssertNotCalled(t, "Interrupt", cpu.NMI)
}

func TestPPUReadingPPUSTATUSSetsVBlanksStartedToFalse(t *testing.T) {
    p := NewPPU()
    p.Status.VBlankStarted = true

    p.Read(PPUSTATUS)

    assert.False(t, p.Status.VBlankStarted)
}

func TestPPUTogglingGenerateNMIOnVBlankWhileInVBlankGeneratesMultipleInterrupts(t *testing.T) {
    p := NewPPU()
    bus := new(MockBus)
    bus.On("Interrupt", cpu.NMI).Return(nil)
    p.Status.VBlankStarted = true
    p.Bus = bus

    p.Write(0x00, PPUCTRL)
    p.Write(0x80, PPUCTRL)

    p.Write(0x00, PPUCTRL)
    p.Write(0x80, PPUCTRL)

    bus.Mock.AssertNumberOfCalls(t, "Interrupt", 2)
}

func TestPPUSettingGenerateNMIOnVBlankWhenAlreadySetDoesntGeneralMultipleInterrupts(t *testing.T) {
    p := NewPPU()
    bus := new(MockBus)
    bus.On("Interrupt", cpu.NMI).Return(nil)
    p.Status.VBlankStarted = true
    p.Ctrl.GenerateNMIOnVBlank = true
    p.Bus = bus

    p.Write(0x80, PPUCTRL)
    p.Write(0x80, PPUCTRL)

    bus.Mock.AssertNotCalled(t, "Interrupt", cpu.NMI)
}

func TestNormalizeMirrorsPPURegisters(t *testing.T) {
    p := NewPPU()

    assert.Equal(t, p.normalize(0x0008), cpu.Address(0x0000))
    assert.Equal(t, p.normalize(0x1456), cpu.Address(0x0006))
}

func TestReadingPPUSTATUSOneCycleBeforeVBlank(t *testing.T) {
    p := NewPPU()
    bus := new(MockBus)
    bus.On("Interrupt", cpu.NMI).Return(nil)
    p.Bus = bus
    p.Scanline = POSTRENDER_SCANLINE + 1
    p.Cycle = 1

    assert.Equal(t, p.Read(PPUSTATUS) & 0x80, byte(0x00))

    p.Step()

    assert.Equal(t, p.Status.VBlankStarted, false)
    bus.Mock.AssertNotCalled(t, "Interrupt", cpu.NMI)
}

func TestReadingPPUSTATUSOnSameCycleAsVBlank(t *testing.T) {
    p := NewPPU()
    bus := new(MockBus)
    bus.On("Cancel", cpu.NMI).Return(nil)
    p.Bus = bus
    p.Scanline = POSTRENDER_SCANLINE + 1
    p.Cycle = 2

    assert.Equal(t, p.Read(PPUSTATUS) & 0x80, byte(0x80))

    p.Step()

    assert.Equal(t, p.Status.VBlankStarted, false)
    assert.Equal(t, p.suppressVBlankStarted, true)
    bus.Mock.AssertCalled(t, "Cancel", cpu.NMI)
}
