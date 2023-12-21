package resp_core

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		testDescription string
		message         RESPMessage
		expected        string
	}{
		{"Null Bulk String", RESPMessage{Type: BulkString, Payload: nil}, "$-1\r\n"},
		{"Array with Ping", RESPMessage{Type: Array, Payload: []RESPMessage{{Type: BulkString, Payload: "ping"}}}, "*1\r\n$4\r\nping\r\n"},
		{"Array with Echo", RESPMessage{Type: Array, Payload: []RESPMessage{{Type: BulkString, Payload: "echo"}, {Type: BulkString, Payload: "hello world"}}}, "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n"},
		{"Array with Get Key", RESPMessage{Type: Array, Payload: []RESPMessage{{Type: BulkString, Payload: "get"}, {Type: BulkString, Payload: "key"}}}, "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n"},
		{"Simple String OK", RESPMessage{Type: SimpleString, Payload: "OK"}, "+OK\r\n"},
		{"Error Message", RESPMessage{Type: Error, Payload: "Error message"}, "-Error message\r\n"},
		{"Empty Bulk String", RESPMessage{Type: BulkString, Payload: ""}, "$0\r\n\r\n"},
		{"Non-Null Bulk String", RESPMessage{Type: BulkString, Payload: "hello world"}, "$11\r\nhello world\r\n"},
		{"Simple String Hello World", RESPMessage{Type: SimpleString, Payload: "hello world"}, "+hello world\r\n"},
	}

	for _, tt := range tests {
		t.Run(
			tt.testDescription,
			func(t *testing.T) {
				serialized := Serialize(tt.message)
				if result := string(serialized); result != tt.expected {
					t.Errorf("Test Case: %s\nActual  : %s\nExpected: %s", tt.testDescription, result, tt.expected)
				}
			},
		)
	}
}

func TestDeserialize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RESPMessage
	}{
		{"Simple String", "+hello world\r\n", RESPMessage{Type: SimpleString, Payload: "hello world"}},
		{"Error", "-Error message\r\n", RESPMessage{Type: Error, Payload: "Error message"}},
		{"Positive Integer", ":42\r\n", RESPMessage{Type: Integer, Payload: 42}},
		{"Negative Integer", ":-123\r\n", RESPMessage{Type: Integer, Payload: -123}},
		{"Bulk String", "$11\r\nhello world\r\n", RESPMessage{Type: BulkString, Payload: "hello world"}},
		{"Null Bulk String", "$-1\r\n", RESPMessage{Type: BulkString, Payload: nil}},
		{"Empty Bulk String", "$0\r\n\r\n", RESPMessage{Type: BulkString, Payload: ""}},
		{"Array with Simple Strings", "*3\r\n+one\r\n+two\r\n+three\r\n", RESPMessage{Type: Array, Payload: []RESPMessage{
			{Type: SimpleString, Payload: "one"},
			{Type: SimpleString, Payload: "two"},
			{Type: SimpleString, Payload: "three"},
		}}},
		{"Array with Null Bulk Strings", "*2\r\n$-1\r\n$-1\r\n", RESPMessage{Type: Array, Payload: []RESPMessage{
			{Type: BulkString, Payload: nil},
			{Type: BulkString, Payload: nil},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			deserialized, err := Deserialize(reader)
			if err != nil {
				t.Errorf("Test Case: %s\nError: %v", tt.name, err)
				return
			}

			if !reflect.DeepEqual(deserialized, tt.expected) {
				t.Errorf("Test Case: %s\nActual  : %+v\nExpected: %+v", tt.name, deserialized, tt.expected)
			}
		})
	}
}
