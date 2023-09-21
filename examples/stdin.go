package main

import (
	"bufio"
	"fmt"
	"loggingdrain/pkg/config"
	templateminer "loggingdrain/pkg/template_miner"
	"os"
)

func main() {
	c := config.LoadConfig("example_config.yaml")
	miner, err := templateminer.NewTemplateMinerWithConfig(c)
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
			break
		default:
			resp := miner.AddLogMessage(input)
			fmt.Printf("\nTemplate: %s\n", resp.TemplateMined)
		}
	}
	fmt.Println(miner.Status())
}
