// server.go
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/yourusername/chatroom/commons"
)

// ChatRoom is the struct that holds the chat room state
type ChatRoom struct {
	messages    []string
	clientConns []*rpc.Client
	mu          sync.RWMutex
}

// ChatService represents the methods available for RPC
type ChatService struct {
	chatRoom *ChatRoom
}

// NewChatRoom creates a new ChatRoom instance
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		messages:    make([]string, 0),
		clientConns: make([]*rpc.Client, 0),
	}
}

// SendMessage is an RPC method for sending messages to the chat room
func (c *ChatService) SendMessage(args *commons.MessageArgs, reply *struct{}) error {
	c.chatRoom.mu.Lock()
	defer c.chatRoom.mu.Unlock()

	// Add the message to the chat room
	c.chatRoom.messages = append(c.chatRoom.messages, args.Message)

	// Broadcast the message to all connected clients
	for _, clientConn := range c.chatRoom.clientConns {
		go func(clientConn *rpc.Client) {
			var reply struct{}
			err := clientConn.Call("ChatService.ReceiveMessage", args, &reply)
			if err != nil {
				fmt.Println("Error broadcasting message to client:", err)
			}
		}(clientConn)
	}

	return nil
}

// GetMessages is an RPC method for retrieving messages from the chat room
func (c *ChatService) GetMessages(_ struct{}, reply *[]string) error {
	c.chatRoom.mu.RLock()
	defer c.chatRoom.mu.RUnlock()
	*reply = append([]string{}, c.chatRoom.messages...)
	return nil
}

// RegisterClient is an RPC method for registering a new client
func (c *ChatService) RegisterClient(_ struct{}, reply *struct{}) error {
	clientConn, _ := rpc.Dial("tcp", commons.GetServerAddress())
	c.chatRoom.mu.Lock()
	defer c.chatRoom.mu.Unlock()
	c.chatRoom.clientConns = append(c.chatRoom.clientConns, clientConn)
	return nil
}

// ReceiveMessage is an RPC method for receiving messages by the client
func (c *ChatService) ReceiveMessage(args *commons.MessageArgs, reply *struct{}) error {
	// You can implement client-specific logic for receiving messages here
	fmt.Printf("Received message: %s\n", args.Message)
	return nil
}

func main() {
	chatRoom := NewChatRoom()
	chatService := &ChatService{chatRoom}

	rpc.Register(chatService)

	// Pooling-Based Implementation
	for {
		conn, err := net.Listen("tcp", commons.GetServerAddress())
		if err != nil {
			fmt.Println("Error starting server:", err)
			return
		}

		clientConn, _ := conn.Accept()
		go func() {
			// Register the client
			var reply struct{}
			err := clientConn.(*net.TCPConn).CloseWrite()
			if err != nil {
				fmt.Println("Error closing write side of TCP connection:", err)
			}
		}()
	}
}
