package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
)

var globalColors = map[string]string{}

// returns rgb as "red;green;blue"
func hexRGB(s string, u string) string {
	red, green, blue := 0, 0, 0
	fmt.Sscanf(s, "#%02x%02x%02x", &red, &green, &blue)
	if red+green+blue == 0 {
		if val, ok := globalColors[u]; ok {
			return val
		} else {
			rand.Seed(time.Now().UnixNano())
			red = rand.Intn(205) + 50
			green = rand.Intn(205) + 50
			blue = rand.Intn(205) + 50
			globalColors[u] = fmt.Sprintf("%d;%d;%d", red, green, blue)
			return globalColors[u]
		}
	}
	return fmt.Sprintf("%d;%d;%d", red, green, blue)
}

// reads lines from stdin and tries to send them as messages to chat.
func sendStdin(c *twitch.Client, twChan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		if text == "exit" {
			fmt.Printf("quitting!\n")
			break
		}
		fmt.Printf("SENDING: %s", text)
		c.Say(twChan, text)
	}
	c.Disconnect()
}

func main() {
	twUser := os.Getenv("TWITCH_USERNAME")
	twToken := os.Getenv("TWITCH_TOKEN")
	client := twitch.NewClient(twUser, twToken)

	/*
				#D2691E
			  sys.stdout.write(u"\u001b[38;5;" + code + "m " + code.ljust(4))
			  print u"\u001b[0m"
		    colorReset := "\033[0m"
	*/

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		colorCode := fmt.Sprintf("\033[38;2;%sm", hexRGB(message.User.Color, message.User.Name))
		fmt.Printf("[%s%s%s]: %s\n", colorCode, message.User.Name, "\033[0m", message.Message)
	})

	client.Join(twUser)
	// client.Join("middleditch")

	go sendStdin(client, twUser)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("heeeeeey")
}
