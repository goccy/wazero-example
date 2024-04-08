package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		content, err := reader.ReadString('\n')
		fmt.Fprintf(os.Stderr, "content = %s. err = %v\n", content, err)
	}
}
