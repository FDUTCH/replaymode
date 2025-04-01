package translator

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (t *Translator) TranslateInventoryContent(pk *packet.InventoryContent) []packet.Packet {
	switch pk.WindowID {
	case protocol.WindowIDOffHand:
		return []packet.Packet{&packet.MobEquipment{
			EntityRuntimeID: t.rid,
			NewItem:         pk.Content[0],
			WindowID:        protocol.WindowIDOffHand,
		}}
	case protocol.WindowIDArmour:
		return []packet.Packet{&packet.MobArmourEquipment{
			EntityRuntimeID: t.rid,
			Helmet:          pk.Content[0],
			Chestplate:      pk.Content[1],
			Leggings:        pk.Content[2],
			Boots:           pk.Content[3],
		}}
	}
	return nil
}
