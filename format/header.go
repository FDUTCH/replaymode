package format

import "github.com/google/uuid"

type Header struct {
	Protocol int32     `nbt:"protocol"`
	ShieldID int32     `nbt:"shieldId"`
	Version  string    `nbt:"version"`
	Uuid     uuid.UUID `nbt:"uuid"`
}
