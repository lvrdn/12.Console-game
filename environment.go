package main

import "sync"

var world *World

type World struct {
	Players   []*Player
	Locations []*Location
	MuPl      *sync.RWMutex
	MuLoc     *sync.RWMutex
}

type Location struct {
	Name              string
	About             string
	LookAround        func(p *Player, l *Location) string
	NextLocationNames string
	LocationCondition bool
	Furnitures        []*Furniture
	Mu                *sync.RWMutex
}

type Furniture struct {
	Name              string
	Items             string
	Action            func() string
	TriggerObjectName string
	Mu                *sync.RWMutex
}

func (w *World) GetPlayer(name string) *Player {
	w.MuPl.RLock()
	defer w.MuPl.RUnlock()
	for _, player := range w.Players {
		if player.Name == name {
			return player
		}
	}
	return nil
}

func (w *World) GetLocation(name string) *Location {
	w.MuLoc.RLock()
	defer w.MuLoc.RUnlock()
	for _, location := range w.Locations {
		if location.Name == name {
			return location
		}

	}
	return nil
}

func (l *Location) GetFurniture(name string) *Furniture {
	l.Mu.RLock()
	defer l.Mu.RUnlock()
	for _, furniture := range l.Furnitures {
		if furniture.Name == name {
			return furniture
		}

	}
	return nil
}
