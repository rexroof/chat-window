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
		if len(text) != 0 {
			c.Say(twChan, text)
		}
	}
	c.Disconnect()
}

func main() {
	twUser := os.Getenv("TWITCH_USERNAME")
	twToken := os.Getenv("TWITCH_TOKEN")
	twChannel := os.Getenv("TWITCH_CHANNEL")
	if len(twChannel) == 0 {
		twChannel = twUser
	}

	paneTitle := os.Getenv("PANE_TITLE")
	if len(paneTitle) > 0 {
		fmt.Printf("\033]2;%s\033\\", paneTitle)
	}

	client := twitch.NewClient(twUser, twToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		colorCode := fmt.Sprintf("\033[38;2;%sm", hexRGB(message.User.Color, message.User.Name))
		bgColorCode := ""
		resetColor := "\033[0m"
		if strings.Contains(message.Tags["msg-id"], "highlighted-message") {
			bgColorCode = fmt.Sprintf("\033[48;2;%sm", hexRGB("#755ebc", "highlight-color-rexroof"))
		}
		fmt.Printf("[%s%s%s]: %s%s%s\n", colorCode, message.User.Name, resetColor, bgColorCode, message.Message, resetColor)
	})

	client.Join(twChannel)
	go sendStdin(client, twChannel)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("[exiting]")
}
