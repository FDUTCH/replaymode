package format

import (
	"encoding/binary"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"io"
	"math"
	"time"
)

type ByteReader interface {
	io.ByteReader
	io.ReadCloser
}

type Reader struct {
	proto      minecraft.Protocol
	r          ByteReader
	timeReader func() (time.Duration, error)
	timeOffset time.Duration
	header     Header
	id         []byte
}

func (reader *Reader) Header() Header {
	return reader.header
}

func NewReader(file string, proto ...minecraft.Protocol) *Reader {

	r, err := Open(file)
	if err != nil {
		panic(err)
	}

	header := Header{}

	err = nbt.NewDecoderWithEncoding(r, nbt.LittleEndian).Decode(&header)
	if err != nil {
		panic(err)
	}

	var p minecraft.Protocol = minecraft.DefaultProtocol

	if len(proto) == 0 {
		proto = append(proto, minecraft.DefaultProtocol)
	}

	for _, ver := range proto {
		if header.Protocol == ver.ID() {
			p = ver
			break
		}
	}

	return &Reader{proto: p, r: r, header: header, timeReader: NewTimeReader(r), timeOffset: time.Duration(math.MaxInt64), id: make([]byte, 4)}
}

func (reader *Reader) Close() error {
	return reader.r.Close()
}

func (reader *Reader) ReadPacket() (pk packet.Packet, err error) {

	t, err := reader.timeReader()

	if err != nil {
		return nil, err
	}

	id, err := reader.readId()
	if err != nil {
		return nil, err
	}

	fn, ok := reader.proto.Packets(false)[id]

	if !ok {
		return nil, fmt.Errorf("unknown packetId (%d)", id)
	}

	pk = fn()

	pk.Marshal(reader.proto.NewReader(reader.r, reader.header.ShieldID, false))

	time.Sleep(t - reader.timeOffset)
	reader.timeOffset = t
	return
}

func (reader *Reader) readId() (uint32, error) {
	_, err := reader.r.Read(reader.id)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(reader.id), nil
}
