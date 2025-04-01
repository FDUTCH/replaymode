package translator

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"time"
)

type Translator struct {
	rid       uint64
	uuid      uuid.UUID
	tolerance time.Time
}

func NewTranslator(rid uint64, uuid uuid.UUID) *Translator {
	return &Translator{rid: rid, uuid: uuid, tolerance: time.Now().Add(time.Second * 2)}
}

func (t *Translator) Translate(p packet.Packet) []packet.Packet {
	switch pk := p.(type) {
	case *packet.PlayerAuthInput:
		return []packet.Packet{TranslatePlayerAuthInput(pk, t.rid)}
	case *packet.Text:
		if pk.TextType != packet.TextTypeChat {
			return []packet.Packet{p}
		}
	case *packet.PlayerSkin:
		return []packet.Packet{t.TranslatePlayerSkin(pk)}
	case *packet.PlayerList:
		return []packet.Packet{t.TranslatePlayerList(pk)}
	case *packet.InventoryContent:
		return t.TranslateInventoryContent(pk)
	case *packet.MovePlayer:
		return t.TranslateMovePlayer(pk)
	case *packet.SetActorData:

	case *packet.ContainerSetData, *packet.BlockPickRequest, *packet.BookEdit, *packet.ClientCacheBlobStatus,
		*packet.CommandRequest, *packet.ContainerClose, *packet.ContainerOpen, *packet.Interact,
		*packet.ItemStackRequest, *packet.ItemStackResponse, *packet.LecternUpdate, *packet.ModalFormResponse,
		*packet.ModalFormRequest, *packet.ClientBoundCloseForm, *packet.NPCRequest, *packet.NPCDialogue,
		*packet.RequestAbility, *packet.RequestChunkRadius, *packet.SubChunkRequest, *packet.ServerBoundLoadingScreen,
		*packet.ServerBoundDiagnostics, *packet.AvailableActorIdentifiers, *packet.BiomeDefinitionList,
		*packet.UpdatePlayerGameType, *packet.InventorySlot:
	default:
		return []packet.Packet{p}
	}
	return nil
}
