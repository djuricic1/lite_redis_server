package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type RESPType byte

const (
	SimpleString RESPType = '+'
	Error        RESPType = '-'
	Integer      RESPType = ':'
	BulkString   RESPType = '$'
	Array        RESPType = '*'
)

// RESPMessage represents a RESP message.
type RESPMessage struct {
	Type    RESPType
	Payload interface{}
}

// String returns the string representation of the RESPMessage.
func (msg RESPMessage) String() string {
	return formatMessage(msg, 0)
}

// formatMessage is a helper function to format the RESPMessage recursively.
func formatMessage(msg RESPMessage, indentLevel int) string {
	// If Payload is also a RESPMessage, format it recursively
	var payloadStr string
	if nestedMsg, ok := msg.Payload.(RESPMessage); ok {
		payloadStr = formatMessage(nestedMsg, indentLevel+1)
	} else {
		payloadStr = fmt.Sprintf("%v", msg.Payload)
	}

	return fmt.Sprintf("RESPMessage{Type: %v, Payload: %s}", msg.Type, payloadStr)
}

func Serialize(message RESPMessage) []byte {
	// the first character is always one of special characters in RESP protocol.
	var buffer bytes.Buffer
	buffer.WriteByte(byte(message.Type))

	switch message.Type {
	case SimpleString, Error:
		buffer.WriteString(message.Payload.(string) + "\r\n")
	case Integer:
		buffer.WriteString(fmt.Sprintf("%d\r\n", message.Payload))
	case BulkString:
		if message.Payload == nil {
			buffer.WriteString(fmt.Sprintf("-1\r\n"))
		} else {
			payload := message.Payload.(string)
			buffer.WriteString(fmt.Sprintf("%d\r\n%s\r\n", len(payload), payload))
		}
	case Array:
		if message.Payload == nil {
			buffer.WriteString(fmt.Sprintf("-1\r\n"))
		} else {
			array := message.Payload.([]RESPMessage)
			buffer.WriteString(fmt.Sprintf("%d\r\n", len(array)))
			for _, element := range array {
				buffer.Write(Serialize(element))
			}
		}
	}

	return buffer.Bytes()
}

func readInt(reader *bufio.Reader) (int, error) {
	var value int
	_, err := fmt.Fscanf(reader, "%d\r\n", &value)
	return value, err
}

// Deserialize reads a RESP message from the provided reader.
func Deserialize(reader *bufio.Reader) (RESPMessage, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return RESPMessage{}, err
	}

	messageType := RESPType(firstByte)
	switch messageType {
	case SimpleString, Error:
		payload, err := reader.ReadString('\n')
		if err != nil {
			return RESPMessage{}, err
		}
		return RESPMessage{Type: messageType, Payload: payload[:len(payload)-2]}, nil
	case Integer:
		var payload int
		_, err := fmt.Fscanf(reader, "%d\r\n", &payload)
		if err != nil {
			return RESPMessage{}, err
		}
		return RESPMessage{Type: Integer, Payload: payload}, nil
	case BulkString:
		length, err := readInt(reader)
		if err != nil {
			return RESPMessage{}, err
		}
		if length == -1 {
			return RESPMessage{Type: BulkString, Payload: nil}, nil
		}
		payload := make([]byte, length)
		_, err = io.ReadFull(reader, payload)
		if err != nil {
			return RESPMessage{}, err
		}
		// Consume trailing '\r\n'
		_, err = reader.Discard(2)
		if err != nil {
			return RESPMessage{}, err
		}
		return RESPMessage{Type: BulkString, Payload: string(payload)}, nil
	case Array:
		length, err := readInt(reader)
		if err != nil {
			return RESPMessage{}, err
		}
		var payload []RESPMessage
		for i := 0; i < length; i++ {
			element, err := Deserialize(reader)
			if err != nil {
				return RESPMessage{}, err
			}
			payload = append(payload, element)
		}
		return RESPMessage{Type: Array, Payload: payload}, nil
	default:
		return RESPMessage{}, fmt.Errorf("unknown message type: %c", firstByte)
	}
}
