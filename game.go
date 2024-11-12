package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {

	lastOutput := map[string]string{
		"Tristan": "",
	}

	players := map[string]*Player{
		"Tristan": NewPlayer("Tristan"),
	}

	mu := &sync.Mutex{}

	go func() {
		output := players["Tristan"].GetOutput()
		for msg := range output {
			mu.Lock()
			lastOutput["Tristan"] = msg
			mu.Unlock()
		}
	}()

	initGame()
	addPlayer(players["Tristan"])

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		command := sc.Text()
		players["Tristan"].HandleInput(command)
		time.Sleep(time.Millisecond)
		runtime.Gosched() // дадим считать ответ
		mu.Lock()
		answer := lastOutput["Tristan"]
		mu.Unlock()
		fmt.Println(answer)
	}

}
