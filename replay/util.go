package replay

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type inventories struct {
	MainHand   protocol.ItemInstance
	Offhand    protocol.ItemInstance
	Helmet     protocol.ItemInstance
	Chestplate protocol.ItemInstance
	Leggings   protocol.ItemInstance
	Boots      protocol.ItemInstance
}

func (i *inventories) send(player *Player) {
	rid := player.gameData.EntityRuntimeID

	//main hand
	player.writePacket(&packet.MobEquipment{
		EntityRuntimeID: rid,
		NewItem:         i.MainHand,
	})

	//offhand
	player.writePacket(&packet.MobEquipment{
		EntityRuntimeID: rid,
		NewItem:         i.Offhand,
		WindowID:        protocol.WindowIDOffHand,
	})

	//armour
	player.writePacket(&packet.MobArmourEquipment{
		EntityRuntimeID: rid,
		Helmet:          i.Helmet,
		Chestplate:      i.Chestplate,
		Leggings:        i.Leggings,
		Boots:           i.Boots,
	})

}

func (i *inventories) updateMobEquipment(pk *packet.MobEquipment) {
	switch pk.WindowID {
	case protocol.WindowIDOffHand:
		i.Offhand = pk.NewItem
	case protocol.WindowIDInventory:
		i.MainHand = pk.NewItem
	}
}

func (i *inventories) updateMobArmourEquipment(pk *packet.MobArmourEquipment) {
	i.Helmet = pk.Helmet
	i.Chestplate = pk.Chestplate
	i.Leggings = pk.Leggings
	i.Boots = pk.Boots
}

type location struct {
	Position mgl32.Vec3
	Pitch    float32
	Yaw      float32
	HeadYaw  float32
}

func (p *location) writeMoveActorAbsolute(pk *packet.MoveActorAbsolute) {
	p.Position = pk.Position
}

func (p *location) writeMovePlayer(pk *packet.MovePlayer) {
	p.Position = pk.Position
	p.Pitch = pk.Pitch
	p.Yaw = pk.Yaw
	p.HeadYaw = pk.HeadYaw
}
