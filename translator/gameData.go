package translator

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func TranslateGameData(game minecraft.GameData, nick string) *packet.AddPlayer {

	return &packet.AddPlayer{
		UUID:            StumpUUID,
		Username:        nick,
		EntityRuntimeID: game.EntityRuntimeID,
		Position:        game.PlayerPosition,
		Pitch:           game.Pitch,
		Yaw:             game.Yaw,
		HeadYaw:         game.Yaw,
		GameType:        game.PlayerGameMode,
		DeviceID:        uuid.New().String(),
		AbilityData: protocol.AbilityData{
			EntityUniqueID: int64(game.EntityRuntimeID),
			Layers: []protocol.AbilityLayer{{
				Type:      protocol.AbilityLayerTypeBase,
				Abilities: protocol.AbilityCount - 1,
			}},
		},
		BuildPlatform: int32(protocol.DeviceAndroid),
	}
}

func (t *Translator) ParseGameData(gameData minecraft.GameData) *packet.StartGame {
	return &packet.StartGame{
		Difficulty:                   gameData.Difficulty,
		WorldName:                    gameData.WorldName,
		WorldSeed:                    gameData.WorldSeed,
		EntityUniqueID:               gameData.EntityUniqueID,
		EntityRuntimeID:              gameData.EntityRuntimeID,
		PlayerGameMode:               gameData.PlayerGameMode,
		BaseGameVersion:              gameData.BaseGameVersion,
		PlayerPosition:               gameData.PlayerPosition,
		Pitch:                        gameData.Pitch,
		Yaw:                          gameData.Yaw,
		Dimension:                    gameData.Dimension,
		WorldSpawn:                   gameData.WorldSpawn,
		EditorWorldType:              gameData.EditorWorldType,
		CreatedInEditor:              gameData.CreatedInEditor,
		ExportedFromEditor:           gameData.ExportedFromEditor,
		PersonaDisabled:              gameData.PersonaDisabled,
		CustomSkinsDisabled:          gameData.CustomSkinsDisabled,
		GameRules:                    gameData.GameRules,
		Time:                         gameData.Time,
		ServerBlockStateChecksum:     gameData.ServerBlockStateChecksum,
		Blocks:                       gameData.CustomBlocks,
		Items:                        gameData.Items,
		PlayerMovementSettings:       gameData.PlayerMovementSettings,
		WorldGameMode:                gameData.WorldGameMode,
		Hardcore:                     gameData.Hardcore,
		ServerAuthoritativeInventory: gameData.ServerAuthoritativeInventory,
		PlayerPermissions:            gameData.PlayerPermissions,
		ChatRestrictionLevel:         gameData.ChatRestrictionLevel,
		DisablePlayerInteractions:    gameData.DisablePlayerInteractions,
		ClientSideGeneration:         gameData.ClientSideGeneration,
		Experiments:                  gameData.Experiments,
		UseBlockNetworkIDHashes:      gameData.UseBlockNetworkIDHashes,
	}
}

func GameDataFromPacket(pk *packet.StartGame) minecraft.GameData {
	return minecraft.GameData{
		Difficulty:                   pk.Difficulty,
		WorldName:                    pk.WorldName,
		WorldSeed:                    pk.WorldSeed,
		EntityUniqueID:               pk.EntityUniqueID,
		EntityRuntimeID:              pk.EntityRuntimeID,
		PlayerGameMode:               pk.PlayerGameMode,
		BaseGameVersion:              pk.BaseGameVersion,
		PlayerPosition:               pk.PlayerPosition,
		Pitch:                        pk.Pitch,
		Yaw:                          pk.Yaw,
		Dimension:                    pk.Dimension,
		WorldSpawn:                   pk.WorldSpawn,
		EditorWorldType:              pk.EditorWorldType,
		CreatedInEditor:              pk.CreatedInEditor,
		ExportedFromEditor:           pk.ExportedFromEditor,
		PersonaDisabled:              pk.PersonaDisabled,
		CustomSkinsDisabled:          pk.CustomSkinsDisabled,
		GameRules:                    pk.GameRules,
		Time:                         pk.Time,
		ServerBlockStateChecksum:     pk.ServerBlockStateChecksum,
		CustomBlocks:                 pk.Blocks,
		Items:                        pk.Items,
		PlayerMovementSettings:       pk.PlayerMovementSettings,
		WorldGameMode:                pk.WorldGameMode,
		Hardcore:                     pk.Hardcore,
		ServerAuthoritativeInventory: pk.ServerAuthoritativeInventory,
		PlayerPermissions:            pk.PlayerPermissions,
		ChatRestrictionLevel:         pk.ChatRestrictionLevel,
		DisablePlayerInteractions:    pk.DisablePlayerInteractions,
		ClientSideGeneration:         pk.ClientSideGeneration,
		Experiments:                  pk.Experiments,
		UseBlockNetworkIDHashes:      pk.UseBlockNetworkIDHashes,
	}
}

func GameDataForViewer(game minecraft.GameData) minecraft.GameData {
	game.PlayerGameMode = 2
	game.WorldName = text.Amethyst + text.Bold + "Replay " + text.Reset + game.WorldName
	game.EntityRuntimeID = Rid
	game.EntityUniqueID = Rid
	return game
}
