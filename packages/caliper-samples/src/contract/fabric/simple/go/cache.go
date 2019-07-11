package main

import (
	"sync"
	"time"
)

var cache = Cache{}

const BLOCK_TIME = 2.0

type TransactionResult struct {
	Balance uint64
	Request uint64
	Rest    uint64
}

type Segment struct {
	BasePage    uint64
	Count       uint64
	Created     time.Time
	TxResult    *TransactionResult
	Transaction *PrototypeTransaction
}

type SegmentMap map[int64]Segment

type Segments struct {
	Mat map[int64]Segment
	L   *sync.RWMutex
}

type Cache struct {
	Mat      map[string]Segments
	GCTicker *time.Ticker
	L        *sync.RWMutex
	Disabled bool
}

func (s *Segment) Life(group *Segments, id int64) {
	time.AfterFunc(BLOCK_TIME*time.Second, func() {
		group.L.Lock()
		defer group.L.Unlock()
		delete(group.Mat, id)
	})
}

func (s *Segments) Push(seg *Segment) {
	id := seg.Created.UnixNano()
	if s.Mat == nil {
		s.Mat = make(SegmentMap)
		s.L = &sync.RWMutex{}
	}
	s.Set(id, *seg)
	seg.Life(s, id)
}

func InitSegments(s *Segments) {
	s.L = &sync.RWMutex{}
	s.Mat = make(SegmentMap)
}

func (c *Cache) Get(address string) (Segments, bool) {
	c.L.RLock()
	defer c.L.RUnlock()
	segments, ok := c.Mat[address]
	return segments, ok
}

func (c *Cache) Set(address string, segments Segments) {
	c.L.Lock()
	defer c.L.Unlock()
	c.Mat[address] = segments
}

func (s *Segments) Get(key int64) (Segment, bool) {
	s.L.RLock()
	defer s.L.RUnlock()
	segment, ok := s.Mat[key]
	return segment, ok
}

func (s *Segments) Set(key int64, segment Segment) {
	s.L.Lock()
	defer s.L.Unlock()
	s.Mat[key] = segment
}

func (c *Cache) GarbageCollector() {
	go func() {
		for {
			<-c.GCTicker.C
			c.L.Lock()
			for key, value := range c.Mat {
				value.L.RLock()
				length := len(value.Mat)
				value.L.RUnlock()
				if length == 0 {
					value.Mat = nil
					value.L = nil
					delete(c.Mat, key)
				}
			}
			c.L.Unlock()
		}
	}()
}

func (c *Cache) Insert(address string, account *PrototypeAccount, count uint64, transaction *PrototypeTransaction, TxResult *TransactionResult) {
	if c.Disabled {
		return
	}
	segments, ok := c.Get(address)
	if ok == false {
		InitSegments(&segments)
	}
	segment := new(Segment)
	segment.BasePage = account.Page
	segment.Count = count
	segment.Created = time.Now()

	segment.Transaction = transaction
	segment.TxResult = TxResult

	segments.Push(segment)
	c.Set(address, segments)
}

func (c *Cache) CacheInit() {
	if c.L == nil {
		logger.Noticef("-- Cache Set --")
		c.GCTicker = time.NewTicker(5 * BLOCK_TIME * time.Second)
		c.Mat = make(map[string]Segments)
		c.L = &sync.RWMutex{}
		c.GarbageCollector()
		return
	}
	logger.Notice("-- Cache Running --")
	logger.Noticef("Status: Mutex: %+v, MapLength: %d, Ticker: %+v", c.L, len(c.Mat), c.GCTicker)
}

func (c *Cache) Query(address string) bool {
	_, ok := c.Get(address)
	return ok
}

// Determination
// return trustable balance
func (c *Cache) Determination(address string) uint64 {
	var balance uint64
	var requests uint64

	segments, ok := c.Get(address)
	if !ok {
		return 0
	}

	segments.L.RLock()
	for _, segment := range segments.Mat {
		if segment.TxResult.Balance > balance {
			balance = segment.TxResult.Balance
		}
		requests += segment.TxResult.Request
	}
	segments.L.RUnlock()

	return balance - requests
}

// Count Transaction in Cache
func (c *Cache) TransactionCount(address string) uint64 {
	var count uint64

	segments, ok := c.Get(address)
	if !ok {
		return 0
	}

	segments.L.RLock()
	for _, segment := range segments.Mat {
		if segment.Count > count {
			count = segment.Count
		}
	}
	segments.L.RUnlock()

	return count
}

// Disable Cache
func (c *Cache) Disable() {
	c.Disabled = true
}

// Enable Cache
func (c *Cache) Enable() {
	c.Disabled = false
}
