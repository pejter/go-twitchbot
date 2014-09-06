package main

import (
	"log"
	"strings"
)

var commands = map[string]string{
	"!cool": "Ice cold!",
	"!info": "I have information if you have coin",
}

func isMod(bot *IRCBot, user string) bool {
	for _, u := range bot.Moderators {
		if u == user {
			return true
		}
	}
	return false
}

func displayInfo(bot *IRCBot, m SimpleMessage) {
	if msg, ok := commands[m.Content]; ok {
		bot.Message(msg)
	}
}

func addInfo(bot *IRCBot, m SimpleMessage) {
	if !isMod(bot, m.User) {
		log.Println("User ", m.User, " isn't a moderator")
		return
	}

	str := strings.Fields(m.Content)
	command := str[1]
	content := strings.Join(str[2:], " ")
	if !strings.HasPrefix(command, "!") {
		command = "!" + command
	}

	if content == "delete" {
		bot.RemoveCallback(command)
	}

	commands[command] = content
	bot.RegisterCallback(command, displayInfo)
	bot.Messagef("Command %s added/changed", command)
	log.Printf("Command %s added. Content: %s", command, commands[str[1]])
}

func initInfo(bot *IRCBot) {
	for c := range commands {
		bot.RegisterCallback(c, displayInfo)
		log.Println("Command ", c, " registered")
	}
	bot.RegisterCallback("!command", addInfo)
	log.Println("module INFO initialized")
}
