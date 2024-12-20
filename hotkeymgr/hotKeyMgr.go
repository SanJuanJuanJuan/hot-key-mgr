package hotkeymgr

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type HotKeyMgr struct {
	sliderWindow *SliderWindow    // 滑动窗口
	heap         *MaxHeap         // 大顶堆存储的topK数据
	cache        map[string]int64 // 本地缓存的热点数据
	mutex        sync.RWMutex     // 读写锁，保证并发性
}

func NewHotKeyMgr(size int, duration time.Duration) *HotKeyMgr {
	return &HotKeyMgr{
		sliderWindow: NewSliderWindow(size, duration),
		heap:         NewMaxHeap(),
		cache:        make(map[string]int64),
	}
}

func (hkm *HotKeyMgr) AddRequest(key string) {
	hkm.sliderWindow.Add(key)
}

func (hkm *HotKeyMgr) MergeBucketsToHeap() {
	hkm.mutex.Lock()
	defer hkm.mutex.Unlock()

	hkm.sliderWindow.MergeBucketsToHeap(hkm.heap)
}

func (hkm *HotKeyMgr) UpdateTopK(topK int, threshold int64) {
	hkm.mutex.Lock()
	defer hkm.mutex.Unlock()

	hkm.cache = make(map[string]int64)

	for hkm.heap.Len() > 0 && topK > 0 {
		hotKey := heap.Pop(hkm.heap).(HotKey)
		if hotKey.Count >= threshold {
			hkm.cache[hotKey.Key] = hotKey.Count
			topK--
		}
	}
}

func (hkm *HotKeyMgr) GetHotKeyFromCache(key string) (int64, bool) {
	hkm.mutex.RLock()
	defer hkm.mutex.RUnlock()

	count, exists := hkm.cache[key]
	return count, exists
}

func (hkm *HotKeyMgr) GetHotKeyCache() map[string]int64 {
	hkm.mutex.RLock()
	defer hkm.mutex.RUnlock()
	return hkm.cache
}

func (hkm *HotKeyMgr) Start(topK, threshold int, duration time.Duration) {
	ticket := time.NewTicker(duration)
	for {
		select {
		case <-ticket.C:
			fmt.Println("ticker")
			hkm.MergeBucketsToHeap()
			hkm.UpdateTopK(topK, int64(threshold))
			hkm.heap = NewMaxHeap()
		}
	}
}
