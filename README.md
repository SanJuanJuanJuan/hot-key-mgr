# hot-key-mgr

## 使用情景
- 对于电商项目的某些banner场景，在运营人员无法预知成为热点的商品时，可通过本工具识别热点数据并将对应数据提级为一级缓存，避免大量的redis请求
- 对于像微博、推特等社交平台的热搜功能
- 等无法提前预知，需要根据用户在过去单位时间的行为来统计TOP N的情景
## 实现思路
> 我们维护一个滑动窗口，比如滑动窗口设置为10s，就是要统计这10s内有哪些key被高频访问，一个滑动窗口中对应多个Bucket，每个Bucket中对应一个map，map的key为商品的id，value为商品对应的请求次数。接着我们可以定时的(比如10s)去统计当前所有Buckets中的key的数据，然后把这些数据导入到大顶堆中，轻而易举的可以从大顶堆中获取topK的key，我们可以设置一个阈值，比如在一个滑动窗口时间内某一个key访问频次超过500次，就认为该key为热点key，从而自动地把该key升级为本地缓存。  
>作者：万俊峰Kevin  
>链接：https://juejin.cn/post/7113730074595541023  
>来源：稀土掘金  
>著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
> 
## 注意事项
1. 支持并发，但不支持分布式（因为时基于本地map存储的缓存）
   - 如需分布式，考虑使用redis
## 使用方法
1. import
`import "github.com/SanJuanJuanJuan/hot-key-mgr/hotkeymgr"`
2. 实例化
`hkm := hotkeymgr.NewHotKeyMgr(5, 10*time.Second)`
3. 在触发请求信息的地方调用add方法
`hkm.AddRequest("product1")`
4. 启用start方法
`go hkm.Start(10, 0, 10*time.Second)`
5. 调用GetHotKeyCache方法获取热点数据
`hkm.GetHotKeyCache()`

## 如何选择size
对热点数据区分精度越高（数据过期误差越低），应选择更大的size，但会增加内存占用
反之，对热点数据区分精度越低（数据过期误差越高），应选择更小的size，能以减少内存占用
