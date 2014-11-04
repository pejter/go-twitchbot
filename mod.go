package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func changeStatus(key, value string) {
	reader := strings.NewReader("channel%5B" + key + "%5D=" + url.QueryEscape(value))
	log.Printf("Created reader: %v", reader)

	req, err := http.NewRequest("PUT", "https://api.twitch.tv/kraken/channels/"+cfg.General.Room, reader)
	if err != nil {
		log.Fatalln("Error occured while creating request: ", err)
	}

	req.Header.Add("Authorization", "OAuth "+cfg.General.Token)
	req.Header.Add("Accept", "application/vnd.twitchtv.v3+json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error occured while creating request: ", err)
	}
	res.Body.Close()
}

func modHandler(m SimpleMessage) {
	str := strings.Fields(m.Content)
	switch str[0] {
	case "!timeout":
		if hasPerm(m.User, "timeout") {
			user := str[1]
			var seconds uint64
			if len(str) == 3 {
				var err error
				seconds, err = strconv.ParseUint(str[2], 10, 32)
				if err != nil {
					log.Printf("%s is not a valid number of seconds", str[2])
					return
				}
			} else {
				seconds = 30
			}
			bot.Messagef("/timeout %s %d", user, seconds)
			log.Printf("Timing out user %s for %d seconds", user, seconds)
		}
	case "!purge":
		if hasPerm(m.User, "purge") {
			user := str[1]
			bot.Messagef("/timeout %s 1", user)
			log.Printf("Timing out user %s for %d seconds", user, str[2])
		}
	}
}

func statusHandler(m SimpleMessage) {
	if !hasPerm(m.User, "status_change") {
		return
	}
	m.Content = strings.TrimPrefix(m.Content, "!")
	status := strings.Fields(m.Content)
	key := status[0]
	value := strings.Join(status[1:], " ")

	if key == "title" {
		key = "status"
	}
	changeStatus(key, value)
	bot.Messagef("Changed current %s to %s", status[0], value)
}

func initMod() {
	bot.RegisterCallback("!timeout", modHandler)
	bot.RegisterCallback("!purge", modHandler)
	bot.RegisterCallback("!title", statusHandler)
	bot.RegisterCallback("!game", statusHandler)
	log.Println("Module MOD initialized")
}
