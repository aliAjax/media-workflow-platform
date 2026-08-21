Bug: 失败任务重排队并完成后，状态转换表、完成条件和保留策略仍按旧失败状态判断，成功任务无法稳定进入并保留终态。

触发: 让 Job 从 failed 重新进入 queued、running 再 completed，随后查询 retryable 与 terminal 状态。

错误信息: 目标检查报告失败任务无法 requeue、running 不能 complete，或 succeeded 未被识别为 terminal。
