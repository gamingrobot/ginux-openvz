package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

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

	m.Run()
}
