package main

import "fmt"
import "bufio"
import "os"
import "strings"
import "io"

func main() {
	consoleReader := bufio.NewReader(os.Stdin)
	quit := false
	for !quit {
		fmt.Print("> ")
		text, e := consoleReader.ReadString('\n')
		if e == io.EOF {
			quit = true
		} else if e != nil {
			fmt.Printf("Error: %v", e);
			return
		}
		text = strings.TrimSuffix(text, "\n")
		if text == "quit" || text == "exit" {quit = true}
		fmt.Println(text)
	}
}
