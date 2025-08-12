# New Ethereum Project

这是一个以太坊开发工具集合项目，包含了常用的区块链操作功能。

## 项目结构

```
new-eth-project/
├── cmd/
│   ├── query-block/        # 查询区块
│   ├── query-transaction/   # 查询交易
│   ├── query-receipt/       # 查询收据
│   ├── create-wallet/       # 创建新钱包
│   ├── eth-transfer/        # ETH转账
│   ├── token-transfer/      # 代币转账
│   ├── query-balance/       # 查询账户余额
│   ├── query-token-balance/ # 查询代币余额
│   ├── subscribe-blocks/    # 订阅区块
│   ├── deploy-contract/     # 部署合约
│   ├── load-contract/       # 加载合约
│   ├── execute-contract/    # 执行合约
│   └── contract-events/     # 合约事件
├── pkg/
│   └── common/              # 公共工具包
├── go.mod
└── README.md
```

## 功能模块

- **query-block**: 查询区块信息
- **query-transaction**: 查询交易详情
- **query-receipt**: 查询交易收据
- **create-wallet**: 创建新的以太坊钱包
- **eth-transfer**: 以太币转账功能
- **token-transfer**: ERC20代币转账
- **query-balance**: 查询账户ETH余额
- **query-token-balance**: 查询账户代币余额
- **subscribe-blocks**: 订阅新区块事件
- **deploy-contract**: 部署智能合约
- **load-contract**: 加载已部署的合约
- **execute-contract**: 执行合约方法
- **contract-events**: 监听合约事件

## 使用方法

每个功能模块都有独立的main.go文件，可以单独运行：

```bash
# 例如运行查询区块功能
cd cmd/query-block
go run main.go
```