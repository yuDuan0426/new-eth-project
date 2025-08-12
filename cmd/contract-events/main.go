package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// StoreABI 智能合约的ABI定义
// ABI（Application Binary Interface）描述了合约的接口
// 包含了合约的所有函数、事件和数据结构的定义
var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	// ===== 第1步：连接以太坊网络 =====
	// 使用WebSocket连接，支持实时事件推送
	// 注意：这里使用的是Rinkeby测试网，现在已经废弃
	// 建议替换为Sepolia测试网或主网的WebSocket端点
	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✅ 成功连接到以太坊网络（WebSocket）")

	// ===== 第2步：设置要监听的合约地址 =====
	// 这是要监听事件的智能合约地址
	// 你需要替换为实际部署的合约地址
	contractAddress := common.HexToAddress("0x2958d15bc5b64b11Ec65e623Ac50C198519f8742")
	fmt.Printf("监听合约地址: %s\n", contractAddress.Hex())

	// ===== 第3步：创建事件过滤器 =====
	// FilterQuery定义了要监听哪些事件
	// Addresses: 指定要监听的合约地址列表
	// Topics: 可以指定要监听的特定事件类型（可选）
	// FromBlock/ToBlock: 可以指定监听的区块范围（可选）
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress}, // 只监听指定合约的事件
	}
	fmt.Println("📡 创建事件过滤器")

	// ===== 第4步：创建事件订阅 =====
	// logs通道用于接收事件日志
	// SubscribeFilterLogs创建一个实时事件订阅
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("创建事件订阅失败:", err)
	}
	fmt.Println("🔔 开始监听合约事件...")

	// ===== 第5步：解析合约ABI =====
	// 将ABI字符串解析为Go可以使用的ABI对象
	// 这个对象用于解码事件数据
	contractAbi, err := abi.JSON(strings.NewReader(string(StoreABI)))
	if err != nil {
		log.Fatal("解析合约ABI失败:", err)
	}

	// ===== 第6步：事件监听主循环 =====
	// 使用select语句同时监听多个通道
	for {
		select {
		// 监听订阅错误
		case err := <-sub.Err():
			log.Fatal("事件订阅出错:", err)

		// 监听新的事件日志
		case vLog := <-logs:
			fmt.Println("\n🎉 收到新事件！")
			fmt.Println("===========================================")

			// ===== 第7步：显示事件基本信息 =====
			// 区块哈希：包含此事件的区块的唯一标识
			fmt.Printf("📦 区块哈希: %s\n", vLog.BlockHash.Hex())
			// 区块号：事件发生的区块编号
			fmt.Printf("🔢 区块号: %d\n", vLog.BlockNumber)
			// 交易哈希：触发此事件的交易的唯一标识
			fmt.Printf("💳 交易哈希: %s\n", vLog.TxHash.Hex())

			// ===== 第8步：解析事件数据 =====
			// 定义事件数据结构，必须与合约中的事件定义匹配
			// ItemSet事件包含key和value两个字段
			event := struct {
				Key   [32]byte // bytes32类型的key
				Value [32]byte // bytes32类型的value
			}{}

			// 使用ABI解码事件数据
			// UnpackIntoInterface将原始事件数据解码到结构体中
			// "ItemSet"是事件名称，vLog.Data包含非indexed参数的数据
			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Fatal("解析事件数据失败:", err)
			}

			// ===== 第9步：显示解析后的事件数据 =====
			// 将bytes32数据转换为十六进制字符串显示
			fmt.Printf("🔑 Key (hex): %s\n", common.Bytes2Hex(event.Key[:]))
			fmt.Printf("💎 Value (hex): %s\n", common.Bytes2Hex(event.Value[:]))

			// 尝试将bytes32数据转换为可读字符串（如果是文本数据）
			// 移除末尾的零字节
			keyStr := strings.TrimRight(string(event.Key[:]), "\x00")
			valueStr := strings.TrimRight(string(event.Value[:]), "\x00")
			if keyStr != "" {
				fmt.Printf("🔑 Key (string): %s\n", keyStr)
			}
			if valueStr != "" {
				fmt.Printf("💎 Value (string): %s\n", valueStr)
			}

			// ===== 第10步：处理事件主题(Topics) =====
			// Topics包含indexed参数和事件签名
			// Topics[0]是事件签名的哈希
			// Topics[1:]是indexed参数的值
			var topics []string
			for i := range vLog.Topics {
				topics = append(topics, vLog.Topics[i].Hex())
			}

			// 显示事件签名（第一个topic）
			fmt.Printf("📋 事件签名: %s\n", topics[0])

			// 显示indexed参数（如果有的话）
			// 在ItemSet事件中，key是indexed参数，所以会出现在topics[1]中
			if len(topics) > 1 {
				fmt.Printf("🏷️  Indexed参数: %v\n", topics[1:])
				// 第一个indexed参数就是key值
				fmt.Printf("🔑 Key (from topic): %s\n", topics[1])
			}

			fmt.Println("===========================================")
		}
	}

	// 注意：这个程序会一直运行，监听新的事件
	// 在实际应用中，你可能需要添加优雅关闭的逻辑
	// 例如监听系统信号，在收到关闭信号时清理资源并退出

	// ===== 事件监听的重要概念说明 =====
	// 1. Indexed vs Non-indexed参数：
	//    - Indexed参数：存储在topics中，可以用于快速搜索和过滤
	//    - Non-indexed参数：存储在data中，包含详细信息但不能用于搜索
	//
	// 2. 事件签名：
	//    - Topics[0]包含事件签名的Keccak256哈希
	//    - 事件签名格式：EventName(type1,type2,...)
	//    - 例如：ItemSet(bytes32,bytes32)
	//
	// 3. WebSocket vs HTTP：
	//    - WebSocket支持实时推送，适合事件监听
	//    - HTTP需要轮询，延迟较高且消耗资源
	//
	// 4. 错误处理：
	//    - 网络断线时订阅会出错，需要重新连接
	//    - 建议实现自动重连机制
	//
	// 5. 性能考虑：
	//    - 可以通过Topics过滤特定事件类型
	//    - 可以设置区块范围避免处理过多历史数据
}