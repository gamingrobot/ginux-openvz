package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/sessions"
	"io"
	"log"
	"net/http"
	"os"
	//"io/ioutil"
	"bufio"
	"os/exec"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

const MAX_VMS = 20

var currentVms int64 = 0

type Config struct {
	Secret string
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	decoder.Decode(&config)

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte(config.Secret))
	m.Use(sessions.Sessions("session", store))

	m.Get("/set", func(session sessions.Session) string {
		session.Set("hello", "world")
		return "OK"
	})

	m.Get("/get", func(session sessions.Session) string {
		v := session.Get("hello")
		if v == nil {
			return ""
		}
		return v.(string)
	})

	m.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		if currentVms >= MAX_VMS {
			http.Error(w, "Over Capacity", 400)
			return
		}
		//spawn console
		cmd := exec.Command("vzctl", "console", "100")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}
		go readLoop(stdout, ws)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Println(err)
			return
		}
		cmd.Start()
		atomic.AddInt64(&currentVms, 1)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				//kill console
				cmd.Process.Kill()
				cmd.Process.Wait()
				atomic.AddInt64(&currentVms, -1)
				//log.Println(err)
				return
			} else {
				//send msg to console
				stdin.Write(message)
			}
		}
	})
	address := "0.0.0.0:3000"
	log.Println("Martini started on", address)
	log.Fatal(http.ListenAndServe(address, m))
}

func readLoop(output io.Reader, ws *websocket.Conn) {
	reader := bufio.NewReader(output)
	buffer := []byte{}
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			return
		}
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		buffer = append(buffer, b)
		valid := utf8.Valid(buffer)
		if valid {
			ws.WriteMessage(websocket.TextMessage, buffer)
			buffer = []byte{}
		}
	}
}
