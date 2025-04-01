package main

import (
	"encoding/json"
	"flag"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
	"log"
	"os"
	"replaymode/format"
	"replaymode/replay"
	"time"
)

var (
	output  = flag.String("o", FormatTime(time.Now()), "output file")
	address = flag.String("a", "play.cubecraft.net:19132", "remote address")
	input   = flag.String("i", "replay", "input file")
)

func main() {
	flag.Parse()
	if flagPresent("-i") {
		Play()
	} else {
		Record()
	}
}

func Play() {
	listener, err := minecraft.Listen("raknet", ":19132")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		play(c.(*minecraft.Conn), *input)
	}
}

func play(conn *minecraft.Conn, file string) {
	player := replay.NewPlayer(conn, format.NewReader(file))
	player.Play()
}

func Record() {
	src := tokenSource()
	listener, err := minecraft.Listen("raknet", ":19132")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		err = record(c.(*minecraft.Conn), *address, src, *output)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func record(conn *minecraft.Conn, remote string, src oauth2.TokenSource, file string) error {
	serverConn, err := minecraft.Dialer{
		ClientData:           conn.ClientData(),
		IdentityData:         conn.IdentityData(),
		TokenSource:          src,
		DownloadResourcePack: skipDownload,
	}.Dial("raknet", remote)
	if err != nil {
		return err
	}
	recorder := replay.NewRecorder(conn, serverConn)
	recorder.Record(file)
	return nil
}

func skipDownload(id uuid.UUID, version string, current int, total int) bool {
	return false
}

func tokenSource() oauth2.TokenSource {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	token := new(oauth2.Token)
	tokenData, err := os.ReadFile("token.tok")
	if err == nil {
		_ = json.Unmarshal(tokenData, token)
	} else {
		token, err = auth.RequestLiveToken()
		check(err)
	}
	src := auth.RefreshTokenSource(token)
	_, err = src.Token()
	if err != nil {
		// The cached refresh token expired and can no longer be used to obtain a new token. We require the
		// user to log in again and use that token instead.
		token, err = auth.RequestLiveToken()
		check(err)
		src = auth.RefreshTokenSource(token)
	}
	tok, _ := src.Token()
	b, _ := json.Marshal(tok)
	_ = os.WriteFile("token.tok", b, 0644)

	return src
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15-04-05")
}

func flagPresent(f string) bool {
	for _, val := range os.Args {
		if val == f {
			return true
		}
	}
	return false
}
