package translator

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TranslatePlayerAuthInput(pk *packet.PlayerAuthInput, rid uint64) *packet.MoveActorAbsolute {
	//pk.Position[1] -= 1.62001

	y := pk.Position.Y()

	yDec := int32(y)

	onGround := y == float32(yDec)

	flags := byte(0)
	if onGround {
		flags |= packet.MoveFlagOnGround
	}

	return &packet.MoveActorAbsolute{
		EntityRuntimeID: rid,
		Flags:           flags,
		Position:        pk.Position,
		Rotation:        mgl32.Vec3{pk.Pitch, pk.Yaw, pk.Yaw},
	}
}
