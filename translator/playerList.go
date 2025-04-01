package translator

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func (t *Translator) TranslatePlayerList(pk *packet.PlayerList) *packet.PlayerList {
	for i, val := range pk.Entries {
		if val.UUID == t.uuid {
			pk.Entries[i].UUID = StumpUUID
		}
	}
	return pk
}
