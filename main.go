package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/adambalogh/autocomp/autocomp"
)

func main() {
	big, err := os.Open("big.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	a := autocomp.New(big)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		command := scanner.Text()
		fmt.Println(a.Predict(command, 5))
	}
}
