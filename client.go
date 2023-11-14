// client.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"net/rpc"

	"github.com/mohammedfekri/chatroom/commons"
)

func main() {
	client, err := rpc.Dial("tcp", commons.GetServerAddress())
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	var reply struct{}
	err = client.Call("ChatService.RegisterClient", struct{}{}, &reply)
	if err != nil {
		fmt.Println("Error registering client:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter message (type 'exit' to quit): ")
		scanner.Scan()
		message := scanner.Text()

		if strings.ToLower(message) == "exit" {
			break
		}

		err := client.Call("ChatService.SendMessage", &commons.MessageArgs{Message: message}, &reply)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}

	// Fetch all messages history
	var messages []string
	err = client.Call("ChatService.GetMessages", struct{}{}, &messages)
	if err != nil {
		fmt.Println("Error fetching messages:", err)
	} else {
		fmt.Println("Message History:", messages)
	}
}
