# Go-Ethereum (go-eth) 常用API参考文档

本文档整理了go-ethereum库中常用的API，按功能模块分类，便于开发者快速查找和使用。

## 目录

1. [客户端连接](#客户端连接)
2. [账户和钱包管理](#账户和钱包管理)
3. [区块查询](#区块查询)
4. [交易操作](#交易操作)
5. [余额查询](#余额查询)
6. [智能合约](#智能合约)
7. [事件监听](#事件监听)
8. [Gas和费用](#gas和费用)
9. [网络信息](#网络信息)
10. [实用工具](#实用工具)

---

## 客户端连接

### 基础连接

```go
import "github.com/ethereum/go-ethereum/ethclient"

// HTTP连接
client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_PROJECT_ID")

// WebSocket连接（支持订阅）
client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID")

// 本地节点连接
client, err := ethclient.Dial("http://localhost:8545")
```

### 常用网络端点

```go
// 主网
"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
"https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY"
"https://cloudflare-eth.com"

// Sepolia测试网
"https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
"https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"

// Goerli测试网
"https://goerli.infura.io/v3/YOUR_PROJECT_ID"
"https://eth-goerli.alchemyapi.io/v2/YOUR_API_KEY"
```

---

## 账户和钱包管理

### 私钥和地址操作

```go
import (
    "crypto/ecdsa"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
)

// 生成新私钥
privateKey, err := crypto.GenerateKey()

// 从十六进制字符串加载私钥
privateKey, err := crypto.HexToECDSA("私钥字符串")

// 导出私钥为十六进制
privateKeyBytes := crypto.FromECDSA(privateKey)
privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]

// 从私钥获取公钥
publicKey := privateKey.Public()
publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

// 从公钥生成地址
address := crypto.PubkeyToAddress(*publicKeyECDSA)

// 地址格式转换
address := common.HexToAddress("0x地址字符串")
addressString := address.Hex()
```

### Nonce管理

```go
// 获取账户当前nonce
nonce, err := client.NonceAt(context.Background(), address, nil)

// 获取待处理nonce（包含内存池中的交易）
pendingNonce, err := client.PendingNonceAt(context.Background(), address)
```

---

## 区块查询

### 区块信息查询

```go
import (
    "context"
    "math/big"
    "github.com/ethereum/go-ethereum/common"
)

// 按区块号查询区块头（轻量级）
header, err := client.HeaderByNumber(context.Background(), big.NewInt(区块号))

// 按区块号查询完整区块（包含交易）
block, err := client.BlockByNumber(context.Background(), big.NewInt(区块号))

// 按区块哈希查询区块
blockHash := common.HexToHash("0x区块哈希")
block, err := client.BlockByHash(context.Background(), blockHash)

// 获取最新区块
block, err := client.BlockByNumber(context.Background(), nil)

// 获取区块中的交易数量
count, err := client.TransactionCount(context.Background(), blockHash)
```

### 区块信息字段

```go
// 区块头信息
header.Number.Uint64()    // 区块号
header.Time              // 时间戳
header.Difficulty        // 难度
header.Hash()            // 区块哈希
header.ParentHash        // 父区块哈希
header.GasLimit          // Gas限制
header.GasUsed           // Gas使用量

// 完整区块信息
block.Transactions()     // 交易列表
block.Size()            // 区块大小
block.Uncles()          // 叔块
```

---

## 交易操作

### 交易查询

```go
// 按交易哈希查询交易
txHash := common.HexToHash("0x交易哈希")
tx, isPending, err := client.TransactionByHash(context.Background(), txHash)

// 按区块和索引查询交易
tx, err := client.TransactionInBlock(context.Background(), blockHash, 交易索引)

// 获取交易回执
receipt, err := client.TransactionReceipt(context.Background(), txHash)
```

### 交易创建和发送

```go
import "github.com/ethereum/go-ethereum/core/types"

// 创建交易
tx := types.NewTransaction(
    nonce,                    // nonce
    toAddress,               // 接收地址
    amount,                  // 转账金额（Wei）
    gasLimit,                // Gas限制
    gasPrice,                // Gas价格
    data,                    // 交易数据
)

// 获取链ID
chainID, err := client.ChainID(context.Background())

// 签名交易
signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

// 发送交易
err = client.SendTransaction(context.Background(), signedTx)
```

### 交易信息字段

```go
// 交易基本信息
tx.Hash()                // 交易哈希
tx.Value()               // 转账金额
tx.Gas()                 // Gas限制
tx.GasPrice()            // Gas价格
tx.Nonce()               // Nonce
tx.Data()                // 交易数据
tx.To()                  // 接收地址

// 交易回执信息
receipt.Status           // 交易状态（1成功，0失败）
receipt.GasUsed          // 实际Gas使用量
receipt.BlockNumber      // 所在区块号
receipt.BlockHash        // 所在区块哈希
receipt.TransactionIndex // 在区块中的索引
receipt.ContractAddress  // 合约地址（合约创建交易）
receipt.Logs             // 事件日志
```

### 地址恢复

```go
import "github.com/ethereum/go-ethereum/core/types"

// 从交易签名恢复发送者地址
sender, err := types.Sender(types.NewEIP155Signer(chainID), tx)
```

---

## 余额查询

### ETH余额查询

```go
// 查询账户ETH余额
balance, err := client.BalanceAt(context.Background(), address, nil)

// 查询指定区块的余额
balance, err := client.BalanceAt(context.Background(), address, big.NewInt(区块号))

// 余额单位转换
import "math"

// Wei转ETH
fEther := new(big.Float)
fEther.SetString(balance.String())
ethValue := new(big.Float).Quo(fEther, big.NewFloat(math.Pow10(18)))
```

### ERC20代币余额查询

```go
// 需要使用智能合约调用
// 1. 准备balanceOf函数调用数据
functionSignature := []byte("balanceOf(address)")
hash := sha3.NewLegacyKeccak256()
hash.Write(functionSignature)
methodID := hash.Sum(nil)[:4]

// 2. 编码参数
paddedAddress := common.LeftPadBytes(address.Bytes(), 32)

// 3. 组合调用数据
var data []byte
data = append(data, methodID...)
data = append(data, paddedAddress...)

// 4. 调用合约
toAddress := common.HexToAddress("代币合约地址")
callMsg := ethereum.CallMsg{
    To:   &toAddress,
    Data: data,
}

result, err := client.CallContract(context.Background(), callMsg, nil)
```

---

## 智能合约

### 合约部署

```go
// 创建合约部署交易
tx := types.NewContractCreation(
    nonce,
    amount,      // 发送给合约的ETH数量
    gasLimit,
    gasPrice,
    contractBytecode, // 合约字节码
)

// 签名并发送
signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
err = client.SendTransaction(context.Background(), signedTx)

// 获取合约地址
receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
contractAddress := receipt.ContractAddress
```

### 合约调用

```go
import "github.com/ethereum/go-ethereum"

// 只读调用（不消耗Gas）
callMsg := ethereum.CallMsg{
    To:   &contractAddress,
    Data: callData, // ABI编码的函数调用数据
}
result, err := client.CallContract(context.Background(), callMsg, nil)

// 写入调用（需要发送交易）
tx := types.NewTransaction(
    nonce,
    contractAddress,
    value,
    gasLimit,
    gasPrice,
    callData,
)
```

### 合约验证

```go
// 检查地址是否为合约
bytecode, err := client.CodeAt(context.Background(), address, nil)
isContract := len(bytecode) > 0

// 获取合约字节码
code, err := client.CodeAt(context.Background(), contractAddress, big.NewInt(区块号))
```

### ABI操作

```go
import "github.com/ethereum/go-ethereum/accounts/abi"

// 解析ABI
contractABI, err := abi.JSON(strings.NewReader(abiString))

// 编码函数调用
data, err := contractABI.Pack("functionName", param1, param2)

// 解码返回值
var result []interface{}
err = contractABI.UnpackIntoInterface(&result, "functionName", returnData)
```

---

## 事件监听

### 实时事件订阅

```go
// WebSocket连接必需
client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID")

// 创建日志通道
logs := make(chan types.Log)

// 创建过滤器
query := ethereum.FilterQuery{
    Addresses: []common.Address{contractAddress},
    Topics:    [][]common.Hash{{eventSignatureHash}},
}

// 订阅日志
sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

// 监听事件
for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case vLog := <-logs:
        // 处理事件日志
        fmt.Printf("事件: %s\n", vLog.TxHash.Hex())
    }
}
```

### 历史事件查询

```go
// 查询历史事件
query := ethereum.FilterQuery{
    FromBlock: big.NewInt(起始区块),
    ToBlock:   big.NewInt(结束区块),
    Addresses: []common.Address{contractAddress},
    Topics:    [][]common.Hash{{eventSignatureHash}},
}

logs, err := client.FilterLogs(context.Background(), query)

for _, vLog := range logs {
    // 处理每个事件日志
    fmt.Printf("区块: %d, 交易: %s\n", vLog.BlockNumber, vLog.TxHash.Hex())
}
```

### 区块订阅

```go
// 订阅新区块头
headers := make(chan *types.Header)
sub, err := client.SubscribeNewHead(context.Background(), headers)

for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case header := <-headers:
        fmt.Printf("新区块: %d\n", header.Number.Uint64())
    }
}
```

---

## Gas和费用

### Gas价格和限制

```go
// 获取建议Gas价格
gasPrice, err := client.SuggestGasPrice(context.Background())

// 估算Gas使用量
toAddress := common.HexToAddress("0x接收地址")
callMsg := ethereum.CallMsg{
    From:     fromAddress,
    To:       &toAddress,
    Value:    amount,
    Data:     data,
}
gasLimit, err := client.EstimateGas(context.Background(), callMsg)

// EIP-1559 费用估算（如果网络支持）
// 获取费用历史
feeHistory, err := client.FeeHistory(context.Background(), 1, big.NewInt(-1), nil)
```

### Gas费用计算

```go
// 传统Gas费用 = gasUsed * gasPrice
totalFee := new(big.Int).Mul(big.NewInt(int64(gasUsed)), gasPrice)

// EIP-1559费用 = gasUsed * (baseFee + priorityFee)
// 其中 effectiveGasPrice = min(maxFeePerGas, baseFee + maxPriorityFeePerGas)
```

---

## 网络信息

### 网络状态

```go
// 获取链ID
chainID, err := client.ChainID(context.Background())

// 获取网络ID
networkID, err := client.NetworkID(context.Background())

// 获取同步状态
syncProgress, err := client.SyncProgress(context.Background())
if syncProgress != nil {
    fmt.Printf("同步进度: %d/%d\n", syncProgress.CurrentBlock, syncProgress.HighestBlock)
}

// 获取节点信息
peers, err := client.PeerCount(context.Background())
```

### 区块链状态

```go
// 获取最新区块号
blockNumber, err := client.BlockNumber(context.Background())

// 检查地址是否为合约
code, err := client.CodeAt(context.Background(), address, nil)
isContract := len(code) > 0
```

---

## 实用工具

### 数据格式转换

```go
import (
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "math/big"
)

// 地址转换
address := common.HexToAddress("0x地址")
addressString := address.Hex()

// 哈希转换
hash := common.HexToHash("0x哈希")
hashString := hash.Hex()

// 十六进制编码/解码
bytes := hexutil.MustDecode("0x数据")
hexString := hexutil.Encode(bytes)

// 大数操作
value := big.NewInt(1000000000000000000) // 1 ETH in Wei
valueFromString, _ := new(big.Int).SetString("1000000000000000000", 10)
```

### 单位转换

```go
import "math"

// Wei转ETH
func WeiToEther(wei *big.Int) *big.Float {
    fEther := new(big.Float)
    fEther.SetString(wei.String())
    return new(big.Float).Quo(fEther, big.NewFloat(math.Pow10(18)))
}

// ETH转Wei
func EtherToWei(eth *big.Float) *big.Int {
    truncInt, _ := eth.Mul(eth, big.NewFloat(math.Pow10(18))).Int(nil)
    return truncInt
}

// Gwei转Wei
func GweiToWei(gwei int64) *big.Int {
    return new(big.Int).Mul(big.NewInt(gwei), big.NewInt(1000000000))
}
```

### 错误处理

```go
import "github.com/ethereum/go-ethereum"

// 检查常见错误
if err == ethereum.NotFound {
    // 交易或区块未找到
}

// 检查交易是否失败
if receipt.Status == 0 {
    // 交易执行失败
}

// 检查是否为合约地址
if len(bytecode) == 0 {
    // 不是合约地址
}
```

---

## 最佳实践

### 1. 连接管理

```go
// 使用连接池
// 避免频繁创建连接
// 适当设置超时时间

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

### 2. 错误处理

```go
// 总是检查错误
if err != nil {
    log.Printf("操作失败: %v", err)
    return err
}

// 重试机制
for i := 0; i < 3; i++ {
    if err := operation(); err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(i+1))
}
```

### 3. Gas优化

```go
// 估算Gas并添加缓冲
estimatedGas, err := client.EstimateGas(context.Background(), callMsg)
gasLimit := estimatedGas * 120 / 100 // 增加20%缓冲

// 监控Gas价格
gasPrice, err := client.SuggestGasPrice(context.Background())
```

### 4. 安全考虑

```go
// 私钥安全
// - 不要在代码中硬编码私钥
// - 使用环境变量或安全的密钥管理系统
// - 在生产环境中使用硬件钱包或HSM

// 交易验证
// - 总是验证交易回执
// - 检查交易状态
// - 验证事件日志
```

---

## 常用常量

```go
// 链ID
const (
    MainnetChainID = 1
    GoerliChainID  = 5
    SepoliaChainID = 11155111
)

// 单位转换
const (
    Wei   = 1
    Gwei  = 1e9
    Ether = 1e18
)

// 零值
var (
    ZeroAddress = common.Address{}
    ZeroHash    = common.Hash{}
)
```

---

## 参考资源

- [go-ethereum官方文档](https://geth.ethereum.org/)
- [以太坊开发者文档](https://ethereum.org/developers/)
- [Solidity文档](https://docs.soliditylang.org/)
- [EIP提案](https://eips.ethereum.org/)

---

*本文档基于go-ethereum v1.10+版本编写，部分API可能在不同版本中有所差异。*