package block

// RPCArgInfo .
type RPCArgInfo struct {
	MID int64
}

// RPCArgBatchInfo .
type RPCArgBatchInfo struct {
	MIDs []int64
}

// RPCResInfo .
type RPCResInfo struct {
	MID         int64
	BlockStatus BlockStatus
	StartTime   int64
	EndTime     int64
}

// Parse .
func (r *RPCResInfo) Parse(b *BlockInfo) {
	r.MID = b.MID
	r.BlockStatus = b.BlockStatus
	r.StartTime = b.StartTime
	r.EndTime = b.EndTime
}

// RPCArgBlock .
// type RPCArgBlock struct {
// }

// RPCArgBatchBlock .
// type RPCArgBatchBlock struct {
// }

// RPCArgRemove .
// type RPCArgRemove struct {
// }

// RPCArgBatchRemove .RPCArgBatchRemove
// type RPCArgBatchRemove struct {
// }
