package replay

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"replaymode/format"
	"sync"
)

type Recorder struct {
	conn, serverConn *minecraft.Conn
	w                *format.Writer
}

func NewRecorder(conn *minecraft.Conn, serverConn *minecraft.Conn) *Recorder {
	return &Recorder{conn: conn, serverConn: serverConn}
}

func (r *Recorder) Record(file string) {
	r.handle(file)
}
func (r *Recorder) writer(file string, data minecraft.GameData, identityData login.IdentityData, proto minecraft.Protocol) {
	r.w = format.NewWriter(data, identityData, file, proto)
	err := r.w.WriteGameData(data)
	if err != nil {
		panic(err)
	}
}

func (r *Recorder) writePacket(pk packet.Packet) {
	for _, p := range r.w.Translator().Translate(pk) {
		err := r.w.WritePacket(p)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Recorder) connTask(wg *sync.WaitGroup) {
	defer r.conn.Close()
	defer wg.Done()
	for {
		pk, err := r.conn.ReadPacket()
		if err != nil {
			return
		}
		err = r.serverConn.WritePacket(pk)
		if err != nil {
			return
		}
		r.writePacket(pk)
	}
}

func (r *Recorder) serverTask(wg *sync.WaitGroup) {
	defer r.serverConn.Close()
	defer wg.Done()
	for {
		pk, err := r.serverConn.ReadPacket()
		if err != nil {
			return
		}
		err = r.conn.WritePacket(pk)
		if err != nil {
			return
		}
		r.writePacket(pk)
	}
}

func (r *Recorder) handle(file string) {
	conn := r.conn
	serverConn := r.serverConn
	var g sync.WaitGroup
	g.Add(2)

	go func() {
		data := serverConn.GameData()
		r.writer(file, data, conn.IdentityData(), conn.Proto())
		if err := conn.StartGame(data); err != nil {
			panic(err)
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()
	g.Wait()

	var wg sync.WaitGroup
	wg.Add(2)

	go r.serverTask(&wg)
	go r.connTask(&wg)
	wg.Wait()
}

func (r *Recorder) Close() error {
	return r.w.Close()
}
