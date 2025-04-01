package format

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"io"
	"replaymode/translator"
	"sync"
	"time"
)

type ByteWriter interface {
	io.ByteWriter
	io.WriteCloser
	Flush() error
}

type Writer struct {
	proto      minecraft.Protocol
	w          ByteWriter
	timeWriter func() error
	header     Header
	id         []byte
	t          *translator.Translator
	mu         sync.Mutex
	close      chan struct{}
}

func NewWriter(data minecraft.GameData, identityData login.IdentityData, file string, proto minecraft.Protocol) *Writer {
	w, err := Create(file)
	if err != nil {
		panic(err)
	}

	if proto == nil {
		proto = minecraft.DefaultProtocol
	}

	uuid, err := uuid.Parse(identityData.Identity)

	if err != nil {
		panic(err)
	}

	header := Header{
		Protocol: proto.ID(),
		Version:  proto.Ver(),
		Uuid:     uuid,
	}

	for _, it := range data.Items {
		if it.Name == "minecraft:shield" {
			header.ShieldID = int32(it.RuntimeID)
			break
		}
	}

	err = nbt.NewEncoderWithEncoding(w, nbt.LittleEndian).Encode(header)
	if err != nil {
		panic(err)
	}

	writer := &Writer{proto: proto, w: w, header: header, timeWriter: NewTimeWriter(w), id: make([]byte, 4), t: translator.NewTranslator(data.EntityRuntimeID, uuid), close: make(chan struct{})}
	go writer.flusher()
	return writer
}

func (writer *Writer) flusher() {
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-writer.close:
		case <-t.C:
			writer.mu.Lock()
			err := writer.w.Flush()
			writer.mu.Unlock()
			if err != nil {
				_ = writer.Close()
			}
		}
	}
}

func (writer *Writer) WritePacket(pk packet.Packet) error {
	var err error
	writer.mu.Lock()
	defer writer.mu.Unlock()
	for _, p := range writer.t.Translate(pk) {
		err = writer.writePacket(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *Writer) WriteGameData(data minecraft.GameData) error {
	return writer.WritePacket(writer.t.ParseGameData(data))
}

func (writer *Writer) Translator() *translator.Translator {
	return writer.t
}

func (writer *Writer) writePacket(pk packet.Packet) error {
	err := writer.timeWriter()
	if err != nil {
		fmt.Println(1)
		return err
	}
	err = writer.writeId(pk.ID())
	if err != nil {
		fmt.Println(2)
		return err
	}

	pk.Marshal(writer.proto.NewWriter(writer.w, writer.header.ShieldID))
	return nil
}

func (writer *Writer) Close() error {
	close(writer.close)
	return writer.w.Close()
}

func (writer *Writer) writeId(id uint32) error {
	binary.LittleEndian.PutUint32(writer.id, id)
	_, err := writer.w.Write(writer.id)
	return err
}
