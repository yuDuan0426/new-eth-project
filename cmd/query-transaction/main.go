// 以太坊交易查询工具
// 本程序演示如何查询以太坊网络上的交易信息
// 包含多种交易查询方式和详细信息解析

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 以太坊交易查询工具 ===")
	fmt.Println("本工具演示如何查询以太坊交易信息，包括多种查询方式\n")

	// ===== 第1步：连接以太坊网络 =====
	// 连接到Sepolia测试网络
	// 注意：请替换<API_KEY>为您的实际API密钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/<API_KEY>")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✓ 成功连接到Sepolia测试网络")

	// ===== 第2步：获取网络链ID =====
	// 链ID用于识别不同的以太坊网络
	// Sepolia测试网的链ID是11155111
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("获取链ID失败:", err)
	}
	fmt.Printf("✓ 网络链ID: %s\n\n", chainID.String())

	// ===== 第3步：通过区块号查询交易 =====
	fmt.Println("=== 方法1：通过区块号查询交易 ===")
	// 指定要查询的区块号
	blockNumber := big.NewInt(5671744)
	fmt.Printf("查询区块号: %s\n", blockNumber.String())

	// 根据区块号获取区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取区块信息失败:", err)
	}
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("区块中交易数量: %d\n\n", len(block.Transactions()))

	// 遍历区块中的所有交易
	for i, tx := range block.Transactions() {
		fmt.Printf("--- 交易 #%d ---\n", i+1)
		
		// ===== 交易基本信息 =====
		fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
		fmt.Printf("转账金额: %s Wei\n", tx.Value().String())
		fmt.Printf("Gas限制: %d\n", tx.Gas())
		fmt.Printf("Gas价格: %d Wei\n", tx.GasPrice().Uint64())
		fmt.Printf("Nonce: %d\n", tx.Nonce())
		fmt.Printf("交易数据: %x\n", tx.Data())
		
		// 接收方地址
		if tx.To() != nil {
			fmt.Printf("接收方地址: %s\n", tx.To().Hex())
		} else {
			fmt.Println("接收方地址: 合约创建交易")
		}

		// ===== 第4步：恢复发送方地址 =====
		// 使用EIP155签名器从交易签名中恢复发送方地址
		if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
			fmt.Printf("发送方地址: %s\n", sender.Hex())
		} else {
			log.Fatal("恢复发送方地址失败:", err)
		}

		// ===== 第5步：获取交易回执 =====
		// 交易回执包含交易执行结果和Gas使用情况
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal("获取交易回执失败:", err)
		}

		// 交易状态：1表示成功，0表示失败
		fmt.Printf("交易状态: %d ", receipt.Status)
		if receipt.Status == 1 {
			fmt.Println("(成功)")
		} else {
			fmt.Println("(失败)")
		}
		
		// 事件日志数量
		fmt.Printf("事件日志数量: %d\n", len(receipt.Logs))
		fmt.Printf("实际Gas使用: %d\n", receipt.GasUsed)
		fmt.Printf("累积Gas使用: %d\n\n", receipt.CumulativeGasUsed)
		
		// 只显示第一个交易的详细信息
		break
	}

	// ===== 第6步：通过区块哈希查询交易 =====
	fmt.Println("=== 方法2：通过区块哈希查询交易 ===")
	// 指定要查询的区块哈希
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	fmt.Printf("查询区块哈希: %s\n", blockHash.Hex())

	// 获取指定区块中的交易数量
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal("获取交易数量失败:", err)
	}
	fmt.Printf("区块中交易数量: %d\n", count)

	// 通过索引获取区块中的特定交易
	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal("获取交易失败:", err)
		}

		fmt.Printf("交易索引 %d: %s\n", idx, tx.Hash().Hex())
		
		// 只显示第一个交易
		break
	}
	fmt.Println()

	// ===== 第7步：通过交易哈希直接查询 =====
	fmt.Println("=== 方法3：通过交易哈希直接查询 ===")
	// 指定要查询的交易哈希
	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	fmt.Printf("查询交易哈希: %s\n", txHash.Hex())

	// 通过交易哈希获取交易信息
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal("获取交易信息失败:", err)
	}

	// 检查交易是否还在待处理状态
	fmt.Printf("交易是否待处理: %t\n", isPending)
	if isPending {
		fmt.Println("交易状态: 待处理(在内存池中)")
	} else {
		fmt.Println("交易状态: 已确认(已打包到区块)")
	}
	fmt.Printf("交易哈希验证: %s\n", tx.Hash().Hex())

	fmt.Println("\n=== 查询完成 ===")
	fmt.Println("\n=== 重要说明 ===")
	fmt.Println("1. 请替换API_KEY为您的实际密钥")
	fmt.Println("2. 示例中的区块号和哈希可能已过时")
	fmt.Println("3. 不同网络的区块结构可能有差异")
	fmt.Println("4. Gas价格和限制会影响交易处理速度")
	fmt.Println("5. 交易状态1表示成功，0表示失败")

	// ===== 技术说明 =====
	// 1. 交易查询方式：
	//    - 通过区块号查询：获取指定区块中的所有交易
	//    - 通过区块哈希查询：使用区块哈希获取交易
	//    - 通过交易哈希查询：直接查询特定交易
	//
	// 2. 交易信息字段：
	//    - Hash: 交易的唯一标识符
	//    - Value: 转账金额(以Wei为单位)
	//    - Gas: 交易的Gas限制
	//    - GasPrice: 每单位Gas的价格
	//    - Nonce: 发送方账户的交易序号
	//    - Data: 交易附带的数据(合约调用参数)
	//    - To: 接收方地址(nil表示合约创建)
	//
	// 3. 地址恢复：
	//    - 使用EIP155签名器从交易签名恢复发送方地址
	//    - 需要正确的链ID来防止重放攻击
	//
	// 4. 交易回执：
	//    - Status: 交易执行状态(1成功，0失败)
	//    - GasUsed: 实际消耗的Gas数量
	//    - Logs: 交易产生的事件日志
	//    - CumulativeGasUsed: 区块中累积的Gas使用量
	//
	// 5. 交易状态：
	//    - Pending: 交易在内存池中等待打包
	//    - Confirmed: 交易已被打包到区块中
	//    - Failed: 交易执行失败但仍被打包
}