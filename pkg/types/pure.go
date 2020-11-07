package types

import (
	"strconv"
)

// Action represent a type of action
type Action string

// Level represent an application status messages
type Level int

// RequestMessage represent request message
type RequestMessage struct {
	ActionStr  string `json:"action_str"`
	DataType   string `json:"data_type"`
	LogTable   string `json:"log_table"`
	RequestMap string `json:"request_map"`
	TransMap   string `json:"trans_map"`
}

// ResponseMessage represent response message
type ResponseMessage struct {
	ActionStr   string `json:"action_str"`
	DataType    string `json:"data_type"`
	LogTable    string `json:"log_table"`
	ResponseMap string `json:"response_map"`
	TransMap    string `json:"trans_map"`
}

// ConfirmMessage represents confirm message
type ConfirmMessage struct {
	ActionStr  string `json:"action_str"`
	DataType   string `json:"data_type"`
	LogTable   string `json:"log_table"`
	ConfirmMap string `json:"confirm_map"`
	TransMap   string `json:"trans_map"`
}

// List of actions types
const (
	NOOP Action = "noop"

	// Request
	CREATE   = "create"
	RETRIEVE = "retrieve"
	UPDATE   = "update"
	DELETE   = "delete"
	FLUSH    = "flush"

	// Response
	CREATED       = "CREATED"
	CREATE_FAIL   = "CREATE_FAIL"
	RETRIEVED     = "RETRIEVED"
	RETRIEVE_FAIL = "RETRIEVE_FAIL"
	UPDATED       = "UPDATED"
	UPDATE_FAIL   = "UPDATE_FAIL"
	DELETED       = "DELETED"
	DELETE_FAIL   = "DELETE_FAIL"

	// Confirmation
	ABORT = "abort"
	DONE  = "done"
	RETRY = "retry"
)

// Application status messages
const (
	LEVEL_0 Level = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
	LEVEL_4
	LEVEL_5
	LEVEL_6
	LEVEL_7
)

var levels = [...]string{
	LEVEL_0: "emerg",
	LEVEL_1: "alert",
	LEVEL_2: "crit",
	LEVEL_3: "error",
	LEVEL_4: "warn",
	LEVEL_5: "notice",
	LEVEL_6: "info",
	LEVEL_7: "debug",
}

func (lev Level) String() string {
	s := ""
	if 0 <= lev && lev < Level(len(levels)) {
		s = levels[lev]
	}
	if s == "" {
		s = "level(" + strconv.Itoa(int(lev)) + ")"
	}
	return s
}
