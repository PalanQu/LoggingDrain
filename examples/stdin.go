package main

import (
	"bufio"
	"fmt"
	"os"

	loggingdrain "github.com/palanqu/loggingdrain"
)

func main() {
	miner, err := loggingdrain.NewTemplateMiner()
	if err != nil {
		panic(err)
	}
	fmt.Println("input q to quit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "q":
			fmt.Println("quit")
			fmt.Println(miner.Status())
			return
		default:
			resp := miner.AddLogMessage(input)
			fmt.Printf("\nTemplate: %s\n", resp.TemplateMined)
		}
	}
}
