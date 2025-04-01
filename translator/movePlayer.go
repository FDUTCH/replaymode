package translator

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (t *Translator) TranslateMovePlayer(pk *packet.MovePlayer) []packet.Packet {

	if t.rid == pk.EntityRuntimeID {

		return []packet.Packet{&packet.MovePlayer{
			EntityRuntimeID: Rid,
			Position:        pk.Position,
			Pitch:           pk.Pitch,
			Yaw:             pk.Yaw,
			HeadYaw:         pk.HeadYaw,
			Mode:            pk.Mode,
			TeleportCause:   pk.TeleportCause,
		}, pk}
		//if t.tolerance.Before(time.Now()) {
		//} else {
		//}
		//pk.EntityRuntimeID = Rid
		//return []packet.Packet{
		//	pk,
		//}

	}
	return []packet.Packet{pk}
}
