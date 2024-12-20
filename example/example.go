package main

import (
	"fmt"
	"hot-key-mgr/hotkeymgr"
	"time"
)

func main() {
	// 创建热点数据管理器，滑动窗口大小5，窗口时间10秒
	fmt.Println("创建热点数据管理器，滑动窗口大小5，窗口时间10秒")
	hkm := hotkeymgr.NewHotKeyMgr(5, 10*time.Second)

	// 模拟记录商品访问
	fmt.Println("模拟记录商品访问")
	hkm.AddRequest("product1")
	hkm.AddRequest("product2")
	hkm.AddRequest("product1")

	// 启动定时器，定期更新热点商品并存入本地缓存
	fmt.Println("启动定时器，定期更新热点商品并存入本地缓存")
	go hkm.Start(10, 0, 10*time.Second)

	// 等待几秒钟让定时器工作
	fmt.Println("等待几秒钟让定时器工作")
	time.Sleep(12 * time.Second)

	// 获取本地缓存中的热点商品
	fmt.Println("获取本地缓存中的热点商品")
	fmt.Println("cache:", hkm.GetHotKeyCache())

	// 模拟记录商品访问
	fmt.Println("模拟记录商品访问")
	hkm.AddRequest("product1")

	// 等待几秒钟让定时器工作
	fmt.Println("等待几秒钟让定时器工作")
	time.Sleep(12 * time.Second)

	// 获取本地缓存中的热点商品
	fmt.Println("获取本地缓存中的热点商品")
	fmt.Println("cache:", hkm.GetHotKeyCache())
}
