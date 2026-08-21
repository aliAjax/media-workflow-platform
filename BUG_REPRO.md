Bug: 构建产物清单和质量结果时，返回切片与调用方输入共享底层数组，原地排序、过滤或修改返回值会污染之前保存的数据。

触发: 传入有容量余量的产物或警告切片，调用清单构建、筛选、质量评估或警告归一化，再修改返回切片并检查原输入。

错误信息: `caller slice reordered`、`selection aliases input`、`quality warnings alias input`、`normalized warnings alias input`。
