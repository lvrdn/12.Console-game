package main

import (
	"fmt"
	"strings"
	"sync"
)

func addPlayer(player *Player) {
	world.MuPl.Lock()
	world.Players = append(world.Players, player)
	world.MuPl.Unlock()
}

func initGame() {
	world = &World{
		Players:   make([]*Player, 0),
		Locations: make([]*Location, 0),
		MuLoc:     &sync.RWMutex{},
		MuPl:      &sync.RWMutex{},
	}

	world.MuLoc.Lock()
	defer world.MuLoc.Unlock()

	world.Locations = append(world.Locations,
		&Location{
			Name:              "улица",
			About:             "на улице весна",
			NextLocationNames: "домой",
			LocationCondition: true,
			Mu:                &sync.RWMutex{},
		},
		&Location{
			Name:  "кухня",
			About: "кухня, ничего интересного",
			LookAround: func(p *Player, l *Location) string {
				answer := "ты находишься на кухне, "

				for _, furniture := range l.Furnitures {
					items := strings.Split(furniture.Items, ",")
					if len(items) == 0 {
						continue
					}
					answer += fmt.Sprintf("на %sе", furniture.Name)
					if len(items) > 1 {
						answer += ":"
					}
					answer += " "

					for _, item := range items {
						answer += fmt.Sprintf("%s, ", item)
					}
				}

				if strings.Contains(p.Items, "рюкзак") {
					answer += "надо идти в универ"
				} else {
					answer += "надо собрать рюкзак и идти в универ"
				}
				return answer
			},
			NextLocationNames: "коридор",
			Furnitures: []*Furniture{
				{
					Name:  "стол",
					Items: "чай",
					Mu:    &sync.RWMutex{},
				},
			},
			Mu: &sync.RWMutex{},
		},
		&Location{
			Name:              "коридор",
			About:             "ничего интересного",
			NextLocationNames: "кухня, комната, улица",
			Furnitures: []*Furniture{
				{
					Name: "дверь",
					Action: func() string {
						world.MuLoc.RLock()
						defer world.MuLoc.RUnlock()

						for _, location := range world.Locations {

							if location.Name == "улица" {
								location.Mu.Lock()
								location.LocationCondition = false
								location.Mu.Unlock()
								return "дверь открыта"
							}
						}
						return "дверь закрыта"
					},
					TriggerObjectName: "ключи",
					Mu:                &sync.RWMutex{},
				},
			},
			LookAround: func(p *Player, l *Location) string {
				return "дверь на улицу или другие помещения"
			},
			Mu: &sync.RWMutex{},
		},
		&Location{
			Name:              "комната",
			About:             "ты в своей комнате",
			NextLocationNames: "коридор",
			Furnitures: []*Furniture{
				{
					Name:  "стол",
					Items: "ключи,конспекты",
					Mu:    &sync.RWMutex{},
				},
				{
					Name:  "стул",
					Items: "рюкзак",
					Mu:    &sync.RWMutex{},
				},
			},
			LookAround: func(p *Player, l *Location) string {
				var answer string
				for _, furniture := range l.Furnitures {
					if len(furniture.Items) == 0 {
						continue
					}
					answer += fmt.Sprintf("на %sе", furniture.Name)
					if furniture.Name == "стол" {
						answer += ": "
					} else {
						answer += " - "
					}
					items := strings.Split(furniture.Items, ",")
					for _, item := range items {
						answer += fmt.Sprintf("%s, ", item)
					}

				}
				answer = strings.TrimSuffix(answer, ", ")
				if answer == "" {
					answer = "пустая комната"
				}
				return answer
			},
			Mu: &sync.RWMutex{},
		},
	)
}
