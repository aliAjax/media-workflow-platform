Bug: 分页、Page 克隆和命令查询返回共享底层数组，调用方 append 或修改元素会污染原任务列表和内部日志。

触发: 取得分页页、Page 副本、List 或 Recent 返回值后追加或改写元素，再读取原列表和下一批次。

错误信息: `selection aliases input` 或分页检查发现下一批次少一条、原列表元素被覆盖。
