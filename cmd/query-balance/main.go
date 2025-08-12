package main

import (
	"context"   // 上下文管理，用于控制请求的生命周期
	"fmt"       // 格式化输入输出
	"log"       // 日志记录
	"math"      // 数学运算
	"math/big"  // 大数运算，处理以太币数量等大整数

	"github.com/ethereum/go-ethereum/common"    // 以太坊通用工具
	"github.com/ethereum/go-ethereum/ethclient" // 以太坊客户端
)

// main函数 - 查询账户余额
// 功能：查询指定地址的ETH余额（当前余额、历史余额、待处理余额）
// 这是一个完整的余额查询实现，展示了多种余额查询方式
func main() {
	// 步骤1：连接以太坊网络
	// 使用Cloudflare提供的以太坊主网节点
	// Cloudflare是一个免费且稳定的以太坊节点提供商
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal("连接以太坊客户端失败:", err)
	}

	// 步骤2：设置要查询的账户地址
	// 这里使用一个示例地址，您可以替换为任何有效的以太坊地址
	account := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	fmt.Printf("查询地址: %s\n\n", account.Hex())

	// 步骤3：查询当前最新余额
	// BalanceAt的第三个参数为nil表示查询最新区块的余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal("查询当前余额失败:", err)
	}
	fmt.Printf("当前余额 (Wei): %s\n", balance.String())

	// 步骤4：查询指定区块高度的历史余额
	// 这里查询区块高度5532993时的余额
	// 这对于查看账户在特定时间点的余额很有用
	blockNumber := big.NewInt(5532993)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal("查询历史余额失败:", err)
	}
	fmt.Printf("区块 %s 时的余额 (Wei): %s\n", blockNumber.String(), balanceAt.String())

	// 步骤5：将Wei转换为ETH单位显示
	// Wei是以太坊的最小单位，1 ETH = 10^18 Wei
	// 为了便于阅读，我们将Wei转换为ETH
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String()) // 将big.Int转换为big.Float
	// 除以10^18将Wei转换为ETH
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("区块 %s 时的余额 (ETH): %s\n", blockNumber.String(), ethValue.String())

	// 步骤6：查询待处理余额
	// 待处理余额包括尚未被打包到区块中的交易影响
	// 这对于查看账户的"即将到账"余额很有用
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	if err != nil {
		log.Fatal("查询待处理余额失败:", err)
	}
	fmt.Printf("待处理余额 (Wei): %s\n", pendingBalance.String())

	// 步骤7：将当前余额也转换为ETH单位显示
	currentBalance := new(big.Float)
	currentBalance.SetString(balance.String())
	currentEthValue := new(big.Float).Quo(currentBalance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("\n=== 余额查询结果汇总 ===\n")
	fmt.Printf("当前余额: %s ETH\n", currentEthValue.String())
	fmt.Printf("历史余额: %s ETH (区块 %s)\n", ethValue.String(), blockNumber.String())

	// 小白说明：
	// 1. Wei是以太坊的最小单位，类似于"分"对于"元"
	// 2. 1 ETH = 1,000,000,000,000,000,000 Wei (10^18)
	// 3. 当前余额：最新区块中的余额
	// 4. 历史余额：指定区块高度时的余额
	// 5. 待处理余额：包含未确认交易的余额
	// 6. 区块高度：以太坊网络中区块的序号，越大越新
}