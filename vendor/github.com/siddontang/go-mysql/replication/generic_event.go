package replication

import (
	"encoding/hex"
	"fmt"
	"io"
)

// we don't parse all event, so some we will use GenericEvent instead
type GenericEvent struct {
	Data []byte
}

func (e *GenericEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Event data: \n%s", hex.Dump(e.Data))
	fmt.Fprintln(w)
}

func (e *GenericEvent) Decode(data []byte) error {
	e.Data = data

	return nil
}

//below events are generic events, maybe later I will consider handle some.

// type StartEventV3 struct {
// 	Version         uint16
// 	ServerVersion   [50]byte
// 	CreateTimestamp uint32
// }

// type StopEvent struct{}

// type LoadEvent struct {
// 	SlaveProxyID uint32
// 	ExecTime     uint32
// 	SkipLines    uint32
// 	TableNameLen uint8
// 	SchemaLen    uint8
// 	NumFileds    uint32
// 	FieldTerm    uint8
// 	EnclosedBy   uint8
// 	LineTerm     uint8
// 	LineStart    uint8
// 	EscapedBy    uint8
// 	OptFlags     uint8
// 	EmptyFlags   uint8

// 	//len = 1 * NumFields
// 	FieldNameLengths []byte

// 	//len = sum(FieldNameLengths) + NumFields
// 	//array of nul-terminated strings
// 	FieldNames []byte

// 	//len = TableNameLen + 1, nul-terminated string
// 	TableName []byte

// 	//len = SchemaLen + 1, nul-terminated string
// 	SchemaName []byte

// 	//string.NUL
// 	FileName []byte
// }

// type NewLoadEvent struct {
// 	SlaveProxyID  uint32
// 	ExecTime      uint32
// 	SkipLines     uint32
// 	TableNameLen  uint8
// 	SchemaLen     uint8
// 	NumFields     uint32
// 	FieldTermLen  uint8
// 	FieldTerm     []byte
// 	EnclosedByLen uint8
// 	EnclosedBy    []byte
// 	LineTermLen   uint8
// 	LineTerm      []byte
// 	LineStartLen  uint8
// 	LineStart     []byte
// 	EscapedByLen  uint8
// 	EscapedBy     []byte
// 	OptFlags      uint8

// 	//len = 1 * NumFields
// 	FieldNameLengths []byte

// 	//len = sum(FieldNameLengths) + NumFields
// 	//array of nul-terminated strings
// 	FieldNames []byte

// 	//len = TableNameLen, nul-terminated string
// 	TableName []byte

// 	//len = SchemaLen, nul-terminated string
// 	SchemaName []byte

// 	//string.EOF
// 	FileName []byte
// }

// type CreateFileEvent struct {
// 	FileID    uint32
// 	BlockData []byte
// }

// type AppendBlockEvent struct {
// 	FileID    uint32
// 	BlockData []byte
// }

// type ExecLoadEvent struct {
// 	FileID uint32
// }

// type BeginLoadQueryEvent struct {
// 	FileID    uint32
// 	BlockData []byte
// }

// type ExecuteLoadQueryEvent struct {
// 	SlaveProxyID     uint32
// 	ExecutionTime    uint32
// 	SchemaLength     uint8
// 	ErrorCode        uint16
// 	StatusVarsLength uint16

// 	FileID           uint32
// 	StartPos         uint32
// 	EndPos           uint32
// 	DupHandlingFlags uint8
// }

// type DeleteFileEvent struct {
// 	FileID uint32
// }

// type RandEvent struct {
// 	Seed1 uint64
// 	Seed2 uint64
// }

// type IntVarEvent struct {
// 	Type  uint8
// 	Value uint64
// }

// type UserVarEvent struct {
// 	NameLength uint32
// 	Name       []byte
// 	IsNull     uint8

// 	//if not is null
// 	Type        uint8
// 	Charset     uint32
// 	ValueLength uint32
// 	Value       []byte

// 	//if more data
// 	Flags uint8
// }

// type IncidentEvent struct {
// 	Type          uint16
// 	MessageLength uint8
// 	Message       []byte
// }

// type HeartbeatEvent struct {
// }
