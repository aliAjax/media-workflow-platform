Bug: 批量清理把删除失败静默吞掉，成功计数包含失败项，并将锁释放推迟到函数结束且无法观察释放状态，导致报告与真实结果不一致。

触发: 清理一组包含不可删除路径和多个锁的资源，检查 CleanupReport、完成计数和每个 FileLock 的释放状态。

错误信息: `lock cleanup count lost`、清理错误列表为空、完成计数等于输入数量或 `IsReleased` 仍为 false。
