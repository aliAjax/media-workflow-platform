Bug: 停止阶段的 FairQueue 读写未统一受锁保护，空队列返回哨兵任务，重试次数边界也多执行一次，worker 因此可能竞态并消费空任务。

触发: 并发 Push、Peek、Reset 和空队列 Pop，同时让重试策略达到最大次数，再观察停止窗口的 worker 行为。

错误信息: `WARNING: DATA RACE`，以及目标检查报告空任务 ID 不为空或重试次数超过上限。
