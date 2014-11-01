package main

import (
	"fmt"
	"log"
	"strings"
)

type Poll struct {
	TotalVotes   uint16
	Choices      map[string]uint16
	AlreadyVoted map[string]struct{}
	Winners      []string
}

func (p *Poll) Start(choices []string) {
	p.TotalVotes = 0
	p.Choices = make(map[string]uint16, len(choices))
	p.AlreadyVoted = make(map[string]struct{})
	p.Winners = nil
	for _, c := range choices {
		p.Choices[c] = 0
		err := bot.RegisterCallback("!"+c, voteHandler)
		if err != nil {
			for c, _ := range p.Choices {
				bot.RemoveCallback("!" + c)
			}
			return
		}
	}
	bot.Messagef("New poll started with choices: %s", strings.Join(choices, ","))
}

func (p *Poll) Stop() {
	if p.TotalVotes == 0 {
		bot.Message("Nobody voted! Poll closed")
		return
	}
	var max uint16 = 0
	for c, n := range p.Choices {
		bot.RemoveCallback("!" + c)
		if n > max {
			max = n
		}
	}

	for c, n := range p.Choices {
		if n == max {
			p.Winners = append(p.Winners, fmt.Sprintf("%s:%d", c, n))
		}
	}
}

var currentPoll Poll

func voteHandler(m SimpleMessage) {
	if _, ok := currentPoll.AlreadyVoted[m.User]; ok {
		return
	}
	currentPoll.TotalVotes++
	currentPoll.Choices[m.Content]++
	currentPoll.AlreadyVoted[m.User] = struct{}{}
}

func pollHandler(m SimpleMessage) {
	if !hasPerm(m.User, "poll_edit") {
		log.Println("User", m.User, "doesn't have permission to edit polls!")
		return
	}

	m.Content = strings.TrimPrefix(m.Content, "!poll ")
	if strings.HasPrefix(m.Content, "open") {
		list := strings.TrimPrefix(m.Content, "open ")
		choices := strings.Split(list, ",")
		currentPoll.Start(choices)
	} else if strings.HasPrefix(m.Content, "close") {
		if currentPoll.Winners != nil {
			bot.Message("No poll currently open")
			return
		}
		currentPoll.Stop()
		log.Println("Winners: ", strings.Join(currentPoll.Winners, ", "))
		bot.Messagef("Poll ended! Total votes: %d Winners: %s", currentPoll.TotalVotes, strings.Join(currentPoll.Winners, ", "))
	}
}

func initPoll() {
	bot.RegisterCallback("!poll", pollHandler)
	log.Println("Module POLL initialized")
}
