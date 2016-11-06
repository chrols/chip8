package chip8

import "testing"

func TestReset(t *testing.T) {
	cpu := Cpu{}
	cpu.Reset()
	if cpu.ProgramCounter != 0x200 {
		t.Error("Expected 0x200, got", cpu.ProgramCounter)
	}
}
