package main

import (
	"bufio"
	"fmt"
	"net"
)

func parseCommand(message RESPMessage) string {
	if message.Type == Array {
		array := message.Payload.([]RESPMessage)
		var command string
		for _, el := range array {
			command += parseCommand(el)
			command += " "
		}
	} else if message.Type == SimpleString {
		return "simplestring"
	} else if message.Type == BulkString {
		return "bulkstring"
	} else if message.Type == Integer {
		return "integer"
	}

	return "-Error"
}

func handleMessage(message RESPMessage) RESPMessage {
	switch message.Type {
	case Array:
		return handleArrayMessageType(message.Payload.([]RESPMessage))
	}
	return RESPMessage{
		Type:    36,
		Payload: "PONG",
	}
}

func handleArrayMessageType(messages []RESPMessage) RESPMessage {
	// first message should be command:
	command_message := messages[0]
	fmt.Println(command_message)
	return RESPMessage{Type: BulkString, Payload: handleCommand(command_message.Payload.(string), messages)}
}

func handleCommand(command string, messages []RESPMessage) string {
	if command == "echo" {
		return messages[1].Payload.(string)
	}
	if command == "ping" {
		return "pong"
	}

	return "unsupported"
}

func HandleConnection(conn net.Conn) {
	// conn represents a network connectio
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection", err)
		}
	}(conn)

	fmt.Println("Client connected:", conn.RemoteAddr())

	// create a buffered reader (reader) and a buffered writer (writer)
	// for the given connection.
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {

		message, err := Deserialize(reader)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		fmt.Println("Received message:", message)

		// here we do some business logic
		response := handleMessage(message)
		fmt.Println("Sending response:", response.Payload)

		_, err = writer.Write(Serialize(response))
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error shutting down server:", err)
		}
	}(listener)

	fmt.Println("Server is listening on port 6379.")

	// todo: implement graceful shutdown.
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go HandleConnection(conn)
	}
}
