package main

import "strings"

func (p *Player) HandleInput(command string) {
	commandSplited := strings.Split(command, " ")

	switch commandSplited[0] {
	case "осмотреться":
		p.LookAround()

	case "идти":
		p.GoNextLocation(commandSplited[1])

	case "одеть":
		p.PutOn(commandSplited[1])
	case "взять":
		p.Take(commandSplited[1])
	case "применить":
		if len(commandSplited) != 3 {
			p.Out <- "введите команду следующим образом: применить предмет1 предмет2"
			return
		}
		p.Use(commandSplited[1], commandSplited[2])
	case "сказать":
		phrase := strings.Join(commandSplited[1:], " ")
		p.Say(phrase)

	case "сказать_игроку":
		playerName := commandSplited[1]
		phrase := strings.Join(commandSplited[2:], " ")
		p.SayToPlayer(phrase, playerName)
	default:
		p.Out <- "неизвестная команда, доступные команды: осмотреться, идти, одеть, взять, применить, сказать, сказать_игроку"
	}
}
