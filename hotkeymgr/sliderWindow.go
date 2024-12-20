package hotkeymgr

import (
	"container/heap"
	"sync"
	"time"
)

type bucket struct {
	mu        sync.RWMutex
	items     map[string]int64
	updatedAt int64
}

type SliderWindow struct {
	buckets  []*bucket     // 每个桶的访问统计
	size     int           // 滑动窗口的桶数
	duration time.Duration // 滑动窗口的时间长度
}

func NewSliderWindow(size int, windowsDuration time.Duration) *SliderWindow {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = &bucket{
			items: make(map[string]int64),
		}
	}

	return &SliderWindow{
		buckets:  buckets,
		size:     size,
		duration: windowsDuration,
	}
}

func (sw *SliderWindow) GetBucketIndex(currentTime int64) int {
	return int(currentTime / int64(sw.duration) % int64(sw.size))
}

// Add 更新当前时间的桶
func (sw *SliderWindow) Add(key string) {
	Now := time.Now().Unix()
	index := sw.GetBucketIndex(Now)

	sw.buckets[index].mu.Lock()
	defer sw.buckets[index].mu.Unlock()

	// 过期的数据需要被重置
	if sw.buckets[index].updatedAt < Now-int64(sw.duration.Seconds()) {
		sw.buckets[index].items = make(map[string]int64)
		sw.buckets[index].updatedAt = Now
	}

	sw.buckets[index].items[key]++
}

// MergeBucketsToHeap 将所有桶中的商品数据合并到大顶堆中
func (sw *SliderWindow) MergeBucketsToHeap(h *MaxHeap) {
	Now := time.Now().Unix()
	for index, bucket := range sw.buckets {
		sw.buckets[index].mu.Lock()
		// 过期的数据不生效
		if sw.buckets[index].updatedAt >= Now-int64(sw.duration.Seconds()) {
			for key, count := range bucket.items {
				heap.Push(h, HotKey{Key: key, Count: count})
			}
		}
		sw.buckets[index].mu.Unlock()
	}
}
