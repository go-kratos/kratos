package replication

const (
	//we only support MySQL 5.0.0+ binlog format, maybe???
	MinBinlogVersion = 4
)

var (
	//binlog header [ fe `bin` ]
	BinLogFileHeader []byte = []byte{0xfe, 0x62, 0x69, 0x6e}

	SemiSyncIndicator byte = 0xef
)

const (
	LOG_EVENT_BINLOG_IN_USE_F            uint16 = 0x0001
	LOG_EVENT_FORCED_ROTATE_F            uint16 = 0x0002
	LOG_EVENT_THREAD_SPECIFIC_F          uint16 = 0x0004
	LOG_EVENT_SUPPRESS_USE_F             uint16 = 0x0008
	LOG_EVENT_UPDATE_TABLE_MAP_VERSION_F uint16 = 0x0010
	LOG_EVENT_ARTIFICIAL_F               uint16 = 0x0020
	LOG_EVENT_RELAY_LOG_F                uint16 = 0x0040
	LOG_EVENT_IGNORABLE_F                uint16 = 0x0080
	LOG_EVENT_NO_FILTER_F                uint16 = 0x0100
	LOG_EVENT_MTS_ISOLATE_F              uint16 = 0x0200
)

const (
	BINLOG_DUMP_NEVER_STOP  uint16 = 0x00
	BINLOG_DUMP_NON_BLOCK   uint16 = 0x01
	BINLOG_THROUGH_POSITION uint16 = 0x02
	BINLOG_THROUGH_GTID     uint16 = 0x04
)

const (
	BINLOG_ROW_IMAGE_FULL    = "FULL"
	BINLOG_ROW_IAMGE_MINIMAL = "MINIMAL"
	BINLOG_ROW_IMAGE_NOBLOB  = "NOBLOB"
)

type EventType byte

const (
	UNKNOWN_EVENT EventType = iota
	START_EVENT_V3
	QUERY_EVENT
	STOP_EVENT
	ROTATE_EVENT
	INTVAR_EVENT
	LOAD_EVENT
	SLAVE_EVENT
	CREATE_FILE_EVENT
	APPEND_BLOCK_EVENT
	EXEC_LOAD_EVENT
	DELETE_FILE_EVENT
	NEW_LOAD_EVENT
	RAND_EVENT
	USER_VAR_EVENT
	FORMAT_DESCRIPTION_EVENT
	XID_EVENT
	BEGIN_LOAD_QUERY_EVENT
	EXECUTE_LOAD_QUERY_EVENT
	TABLE_MAP_EVENT
	WRITE_ROWS_EVENTv0
	UPDATE_ROWS_EVENTv0
	DELETE_ROWS_EVENTv0
	WRITE_ROWS_EVENTv1
	UPDATE_ROWS_EVENTv1
	DELETE_ROWS_EVENTv1
	INCIDENT_EVENT
	HEARTBEAT_EVENT
	IGNORABLE_EVENT
	ROWS_QUERY_EVENT
	WRITE_ROWS_EVENTv2
	UPDATE_ROWS_EVENTv2
	DELETE_ROWS_EVENTv2
	GTID_EVENT
	ANONYMOUS_GTID_EVENT
	PREVIOUS_GTIDS_EVENT
)

const (
	// MariaDB event starts from 160
	MARIADB_ANNOTATE_ROWS_EVENT EventType = 160 + iota
	MARIADB_BINLOG_CHECKPOINT_EVENT
	MARIADB_GTID_EVENT
	MARIADB_GTID_LIST_EVENT
)

func (e EventType) String() string {
	switch e {
	case UNKNOWN_EVENT:
		return "UnknownEvent"
	case START_EVENT_V3:
		return "StartEventV3"
	case QUERY_EVENT:
		return "QueryEvent"
	case STOP_EVENT:
		return "StopEvent"
	case ROTATE_EVENT:
		return "RotateEvent"
	case INTVAR_EVENT:
		return "IntVarEvent"
	case LOAD_EVENT:
		return "LoadEvent"
	case SLAVE_EVENT:
		return "SlaveEvent"
	case CREATE_FILE_EVENT:
		return "CreateFileEvent"
	case APPEND_BLOCK_EVENT:
		return "AppendBlockEvent"
	case EXEC_LOAD_EVENT:
		return "ExecLoadEvent"
	case DELETE_FILE_EVENT:
		return "DeleteFileEvent"
	case NEW_LOAD_EVENT:
		return "NewLoadEvent"
	case RAND_EVENT:
		return "RandEvent"
	case USER_VAR_EVENT:
		return "UserVarEvent"
	case FORMAT_DESCRIPTION_EVENT:
		return "FormatDescriptionEvent"
	case XID_EVENT:
		return "XIDEvent"
	case BEGIN_LOAD_QUERY_EVENT:
		return "BeginLoadQueryEvent"
	case EXECUTE_LOAD_QUERY_EVENT:
		return "ExectueLoadQueryEvent"
	case TABLE_MAP_EVENT:
		return "TableMapEvent"
	case WRITE_ROWS_EVENTv0:
		return "WriteRowsEventV0"
	case UPDATE_ROWS_EVENTv0:
		return "UpdateRowsEventV0"
	case DELETE_ROWS_EVENTv0:
		return "DeleteRowsEventV0"
	case WRITE_ROWS_EVENTv1:
		return "WriteRowsEventV1"
	case UPDATE_ROWS_EVENTv1:
		return "UpdateRowsEventV1"
	case DELETE_ROWS_EVENTv1:
		return "DeleteRowsEventV1"
	case INCIDENT_EVENT:
		return "IncidentEvent"
	case HEARTBEAT_EVENT:
		return "HeartbeatEvent"
	case IGNORABLE_EVENT:
		return "IgnorableEvent"
	case ROWS_QUERY_EVENT:
		return "RowsQueryEvent"
	case WRITE_ROWS_EVENTv2:
		return "WriteRowsEventV2"
	case UPDATE_ROWS_EVENTv2:
		return "UpdateRowsEventV2"
	case DELETE_ROWS_EVENTv2:
		return "DeleteRowsEventV2"
	case GTID_EVENT:
		return "GTIDEvent"
	case ANONYMOUS_GTID_EVENT:
		return "AnonymousGTIDEvent"
	case PREVIOUS_GTIDS_EVENT:
		return "PreviousGTIDsEvent"
	case MARIADB_ANNOTATE_ROWS_EVENT:
		return "MariadbAnnotateRowsEvent"
	case MARIADB_BINLOG_CHECKPOINT_EVENT:
		return "MariadbBinLogCheckPointEvent"
	case MARIADB_GTID_EVENT:
		return "MariadbGTIDEvent"
	case MARIADB_GTID_LIST_EVENT:
		return "MariadbGTIDListEvent"

	default:
		return "UnknownEvent"
	}
}

const (
	BINLOG_CHECKSUM_ALG_OFF byte = 0 // Events are without checksum though its generator
	// is checksum-capable New Master (NM).
	BINLOG_CHECKSUM_ALG_CRC32 byte = 1 // CRC32 of zlib algorithm.
	//  BINLOG_CHECKSUM_ALG_ENUM_END,  // the cut line: valid alg range is [1, 0x7f].
	BINLOG_CHECKSUM_ALG_UNDEF byte = 255 // special value to tag undetermined yet checksum
	// or events from checksum-unaware servers
)
