package replay

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"replaymode/format"
	"replaymode/translator"
	"time"
)

type Player struct {
	conn     *minecraft.Conn
	reader   *format.Reader
	gameData minecraft.GameData
}

func NewPlayer(conn *minecraft.Conn, reader *format.Reader) *Player {
	return &Player{conn: conn, reader: reader}
}

func (player *Player) Play() {

	defer player.reader.Close()

	for {
		p, ok := player.readPacket()
		if !ok {
			return
		}

		switch pk := p.(type) {

		case *packet.StartGame:
			player.gameData = translator.GameDataFromPacket(pk)
		case *packet.ItemRegistry:
			player.gameData.Items = pk.Items
			player.start()
		case *packet.MovePlayer:
			if pk.EntityRuntimeID == player.gameData.EntityRuntimeID {
				time.AfterFunc(time.Second/20, func() {
					player.writePacket(pk)
				})
			} else {
				player.writePacket(pk)
			}
		default:
			if !player.writePacket(pk) {
				return
			}
		}

	}
}

func (player *Player) start() {

	viewerData := translator.GameDataForViewer(player.gameData)
	err := player.conn.StartGame(viewerData)
	if err != nil {
		player.disconnectErr(err)
		return
	}

	player.sendAbilities()
	addPacket := translator.TranslateGameData(player.gameData, player.conn.IdentityData().DisplayName)
	if !player.writePacket(addPacket) {
		return
	}
}

//func (player *Player) resend(addPacket *packet.AddPlayer) {
//	player.writePacket(&packet.RemoveActor{EntityUniqueID: player.gameData.EntityUniqueID})
//	addPacket.Pitch = player.Pitch
//	addPacket.Yaw = player.Yaw
//	addPacket.HeadYaw = player.HeadYaw
//	addPacket.Position = player.Position
//	player.writePacket(addPacket)
//	player.writePacket(&packet.MovePlayer{
//		EntityRuntimeID: addPacket.EntityRuntimeID,
//		Position:        player.Position,
//		Pitch:           player.Pitch,
//		Yaw:             player.Yaw,
//		HeadYaw:         player.HeadYaw,
//		Mode:            packet.MoveModeTeleport,
//		OnGround:        true,
//	})
//	player.send(player)
//}

func (player *Player) writePacket(pk packet.Packet) bool {
	return player.conn.WritePacket(pk) == nil
}

func (player *Player) readPacket() (packet.Packet, bool) {
	pk, err := player.reader.ReadPacket()
	player.disconnectErr(err)
	return pk, err == nil
}

func (player *Player) disconnectErr(err error) {
	if err != nil {
		player.writePacket(&packet.Disconnect{
			Message: err.Error(),
		})
		_ = player.conn.Close()
	}
}

func (player *Player) log(msg string) {
	player.writePacket(&packet.Text{
		TextType:   packet.TextTypeRaw,
		SourceName: "Replay",
		Message:    msg,
	})
}

func (player *Player) sendAbilities() {

	conn := player.conn
	m := protocol.NewEntityMetadata()
	//m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagHasGravity)
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagBreathing)
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible)
	conn.WritePacket(&packet.SetActorData{
		EntityRuntimeID: translator.Rid,
		EntityMetadata:  m,
	})

	abilities := uint32(0)
	abilities |= protocol.AbilityMayFly
	abilities |= protocol.AbilityNoClip
	//abilities |= protocol.AbilityInvulnerable

	conn.WritePacket(&packet.UpdateAbilities{AbilityData: protocol.AbilityData{
		EntityUniqueID:     translator.Rid,
		PlayerPermissions:  packet.PermissionLevelVisitor,
		CommandPermissions: packet.CommandPermissionLevelNormal,
		Layers: []protocol.AbilityLayer{
			{
				Type:      protocol.AbilityLayerTypeBase,
				Abilities: protocol.AbilityCount - 1,
				Values:    abilities,
				FlySpeed:  float32(0.05),
				WalkSpeed: float32(0.1),
			}, {
				Type:      protocol.AbilityLayerTypeSpectator,
				Abilities: protocol.AbilityCount - 1,
				Values:    protocol.AbilityFlying | protocol.AbilityNoClip,
				FlySpeed:  float32(0.2),
				WalkSpeed: float32(0.1),
			},
		},
	}})
	conn.WritePacket(&packet.SetHud{
		Elements: []int32{packet.HudElementHealth, packet.HudElementHunger, packet.HudElementAirBubbles, packet.HudElementArmour, packet.HudElementProgressBar},
	})
}

//func (player *Player) chunkPos() protocol.ChunkPos {
//	return protocol.ChunkPos{int32(player.Position[0]) >> 4, int32(player.Position[2]) >> 4}
//}
