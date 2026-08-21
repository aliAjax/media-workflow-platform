Bug: 首次启动没有显式 blob 目录时，默认存储根目录和多个零值构造未初始化内部状态，创建资产或写入配额会落到空路径并触发 nil map panic。

触发: 使用默认构造创建 blob、上传会话、配额账本和资源标签，再执行第一次写入；显式提供目录的路径不受影响。

错误信息: `panic: assignment to entry in nil map`，栈中可见 `internal/domain.(*QuotaLedger).Set`。
