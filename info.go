package main

import (
	"fmt"
	"log"
	"strings"
)

type commandStruct struct {
	Command, Content string
}

var commands map[string]string

func displayInfo(m SimpleMessage) {
	if msg, ok := commands[m.Content]; ok {
		bot.Message(msg)
	}
}

func addInfo(m SimpleMessage) {
	if !hasPerm(m.User, "info_add") {
		log.Println("User", m.User, "doesn't have permission to add info commands!")
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
		bot.Messagef("Command %s deleted", command)
		return
	}

	commands[command] = content
	_, err := db.Exec(fmt.Sprintf("insert into commands(command, content) values ('%s', '%s')", command, content))
	bot.RegisterCallback(command, displayInfo)
	bot.Messagef("Command %s added/changed", command)
	log.Printf("Command %s added. Content: %s", command, commands[command])
	if err != nil {
		log.Println("Permanently saving command failed: ", err)
	}
}

func initInfo() {
	var com []commandStruct
	db.Select(&com, "SELECT * from commands")
	commands = make(map[string]string)
	for _, c := range com {
		commands[c.Command] = c.Content
	}

	for c := range commands {
		bot.RegisterCallback(c, displayInfo)
		log.Println("Command ", c, " registered")
	}
	bot.RegisterCallback("!com", addInfo)
	log.Println("Module INFO initialized")
}
