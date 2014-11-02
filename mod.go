package main

import (
	"log"
	"strconv"
	"strings"
)

func modHandler(m SimpleMessage) {
	str := strings.Fields(m.Content)
	switch str[0] {
	case "!timeout":
		if hasPerm(m.User, "timeout") {
			user := str[1]
			seconds, err := strconv.ParseUint(str[2], 10, 32)
			if err != nil {
				log.Printf("%s is not a valid number of seconds", str[2])
				return
			}
			bot.Messagef("/timeout %s %d", user, seconds)
			log.Printf("Timing out user %s for %d seconds", user, str[2])
		}
	case "!purge":
		if hasPerm(m.User, "purge") {
			user := str[1]
			bot.Messagef("/timeout %s 1", user)
			log.Printf("Timing out user %s for %d seconds", user, str[2])
		}
	}
}

func initMod() {
	bot.RegisterCallback("!timeout", modHandler)
	log.Println("Module MOD initialized")
}
