package main

import (
	"log"
	"regexp"
)

func checkSpam(m SimpleMessage) {
	if hasPerm(m.User, "%") {
		//return
	}
	linkPattern := regexp.MustCompile(`(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`)
	isLink := linkPattern.MatchString(m.Content)
	if isLink {
		bot.Messagef("/timeout %s 1", m.User)
		bot.Messagef("Purging %s. Reason: posting links not allowed", m.User)
		log.Printf("Purging %s. Reason: posting links not allowed", m.User)
	}
	return
}
