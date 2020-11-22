package forza

import (
	"fmt"
	"testing"
)

func TestDeserialize(t *testing.T) {
	dataPaused := []byte{0x00, 0x00, 0x00, 0x00, 0x82, 0x9D, 0xC8, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	dataRunning := []byte{0x01, 0x00, 0x00, 0x00, 0x14, 0x9D, 0xC8, 0x08, 0xF8, 0xBF, 0xDA, 0x45, 0xF8, 0xFF, 0x47, 0x44, 0x14, 0x17, 0x48, 0x44, 0x91, 0x64, 0xB6, 0xBC, 0xB6, 0x71, 0x4F, 0x3E, 0xAA, 0xC6, 0xF5, 0xBD, 0x00, 0x53, 0xEE, 0x3B, 0x10, 0x9D, 0xA5, 0x3C, 0x51, 0x36, 0x16, 0x40, 0x82, 0xD0, 0xCB, 0xBC, 0xE8, 0xA3, 0x39, 0xBA, 0x98, 0xD9, 0x01, 0xBC, 0xA1, 0xB3, 0x3C, 0x40, 0x47, 0x63, 0xA1, 0xBC, 0x2A, 0xE2, 0x16, 0xBD, 0x27, 0x3F, 0xAC, 0x3E, 0x63, 0x86, 0xA1, 0x3E, 0xB9, 0x22, 0x75, 0x3E, 0x7E, 0xCB, 0x8B, 0x3E, 0x4F, 0xE2, 0xAB, 0x39, 0x82, 0xED, 0xA6, 0x39, 0xFD, 0x96, 0xAF, 0x39, 0x8D, 0x28, 0xA9, 0x39, 0x36, 0xFF, 0xE5, 0x40, 0x2E, 0x18, 0xE6, 0x40, 0x8F, 0xFD, 0xE5, 0x40, 0x54, 0x01, 0xE6, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x90, 0xC2, 0xF5, 0x3D, 0x90, 0xC2, 0xF5, 0x3D, 0x90, 0xC2, 0xF5, 0x3D, 0x90, 0xC2, 0xF5, 0x3D, 0x32, 0xBD, 0xC9, 0xBB, 0x97, 0x46, 0x97, 0xBB, 0x84, 0xD6, 0x0D, 0xBC, 0xD3, 0x70, 0xF3, 0xBB, 0x5E, 0x06, 0xCA, 0x3B, 0x95, 0xA2, 0x97, 0x3B, 0xAD, 0xF1, 0x0D, 0x3C, 0x91, 0xAB, 0xF3, 0x3B, 0x80, 0x48, 0x37, 0x3B, 0x00, 0xB8, 0xCF, 0x3A, 0x00, 0xD1, 0x77, 0x3B, 0xE0, 0x2C, 0xF9, 0x3B, 0x01, 0x02, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x6F, 0x03, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x1C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B, 0x6B, 0xFC, 0xC4, 0xBF, 0xEE, 0x9B, 0x43, 0x35, 0x8F, 0x12, 0xC3, 0xED, 0x37, 0x16, 0x40, 0xD9, 0xC0, 0x01, 0xC0, 0xA1, 0x28, 0xC6, 0xBC, 0x9A, 0x83, 0xD0, 0x42, 0x48, 0x81, 0xD0, 0x42, 0x29, 0x03, 0xDD, 0x42, 0x29, 0x03, 0xDD, 0x42, 0x64, 0x66, 0x30, 0xC1, 0x00, 0x00, 0x80, 0x3F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xE7, 0x9C, 0xA0, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}

	fmt.Printf("Paused: %x\n\n", dataPaused)
	fmt.Printf("Running: %x\n\n", dataRunning)

	var packet DataPacket
	packet.FromBytes(dataPaused)
	fmt.Printf("Paused: %+v\n\n", packet)

	packet.FromBytes(dataRunning)
	fmt.Printf("Running: %+v\n\n", packet)
}
