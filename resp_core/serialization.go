package resp_core

import (
	"bytes"
	"fmt"
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

func Serialize(message RESPMessage) []byte {
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
