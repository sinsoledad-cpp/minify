package shortener

import (
	"sync"
	"time"
)

const (
	// 起始时间戳 (2023-01-01 00:00:00 UTC)
	epoch int64 = 1672531200000

	// 机器ID所占位数
	workerIDBits int64 = 10

	// 序列号所占位数
	sequenceBits int64 = 12

	// 机器ID最大值 (1023)
	maxWorkerID int64 = -1 ^ (-1 << workerIDBits)

	// 序列号最大值 (4095)
	maxSequence int64 = -1 ^ (-1 << sequenceBits)

	// 机器ID左移位数
	workerIDShift = sequenceBits

	// 时间戳左移位数
	timestampShift = sequenceBits + workerIDBits
)

// Snowflake 雪花算法ID生成器
type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

// NewSnowflake 创建一个新的雪花算法ID生成器
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		workerID = workerID & maxWorkerID
	}

	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// 如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过，抛出异常
	if now < s.timestamp {
		// 这里可以根据实际情况处理时钟回拨问题
		// 简单处理：等待直到时间追上来
		for now <= s.timestamp {
			now = time.Now().UnixMilli()
		}
	}

	// 如果是同一时间生成的，则进行序列号递增
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & maxSequence
		// 序列号已经达到最大值，等待下一毫秒
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 时间戳改变，序列重置
		s.sequence = 0
	}

	s.timestamp = now

	// 生成ID
	id := ((now - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}