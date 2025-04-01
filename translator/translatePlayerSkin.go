package translator

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (t *Translator) TranslatePlayerSkin(pk *packet.PlayerSkin) *packet.PlayerSkin {
	if pk.UUID == t.uuid {
		pk.UUID = StumpUUID
	}
	return pk
}
