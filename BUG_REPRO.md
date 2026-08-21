Bug: 并发读取资产、Pipeline、Job 和产物列表时，读路径与写路径同时访问共享状态，可能触发 Go 的 data race 并使服务退出。

触发: 在列表遍历尚未完成时并发创建或更新对应资源，重复执行并发列表操作即可复现；停止调度器时重复停止也会暴露停止状态问题。

错误信息: `WARNING: DATA RACE`，并可能伴随 `fatal error: concurrent map read and map write`；重复停止时会出现 `panic: close of closed channel`。
