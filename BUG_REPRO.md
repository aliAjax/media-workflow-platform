Bug: 执行器和健康检查收到调用方 context 后，部分路径改用 Background context，取消和 deadline 无法传到实际工作。

触发: 先取消请求 context，再执行单阶段、阶段序列和健康检查；工作仍会继续，检查结果也会保持 ok。

错误信息: 定向检查出现 `err=<nil>`、`<nil>` 或 `Status:ok`，而调用方已经取消。
