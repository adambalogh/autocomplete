package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/adambalogh/autocomp/autocomp"
)

func main() {
	a := autocomp.New()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		command := scanner.Text()
		fmt.Println(a.Predict(command, 5))
	}
}
