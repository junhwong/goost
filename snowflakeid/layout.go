package snowflakeid

import "fmt"

// Layout 用于定义ID的二进制布局和起始Epoch
type Layout struct {
	// TimeBits 表示时间的位数
	TimeBits uint8
	// TopologyBits 表示工作节点拓扑的位数
	TopologyBits uint8
	// WorkerBits 表示工作节点的位数
	WorkerBits uint8
	// SequenceBits 表示序号的位数
	SequenceBits uint8
	// Epoch 起始时间戳
	Epoch uint64
}

func (l *Layout) Validate() (idMax, epochMax, topologyMax, workerMax, sequenceMax, idShift, topologyShift int64, err error) {
	bits := l.TimeBits + l.TopologyBits + l.WorkerBits + l.SequenceBits
	if l.TimeBits < 16 || bits > 63 {
		err = fmt.Errorf("bad bits")
	}

	idMax = int64(1)<<bits - 1
	epochMax = int64(1)<<l.TimeBits - 1
	topologyMax = int64(1)<<l.TopologyBits - 1
	workerMax = int64(1)<<l.WorkerBits - 1
	sequenceMax = int64(1)<<l.SequenceBits - 1
	idShift = int64(l.TopologyBits + l.WorkerBits + l.SequenceBits)
	topologyShift = int64(l.WorkerBits + l.SequenceBits)

	return
}

// DefaultLayout 符合 web(ecmascript) 数值安全的布局。
// see also: http://www.ecma-international.org/ecma-262/6.0/index.html#sec-ecmascript-language-types-number-type
var DefaultLayout = Layout{
	TimeBits:     41,
	TopologyBits: 0,
	WorkerBits:   5,
	SequenceBits: 10,
	Epoch:        EPOCH2020,
}
