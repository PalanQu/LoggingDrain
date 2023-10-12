package main

import (
	"bufio"
	"fmt"
	"os"

	loggingdrain3 "github.com/palanqu/loggingdrain3"
)

func main() {
	miner, err := loggingdrain3.NewTemplateMiner()
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
