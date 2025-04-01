package test

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"replaymode/format"
	"testing"
)

func Test_marshal(t *testing.T) {
	name := "some"
	w := format.NewWriter(minecraft.GameData{}, login.IdentityData{Identity: uuid.New().String()}, name, nil)
	for _, p := range minecraft.DefaultProtocol.Packets(false) {
		pk := p()
		switch pk.ID() {
		case packet.IDInventoryTransaction, packet.IDEvent:
		default:
			err := w.WritePacket(pk)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	err := w.Flush()
	if err != nil {
		t.Fatal(err)
	}
	w.Close()
	r := format.NewReader(name)
	for {
		_, err := r.ReadPacket()
		if err != nil {
			t.Fatal(err)
		}

	}
}
