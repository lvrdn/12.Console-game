package main

import (
	"fmt"
	"strings"
	"sync"
)

type Player struct {
	Name                string
	CurrentLocationName string
	TakeItems           bool
	Items               string
	Mu                  *sync.RWMutex
	Out                 chan string
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:                name,
		CurrentLocationName: "кухня",
		Out:                 make(chan string),
		Mu:                  &sync.RWMutex{},
		Items:               "",
		TakeItems:           false,
	}
}

func (p *Player) GetOutput() chan string {
	return p.Out
}

func (p *Player) LookAround() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	currentLocation := world.GetLocation(p.CurrentLocationName)

	currentLocation.Mu.RLock()
	defer currentLocation.Mu.RUnlock()

	answer := currentLocation.LookAround(p, currentLocation)

	answer += fmt.Sprintf(". можно пройти - %s", currentLocation.NextLocationNames)

	world.MuPl.RLock()
	fl := true
	for _, player := range world.Players {
		if player.Name != p.Name && player.CurrentLocationName == p.CurrentLocationName {
			if fl {
				answer += ". Кроме вас тут ещё "
				fl = false
			}
			answer += player.Name + ", "
		}
	}
	world.MuPl.RUnlock()
	answer = strings.TrimSuffix(answer, ", ")
	p.Out <- answer

}

func (p *Player) GoNextLocation(NextLocationName string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	currentLocation := world.GetLocation(p.CurrentLocationName)

	if !strings.Contains(currentLocation.NextLocationNames, NextLocationName) {
		answer := fmt.Sprintf("нет пути в %s", NextLocationName)
		p.Out <- answer
		return
	}

	nextLocation := world.GetLocation(NextLocationName)
	nextLocation.Mu.RLock()
	defer nextLocation.Mu.RUnlock()

	if nextLocation.LocationCondition {
		p.Out <- "дверь закрыта"
		return
	}

	p.CurrentLocationName = nextLocation.Name

	answer := fmt.Sprintf("%s. можно пройти - %s", nextLocation.About, nextLocation.NextLocationNames)

	p.Out <- answer

}

func (p *Player) PutOn(item string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	currentLocation := world.GetLocation(p.CurrentLocationName)

	currentLocation.Mu.RLock()
	defer currentLocation.Mu.RUnlock()

	for _, furniture := range currentLocation.Furnitures {
		furniture.Mu.RLock()
		defer furniture.Mu.RUnlock()

		if strings.Contains(furniture.Items, item) {

			if item == "рюкзак" {
				p.TakeItems = true
			}
			if p.Items == "" {
				p.Items += item
			} else {
				p.Items += "," + item
			}

			switch {
			case len(furniture.Items) == len(item):
				furniture.Items = ""
			case strings.Index(furniture.Items, item) == 0:
				furniture.Items = strings.TrimPrefix(furniture.Items, item+",")
			case strings.Index(furniture.Items, item)+len(item) == len(furniture.Items):
				furniture.Items = strings.TrimSuffix(furniture.Items, ","+item)
			default:
				furniture.Items = strings.Replace(furniture.Items, item+",", "", 1)
			}

			answer := "вы одели: " + item
			p.Out <- answer
			return
		}
	}

	p.Out <- "нет такого"
}

func (p *Player) Take(item string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if !p.TakeItems {
		p.Out <- "некуда класть"
		return
	}

	currentLocation := world.GetLocation(p.CurrentLocationName)

	currentLocation.Mu.RLock()
	defer currentLocation.Mu.RUnlock()

	for _, furniture := range currentLocation.Furnitures {
		furniture.Mu.Lock()
		defer furniture.Mu.Unlock()

		if strings.Contains(furniture.Items, item) {

			if p.Items == "" {
				p.Items += item
			} else {
				p.Items += "," + item
			}

			switch {
			case len(furniture.Items) == len(item):
				furniture.Items = ""
			case strings.Index(furniture.Items, item) == 0:
				furniture.Items = strings.TrimPrefix(furniture.Items, item+",")
			case strings.Index(furniture.Items, item)+len(item) == len(furniture.Items):
				furniture.Items = strings.TrimSuffix(furniture.Items, ","+item)
			default:
				furniture.Items = strings.Replace(furniture.Items, item+",", "", 1)
			}

			answer := "предмет добавлен в инвентарь: " + item
			p.Out <- answer
			return

		}

	}

	p.Out <- "нет такого"
}

func (p *Player) Use(item, object string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if !strings.Contains(p.Items, item) {
		answer := "нет предмета в инвентаре - " + item
		p.Out <- answer
		return
	}

	currentLocation := world.GetLocation(p.CurrentLocationName)

	findedFurniture := currentLocation.GetFurniture(object)
	if findedFurniture == nil {
		p.Out <- "не к чему применить"
		return
	}

	if findedFurniture.TriggerObjectName != item {
		p.Out <- "не к чему применить"
		return
	}

	answer := findedFurniture.Action()
	p.Out <- answer

}

func (p *Player) Say(phrase string) {

	world.MuPl.RLock()

	for _, player := range world.Players {
		if player.CurrentLocationName == p.CurrentLocationName {
			answer := fmt.Sprintf("%s говорит: %s", p.Name, phrase)

			player.Mu.Lock()
			player.Out <- answer
			player.Mu.Unlock()
		}
	}
	world.MuPl.RUnlock()
}

func (p *Player) SayToPlayer(phrase string, playerName string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	player := world.GetPlayer(playerName)

	player.Mu.Lock()
	defer player.Mu.Unlock()

	if p.CurrentLocationName != player.CurrentLocationName {
		p.Out <- "тут нет такого игрока"
		return
	}

	var answer string
	if phrase == "" {
		answer = fmt.Sprintf("%s выразительно молчит, смотря на вас", p.Name)
	} else {
		answer = fmt.Sprintf("%s говорит вам: %s", p.Name, phrase)
	}

	player.Out <- answer
}
