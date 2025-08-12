package main

import (
	"context" // 上下文管理，用于控制请求的生命周期
	"fmt"     // 格式化输入输出
	"log"     // 日志记录

	"github.com/ethereum/go-ethereum/core/types" // 以太坊核心类型
	"github.com/ethereum/go-ethereum/ethclient"  // 以太坊客户端
)

// main函数 - 订阅以太坊新区块
// 功能：实时订阅以太坊网络的新区块，并显示区块详细信息
// 这是一个完整的区块订阅实现，使用WebSocket连接
func main() {
	// 步骤1：连接以太坊网络
	// 使用Infura提供的Ropsten测试网络WebSocket节点
	// 注意：Ropsten测试网已经停用，建议使用Sepolia或Goerli测试网
	// 您需要将URL替换为有效的WebSocket端点
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		log.Fatal("连接以太坊客户端失败:", err)
	}
	fmt.Println("成功连接到以太坊网络，开始订阅新区块...")

	// 步骤2：创建区块头通道
	// 这个通道将接收新区块的头部信息
	headers := make(chan *types.Header)

	// 步骤3：订阅新区块头
	// SubscribeNewHead会在有新区块产生时发送区块头到通道
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("订阅新区块失败:", err)
	}
	fmt.Println("订阅成功！等待新区块...")

	// 步骤4：无限循环监听新区块
	for {
		select {
		// 处理订阅错误
		case err := <-sub.Err():
			log.Fatal("订阅出现错误:", err)

		// 处理新区块头
		case header := <-headers:
			// 显示区块哈希
			fmt.Printf("\n=== 新区块到达 ===\n")
			fmt.Printf("区块哈希: %s\n", header.Hash().Hex())

			// 步骤5：根据区块哈希获取完整区块信息
			// 区块头只包含基本信息，要获取交易列表需要查询完整区块
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal("获取完整区块信息失败:", err)
			}

			// 步骤6：显示区块详细信息
			fmt.Printf("区块哈希: %s\n", block.Hash().Hex())        // 区块的唯一标识符
			fmt.Printf("区块高度: %d\n", block.Number().Uint64())   // 区块在链上的序号
			fmt.Printf("时间戳: %d\n", block.Time())               // 区块创建时间（Unix时间戳）
			fmt.Printf("随机数: %d\n", block.Nonce())             // 挖矿时使用的随机数
			fmt.Printf("交易数量: %d\n", len(block.Transactions())) // 区块中包含的交易数量

			// 可选：显示更多区块信息
			fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex()) // 前一个区块的哈希
			fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())    // 挖出这个区块的矿工地址
			fmt.Printf("Gas使用量: %d\n", block.GasUsed())          // 区块中所有交易消耗的Gas总量
			fmt.Printf("Gas限制: %d\n", block.GasLimit())          // 区块的Gas限制
			fmt.Println("---")
		}
	}

	// 小白说明：
	// 1. WebSocket连接允许实时接收数据，比HTTP轮询更高效
	// 2. 区块头包含区块的基本信息，但不包含交易详情
	// 3. 要获取完整区块信息（包括交易），需要额外调用BlockByHash
	// 4. select语句用于同时监听多个通道，这里监听错误和新区块
	// 5. 区块哈希是区块的唯一标识符，由区块内容计算得出
	// 6. 区块高度表示区块在区块链中的位置，从0开始递增
	// 7. 时间戳记录了区块被挖出的时间
	// 8. Nonce是矿工在挖矿过程中尝试的随机数
	// 9. 程序会一直运行直到手动停止（Ctrl+C）
}