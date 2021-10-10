package forzaprotocol

import (
	"bytes"
	"encoding/binary"
)

// NOTE: unsafe.Sizeof returns a packet size of 324, which is wrong. Not sure why as of yet.
var PacketSize = 323 // unsafe.Sizeof(Packet{})

// Packet represents a single packet from the UDP dataport.
type Packet struct {
	Running                          int32 // 0 in menus, 1 when not.
	Timestamp                        uint32
	EngineMaxRpm                     float32
	EngineIdleRpm                    float32
	CurrentEngineRpm                 float32
	AccelerationX                    float32
	AccelerationY                    float32
	AccelerationZ                    float32
	VelocityX                        float32
	VelocityY                        float32
	VelocityZ                        float32
	AngularVelocityX                 float32
	AngularVelocityY                 float32
	AngularVelocityZ                 float32
	Yaw                              float32
	Pitch                            float32
	Roll                             float32
	NormalizedSuspensionFrontLeft    float32 // 0.0 - relaxed, 1.0 - compressed.
	NormalizedSuspensionFrontRight   float32
	NormalizedSuspensionRearLeft     float32
	NormalizedSuspensionRearRight    float32
	TireSlipRatioFrontLeft           float32 // 0.0 - max grip, 1.0+ - no grip.
	TireSlipRatioFrontRight          float32
	TireSlipRatioRearLeft            float32
	TireSlipRatioRearRight           float32
	WheelRotationSpeedFrontLeft      float32 // Radians p/sec.
	WheelRotationSpeedFrontRight     float32
	WheelRotationSpeedRearLeft       float32
	WheelRotationSpeedRearRight      float32
	WheelOnRumbleStripFrontLeft      int32 // 0 when not on rumble strip, 1 when on.
	WheelOnRumbleStripFrontRight     int32
	WheelOnRumbleStripRearLeft       int32
	WheelOnRumbleStripRearRight      int32
	WheelInPuddleDepthFrontLeft      float32 // 0.0 - no puddle, 1.0 - deep puddle.
	WheelInPuddleDepthFrontRight     float32
	WheelInPuddleDepthRearLeft       float32
	WheelInPuddleDepthRearRight      float32
	SurfaceRumbleFrontLeft           float32
	SurfaceRumbleFrontRight          float32
	SurfaceRumbleRearLeft            float32
	SurfaceRumbleRearRight           float32
	TireSlipAngleFrontLeft           float32 // 0.0 - max grip, 1.0+ - no grip.
	TireSlipAngleFrontRight          float32
	TireSlipAngleRearLeft            float32
	TireSlipAngleRearRight           float32
	TireCombinedSlipFrontLeft        float32 // 0.0 - max grip, 1.0+ - no grip.
	TireCombinedSlipFrontRight       float32
	TireCombinedSlipRearLeft         float32
	TireCombinedSlipRearRight        float32
	SuspensionTravelMetersFrontLeft  float32
	SuspensionTravelMetersFrontRight float32
	SuspensionTravelMetersRearLeft   float32
	SuspensionTravelMetersRearRight  float32
	CarID                            int32
	CarClass                         int32 // 0: D, 1: C, 2: B, 3: A, 4: S1, 5: S2, 6: R, 7: X.
	CarPerformanceIndex              int32
	Drivetrain                       int32 // 0: FWD, 1: RWD, 2: AWD.
	NumCylinders                     int32
	Unknown                          [12]byte
	PositionX                        float32
	PositionY                        float32
	PositionZ                        float32
	Speed                            float32 // Meters per second.
	Power                            float32 // Watts.
	Torque                           float32 // Newtonmeter.
	TireTempFrontLeft                float32
	TireTempFrontRight               float32
	TireTempRearLeft                 float32
	TireTempRearRight                float32
	Boost                            float32
	Fuel                             float32
	DistanceTraveled                 float32
	BestLap                          float32
	LastLap                          float32
	CurrentLap                       float32
	CurrentRaceTime                  float32
	LapNumber                        uint16
	RacePosition                     uint8
	Accel                            uint8
	Brake                            uint8
	Clutch                           uint8
	Handbrake                        uint8
	Gear                             uint8
	Steer                            int8
	NormalizedDrivingLine            int8
	NormalizedAIBrakeDifference      int8
}

// FromBytes converts a raw packet into this struct.
func (p *Packet) FromBytes(data []byte) {
	binary.Read(bytes.NewReader(data), binary.LittleEndian, p)
}

// ToBytes converts a packet to raw bytes.
func (p *Packet) ToBytes() []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.LittleEndian, *p)
	return buffer.Bytes()
}
