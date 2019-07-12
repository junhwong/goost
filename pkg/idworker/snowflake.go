package idworker

import "fmt"

// twitter-snowflake 雪花算法
// 把时间戳,工作机器ID, 序列号组合成一个 64位数字.
// 第一位置零, [2,42]这41位存放时间戳,[43,52]这10位存放机器id,[53,64]最后12位存放序列号.
// 参考：
// https://blog.csdn.net/u012488504/article/details/82194495
// https://www.cnblogs.com/Hollson/p/9116218.html
// https://tech.meituan.com/2017/04/21/mt-leaf.html
// https://blog.csdn.net/X5fnncxzq4/article/details/79549514
// https://github.com/twitter-archive/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala

var (
	twepoch       int64 // 时间戳起始点，警告：中途修改只能调大，否则造成生成的ID大于历史ID。
	workerIdShift uint  = 12
	maxWorkerId   int64 = 1024
)

// Next 返回一个新的ID
func genSnowflakeId(workerID, epoch, sn int64) int64 {
	// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
	// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
	stamp := epoch & 0x1FFFFFFFFFF
	// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
	stamp <<= 22
	id := stamp | (workerID << workerIdShift) | sn
	if id <= 0 {
		fmt.Printf("Error: workerID:%v,epoch:%v,sn:%v\n", workerID, epoch, sn)
	}
	return id
}

type snowflakeIdWorker struct {
	workerID int64
}

func (g snowflakeIdWorker) NextId() int64 {
	epoch, sn := Next()
	epoch -= twepoch
	return genSnowflakeId(g.workerID, epoch, sn)
}

func NewSnowflakeIdWorker(workerID int64) (IdWorker, error) {
	if workerID > maxWorkerId || workerID < 0 {
		return nil, fmt.Errorf("worker Id can't be greater than %d or less than 0", maxWorkerId)
	}
	return &snowflakeIdWorker{
		workerID: workerID,
	}, nil
}
