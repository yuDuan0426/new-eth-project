# go-ethereum ethclient.Client 方法总结

本文档详细总结了 `github.com/ethereum/go-ethereum/ethclient` 包中 `Client` 类型的所有主要方法，帮助开发者快速了解和使用以太坊客户端的各种功能。

## 目录

1. [连接和基础信息](#连接和基础信息)
2. [区块相关方法](#区块相关方法)
3. [交易相关方法](#交易相关方法)
4. [收据相关方法](#收据相关方法)
5. [账户和余额查询](#账户和余额查询)
6. [合约调用方法](#合约调用方法)
7. [Gas 相关方法](#gas-相关方法)
8. [网络和同步状态](#网络和同步状态)
9. [事件日志查询](#事件日志查询)
10. [其他实用方法](#其他实用方法)

---

## 连接和基础信息

### 1. Dial() - 连接以太坊节点

```go
client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
```

- **功能**: 连接到以太坊节点（主网、测试网或本地节点）
- **参数**: RPC URL 字符串
- **返回**: Client 实例和错误信息

### 2. ChainID() - 获取链ID

```go
chainID, err := client.ChainID(context.Background())
```

- **功能**: 获取当前网络的链ID（主网=1，Sepolia=11155111等）
- **用途**: 用于交易签名和网络识别

### 3. NetworkID() - 获取网络ID

```go
networkID, err := client.NetworkID(context.Background())
```

- **功能**: 获取网络ID（通常与链ID相同）

---

## 区块相关方法

### 4. BlockNumber() - 获取最新区块号

```go
blockNumber, err := client.BlockNumber(context.Background())
```

- **功能**: 获取当前最新的区块号
- **返回**: uint64 类型的区块号

### 5. BlockByNumber() - 通过区块号获取区块

```go
block, err := client.BlockByNumber(context.Background(), big.NewInt(18000000))
```

- **功能**: 根据区块号获取完整区块信息
- **参数**: 区块号（nil表示最新区块）
- **返回**: Block 结构体，包含区块头、交易列表等

### 6. BlockByHash() - 通过区块哈希获取区块

```go
block, err := client.BlockByHash(context.Background(), blockHash)
```

- **功能**: 根据区块哈希获取区块信息

### 7. HeaderByNumber() - 获取区块头

```go
header, err := client.HeaderByNumber(context.Background(), big.NewInt(18000000))
```

- **功能**: 只获取区块头信息（不包含交易列表）
- **优势**: 比获取完整区块更快，数据量更小

### 8. HeaderByHash() - 通过哈希获取区块头

```go
header, err := client.HeaderByHash(context.Background(), blockHash)
```

---

## 交易相关方法

### 9. TransactionByHash() - 通过哈希获取交易

```go
tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
```

- **功能**: 根据交易哈希获取交易详情
- **返回**: 交易对象、是否待处理、错误信息

### 10. TransactionInBlock() - 获取区块中的特定交易

```go
tx, err := client.TransactionInBlock(context.Background(), blockHash, txIndex)
```

- **功能**: 根据区块哈希和交易索引获取交易

### 11. TransactionCount() - 获取区块中的交易数量

```go
count, err := client.TransactionCount(context.Background(), blockHash)
```

- **功能**: 获取指定区块包含的交易数量

### 12. SendTransaction() - 发送交易

```go
err := client.SendTransaction(context.Background(), signedTx)
```

- **功能**: 将已签名的交易广播到网络
- **注意**: 交易必须先签名

### 13. PendingTransactionCount() - 获取待处理交易数量

```go
count, err := client.PendingTransactionCount(context.Background())
```

- **功能**: 获取交易池中待处理的交易数量

---

## 收据相关方法

### 14. TransactionReceipt() - 获取交易收据

```go
receipt, err := client.TransactionReceipt(context.Background(), txHash)
```

- **功能**: 获取已确认交易的执行收据
- **包含**: 执行状态、Gas使用量、事件日志等

### 15. BlockReceipts() - 批量获取区块收据

```go
receipts, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(blockNum))
```

- **功能**: 获取指定区块中所有交易的收据
- **优势**: 批量获取，效率更高

---

## 账户和余额查询

### 16. BalanceAt() - 查询账户余额

```go
balance, err := client.BalanceAt(context.Background(), address, nil)
```

- **功能**: 查询指定地址在指定区块的ETH余额
- **参数**: 地址、区块号（nil表示最新）
- **返回**: *big.Int 类型的余额（单位：Wei）

### 17. NonceAt() - 获取账户Nonce

```go
nonce, err := client.NonceAt(context.Background(), address, nil)
```

- **功能**: 获取账户的交易计数（Nonce）
- **用途**: 用于构造新交易

### 18. PendingNonceAt() - 获取待处理Nonce

```go
nonce, err := client.PendingNonceAt(context.Background(), address)
```

- **功能**: 获取包含待处理交易的Nonce
- **用途**: 连续发送多笔交易时使用

### 19. CodeAt() - 获取合约代码

```go
code, err := client.CodeAt(context.Background(), contractAddress, nil)
```

- **功能**: 获取指定地址的合约字节码
- **用途**: 验证地址是否为合约地址

### 20. StorageAt() - 读取合约存储

```go
value, err := client.StorageAt(context.Background(), contractAddress, storageKey, nil)
```

- **功能**: 读取合约的存储槽数据
- **参数**: 合约地址、存储键、区块号

---

## 合约调用方法

### 21. CallContract() - 调用合约方法

```go
result, err := client.CallContract(context.Background(), callMsg, nil)
```

- **功能**: 执行只读的合约方法调用
- **用途**: 查询合约状态，不消耗Gas

### 22. PendingCallContract() - 待处理状态调用

```go
result, err := client.PendingCallContract(context.Background(), callMsg)
```

- **功能**: 在待处理状态下调用合约

### 23. EstimateGas() - 估算Gas消耗

```go
gas, err := client.EstimateGas(context.Background(), callMsg)
```

- **功能**: 估算交易或合约调用需要的Gas量
- **用途**: 在发送交易前预估成本

---

## Gas 相关方法

### 24. SuggestGasPrice() - 建议Gas价格

```go
gasPrice, err := client.SuggestGasPrice(context.Background())
```

- **功能**: 获取网络建议的Gas价格
- **用途**: 设置交易的Gas价格

### 25. SuggestGasTipCap() - 建议小费上限

```go
tipCap, err := client.SuggestGasTipCap(context.Background())
```

- **功能**: 获取EIP-1559交易的建议小费上限
- **用途**: 用于Type 2交易

### 26. FeeHistory() - 获取费用历史

```go
feeHistory, err := client.FeeHistory(context.Background(), blockCount, lastBlock, rewardPercentiles)
```

- **功能**: 获取历史区块的费用信息
- **用途**: 分析Gas价格趋势

---

## 网络和同步状态

### 27. SyncProgress() - 获取同步进度

```go
progress, err := client.SyncProgress(context.Background())
```

- **功能**: 获取节点同步进度
- **返回**: 同步状态信息或nil（已同步）

### 28. PeerCount() - 获取对等节点数量

```go
peerCount, err := client.PeerCount(context.Background())
```

- **功能**: 获取连接的对等节点数量

---

## 事件日志查询

### 29. FilterLogs() - 过滤日志

```go
logs, err := client.FilterLogs(context.Background(), filterQuery)
```

- **功能**: 根据过滤条件查询事件日志
- **参数**: FilterQuery 结构体

### 30. SubscribeFilterLogs() - 订阅日志

```go
sub, err := client.SubscribeFilterLogs(context.Background(), filterQuery, logsCh)
```

- **功能**: 实时订阅符合条件的事件日志
- **用途**: 监听合约事件

### 31. SubscribeNewHead() - 订阅新区块头

```go
sub, err := client.SubscribeNewHead(context.Background(), headersCh)
```

- **功能**: 订阅新区块头
- **用途**: 实时监听新区块

---

## 其他实用方法

### 32. Close() - 关闭连接

```go
client.Close()
```

- **功能**: 关闭与以太坊节点的连接
- **用途**: 释放资源

### 33. Client() - 获取底层RPC客户端

```go
rpcClient := client.Client()
```

- **功能**: 获取底层的RPC客户端
- **用途**: 执行自定义RPC调用

---

## 使用示例

### 基本连接和查询

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊节点
    client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 获取链ID
    chainID, err := client.ChainID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("链ID: %s\n", chainID.String())
    
    // 获取最新区块号
    blockNumber, err := client.BlockNumber(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("最新区块号: %d\n", blockNumber)
    
    // 查询账户余额
    address := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9db96c4b4d8b6")
    balance, err := client.BalanceAt(context.Background(), address, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("账户余额: %s Wei\n", balance.String())
}
```

### 合约调用示例

```go
// 调用合约的只读方法
func callContractMethod(client *ethclient.Client, contractAddr common.Address) {
    // 构造调用消息
    callMsg := ethereum.CallMsg{
        To:   &contractAddr,
        Data: methodData, // ABI编码的方法调用数据
    }
    
    // 执行调用
    result, err := client.CallContract(context.Background(), callMsg, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 解码结果
    fmt.Printf("调用结果: %x\n", result)
}
```

---

## 最佳实践

1. **连接管理**
   - 使用连接池管理多个连接
   - 及时关闭不用的连接
   - 处理连接断开和重连

2. **错误处理**
   - 区分网络错误和业务错误
   - 实现重试机制
   - 记录详细的错误日志

3. **性能优化**
   - 使用批量查询减少网络请求
   - 缓存不变的数据（如历史区块）
   - 选择合适的查询方法（如HeaderByNumber vs BlockByNumber）

4. **安全考虑**
   - 验证输入参数
   - 使用HTTPS连接
   - 保护私钥和敏感信息

---

## 常用常量

```go
// 网络链ID
const (
    MainnetChainID = 1
    SepoliaChainID = 11155111
    GoerliChainID  = 5
)

// 单位转换
const (
    Wei   = 1
    GWei  = 1e9
    Ether = 1e18
)
```

---

*本文档基于 go-ethereum v1.13.x 版本整理，具体API可能因版本而异。*