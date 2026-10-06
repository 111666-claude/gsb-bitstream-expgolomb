# bitstream-expgolomb

按位读一个字节缓冲：定长位段、无符号指数哥伦布码、有符号指数哥伦布码。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/bitstream --sample=bitorder
go run ./cmd/bitstream --sample=ue
go run ./cmd/bitstream --sample=se
go run ./cmd/bitstream --sample=bounds
go run ./cmd/bitstream --sample=work
```

## 口径

- **位序**：位从每个字节的最高位开始读。
- **游标**：读过的位不再读；顺序读取的总代价与位宽之和成正比。
- **UE**：先数前导零 `z`，再读 `z` 位 `b`，值等于 `2^z - 1 + b`。
- **SE**：先读 UE 得到 `k`，`k` 偶数得 `-k/2`，奇数得 `(k+1)/2`。
- **越界**：请求的位超出缓冲时报 `eof`。

## 不变量

- 读到的值只由已读位置之后、请求宽度之内的位决定。
- 越界一律报错，不返回凑出来的值。
- `scanned` 不随读取次数平方级放大：连续读四百个字节的 `scanned` 不超过 6000。

## 输出契约

`Read(n)` / `UE` / `SE` 顺序读取；出错返回 `eof`。场景打印一行 JSON。
`--sample=work` 打印 `{"reads": N, "scanned": N}`。
