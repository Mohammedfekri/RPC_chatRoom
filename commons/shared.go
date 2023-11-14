
package commons

// MessageArgs represents the arguments for sending a message
type MessageArgs struct {
	Message string
}

// ChatService represents the methods available for RPC
type ChatService interface {
	SendMessage(args *MessageArgs, reply *struct{}) error
	GetMessages(_ struct{}, reply *[]string) error
	RegisterClient(_ struct{}, reply *struct{}) error
	ReceiveMessage(args *MessageArgs, reply *struct{}) error
}

// GetServerAddress returns the fixed server address for the RPC server
func GetServerAddress() string {
	return "0.0.0.0:9999"
}
