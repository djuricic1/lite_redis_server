package resp_core

import "testing"

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
