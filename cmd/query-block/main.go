// 以太坊区块查询工具
// 本程序演示如何查询以太坊网络上的区块信息
// 包含区块头信息、完整区块信息和交易数量统计

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 以太坊区块查询工具 ===")
	fmt.Println("本工具演示如何查询以太坊区块信息，包括区块头和完整区块数据\n")

	// ===== 第1步：连接以太坊网络 =====
	// 连接到Sepolia测试网络
	// 注意：请替换<API_KEY>为您的实际API密钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/<API_KEY>")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✓ 成功连接到Sepolia测试网络")

	// ===== 第2步：设置查询参数 =====
	// 指定要查询的区块号
	blockNumber := big.NewInt(5671744)
	fmt.Printf("查询区块号: %s\n\n", blockNumber.String())

	// ===== 第3步：查询区块头信息 =====
	fmt.Println("=== 方法1：查询区块头信息 ===")
	// 区块头包含区块的基本元数据，不包含交易详情
	// 这种方式查询速度更快，数据量更小
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取区块头失败:", err)
	}

	// 显示区块头信息
	fmt.Printf("区块号: %d\n", header.Number.Uint64())
	fmt.Printf("时间戳: %d\n", header.Time)
	fmt.Printf("难度: %d\n", header.Difficulty.Uint64())
	fmt.Printf("区块哈希: %s\n\n", header.Hash().Hex())

	// ===== 第4步：查询完整区块信息 =====
	fmt.Println("=== 方法2：查询完整区块信息 ===")
	// 完整区块包含所有交易数据，数据量较大但信息完整
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取完整区块失败:", err)
	}

	// 显示完整区块信息
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("难度: %d\n", block.Difficulty().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("交易数量: %d\n\n", len(block.Transactions()))

	// ===== 第5步：通过区块哈希查询交易数量 =====
	fmt.Println("=== 方法3：通过区块哈希查询交易数量 ===")
	// 使用区块哈希来查询交易数量，这是另一种验证方式
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal("获取交易数量失败:", err)
	}
	fmt.Printf("通过区块哈希查询的交易数量: %d\n", count)

	// ===== 第6步：数据一致性验证 =====
	fmt.Println("\n=== 数据一致性验证 ===")
	fmt.Printf("区块头中的区块号: %d\n", header.Number.Uint64())
	fmt.Printf("完整区块中的区块号: %d\n", block.Number().Uint64())
	fmt.Printf("区块头中的哈希: %s\n", header.Hash().Hex())
	fmt.Printf("完整区块中的哈希: %s\n", block.Hash().Hex())
	fmt.Printf("完整区块中的交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("通过哈希查询的交易数量: %d\n", count)

	// 验证数据一致性
	if header.Number.Uint64() == block.Number().Uint64() &&
		header.Hash().Hex() == block.Hash().Hex() &&
		uint(len(block.Transactions())) == count {
		fmt.Println("\n✓ 所有查询结果一致，数据验证通过！")
	} else {
		fmt.Println("\n✗ 查询结果不一致，请检查网络连接或数据")
	}

	fmt.Println("\n=== 查询完成 ===")
	fmt.Println("\n=== 重要说明 ===")
	fmt.Println("1. 请替换API_KEY为您的实际密钥")
	fmt.Println("2. 区块号可以根据需要调整")
	fmt.Println("3. 区块头查询速度快，适合获取基本信息")
	fmt.Println("4. 完整区块查询包含所有交易，数据量大")
	fmt.Println("5. 难度为0表示使用权益证明(PoS)共识")

	// ===== 技术说明 =====
	// 1. 区块查询方式对比：
	//    - HeaderByNumber: 只获取区块头，速度快，数据量小
	//    - BlockByNumber: 获取完整区块，包含所有交易，数据量大
	//    - TransactionCount: 通过区块哈希获取交易数量
	//
	// 2. 区块头信息字段：
	//    - Number: 区块在区块链中的序号
	//    - Time: 区块创建的Unix时间戳
	//    - Difficulty: 挖矿难度(PoS网络中为0)
	//    - Hash: 区块的唯一标识符
	//
	// 3. 完整区块额外信息：
	//    - Transactions: 区块中包含的所有交易
	//    - GasLimit: 区块的Gas限制
	//    - GasUsed: 区块实际使用的Gas
	//    - Size: 区块的字节大小
	//
	// 4. 时间戳说明：
	//    - Unix时间戳格式，表示自1970年1月1日以来的秒数
	//    - 可以转换为人类可读的日期时间格式
	//
	// 5. 难度机制：
	//    - 工作量证明(PoW)：难度动态调整挖矿难度
	//    - 权益证明(PoS)：难度为0，使用验证者机制
	//    - Sepolia测试网使用PoS，所以难度为0
	//
	// 6. 性能考虑：
	//    - 区块头查询适合快速获取基本信息
	//    - 完整区块查询适合需要交易详情的场景
	//    - 根据实际需求选择合适的查询方式


}