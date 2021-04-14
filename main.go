package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
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
	} else if red+green+blue < 150 {
		// bump up colors if they are too dark.
		//  potentially update this to instead mask background
		// 24 24 27 = 75  [xmetrix]
		// 14 12 19 = 45  [nyxiative]
		factor := int((160 - red - green - blue) / 3)
		red += factor
		green += factor
		blue += factor
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
		bits := message.Bits
		colorCode := fmt.Sprintf("\033[38;2;%sm", hexRGB(message.User.Color, message.User.Name))
		bgColorCode := ""
		resetColor := "\033[0m"
		if bits > 0 {
			bgColorCode = fmt.Sprintf("\033[48;2;%sm", hexRGB("#f36f00", "highlight-color-bitsbits"))
		} else if strings.Contains(message.Tags["msg-id"], "highlighted-message") {
			bgColorCode = fmt.Sprintf("\033[48;2;%sm", hexRGB("#755ebc", "highlight-color-highlight"))
		} else if strings.Contains(message.Message, "rexroof") {
			bgColorCode = fmt.Sprintf("\033[48;2;%sm", hexRGB("#b54624", "highlight-color-rexroof"))
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
