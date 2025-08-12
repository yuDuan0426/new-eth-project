// 历史事件查询示例
// 本文件演示如何查询智能合约的历史事件
// 与contract-events目录中的实时事件监听不同，这里是一次性查询历史数据
//
// 使用方法：
// 1. 在此目录下运行: go run main.go
// 2. 或者构建后运行: go build && ./query-history-events

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// StoreABI 智能合约的ABI定义
// 定义了合约的接口，包含构造函数、事件和函数的完整定义
// 这个ABI描述了ItemSet事件的结构：key(indexed), value(non-indexed)
var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

// main函数 - 查询历史合约事件
// 功能：查询指定区块范围内的历史事件，而不是实时监听
// 适用场景：数据分析、历史记录查询、事件回放等
func main() {
	fmt.Println("=== 智能合约历史事件查询工具 ===")
	fmt.Println("本工具用于查询指定区块范围内的合约事件")
	fmt.Println("与实时事件监听不同，这是一次性的历史数据查询\n")

	// ===== 第1步：连接以太坊网络 =====
	// 使用HTTP连接，适合一次性查询操作
	// 对于历史数据查询，HTTP连接比WebSocket更稳定
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/<API_KEY>")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✓ 成功连接到以太坊Sepolia测试网")

	// ===== 第2步：设置查询参数 =====
	// 指定要查询的合约地址
	contractAddress := common.HexToAddress("0x2958d15bc5b64b11Ec65e623Ac50C198519f8742")
	fmt.Printf("✓ 目标合约地址: %s\n", contractAddress.Hex())

	// ===== 第3步：创建事件过滤器 =====
	// FilterQuery定义了查询条件
	query := ethereum.FilterQuery{
		// FromBlock: 查询的起始区块号
		// 注意：区块范围不要太大，避免查询超时
		FromBlock: big.NewInt(6920583),
		// ToBlock: 查询的结束区块号（注释掉表示查询到最新区块）
		// ToBlock:   big.NewInt(2394201),
		// Addresses: 要监听的合约地址列表
		Addresses: []common.Address{
			contractAddress,
		},
		// Topics: 事件主题过滤器（注释掉表示查询所有事件）
		// Topics可以用来过滤特定的事件类型或参数
		// Topics: [][]common.Hash{
		//  {}, // 第一个元素是事件签名哈希
		//  {}, // 后续元素是indexed参数的值
		// },
	}
	fmt.Printf("✓ 查询区块范围: 从 %d 到最新区块\n", query.FromBlock.Int64())

	// ===== 第4步：执行日志查询 =====
	// FilterLogs返回符合条件的所有历史日志
	// 这是一次性操作，不同于实时订阅
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal("查询事件日志失败:", err)
	}
	fmt.Printf("✓ 找到 %d 个事件日志\n\n", len(logs))

	// ===== 第5步：解析合约ABI =====
	// 将ABI字符串解析为Go可以使用的ABI对象
	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal("解析合约ABI失败:", err)
	}
	fmt.Println("✓ 成功解析合约ABI")

	// ===== 第6步：遍历和解析每个事件日志 =====
	for i, vLog := range logs {
		fmt.Printf("\n=== 事件 #%d ===\n", i+1)
		
		// 显示事件的基本信息
		fmt.Printf("区块哈希: %s\n", vLog.BlockHash.Hex())
		fmt.Printf("区块号: %d\n", vLog.BlockNumber)
		fmt.Printf("交易哈希: %s\n", vLog.TxHash.Hex())
		fmt.Printf("日志索引: %d\n", vLog.Index)

		// ===== 第7步：解析事件数据 =====
		// 定义与ItemSet事件匹配的结构体
		// 注意：这里只包含non-indexed参数（value）
		// indexed参数（key）在Topics中处理
		event := struct {
			Key   [32]byte // 这个字段实际不会被UnpackIntoInterface填充
			Value [32]byte // 这个字段会被填充
		}{}
		
		// UnpackIntoInterface解析事件的Data字段
		// 只解析non-indexed参数
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Printf("解析事件数据失败: %v", err)
			continue
		}

		// 显示解析后的事件数据
		fmt.Printf("Value (non-indexed): 0x%s\n", common.Bytes2Hex(event.Value[:]))

		// ===== 第8步：处理事件Topics =====
		// Topics[0]: 事件签名的Keccak256哈希
		// Topics[1:]: indexed参数的值
		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}

		// 显示事件签名
		fmt.Printf("事件签名哈希: %s\n", topics[0])
		
		// 显示indexed参数
		if len(topics) > 1 {
			fmt.Printf("Key (indexed): %s\n", topics[1])
			fmt.Printf("所有indexed参数: %v\n", topics[1:])
		}
	}

	// ===== 第9步：验证事件签名 =====
	// 计算ItemSet事件的签名哈希，用于验证
	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Printf("\n=== 事件签名验证 ===\n")
	fmt.Printf("计算得到的事件签名哈希: %s\n", hash.Hex())
	fmt.Println("此哈希应该与所有事件的Topics[0]匹配")

	fmt.Println("\n=== 查询完成 ===")
	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 替换<API_KEY>为实际的Alchemy或Infura密钥")
	fmt.Println("2. 调整合约地址为要查询的实际合约")
	fmt.Println("3. 根据需要调整查询的区块范围（FromBlock/ToBlock）")
	fmt.Println("4. 如需查询其他事件，请更新ABI和事件结构")
	fmt.Println("5. 可以使用Topics过滤器来查询特定的事件或参数值")

	// ===== 技术说明 =====
	// 1. 历史查询 vs 实时监听：
	//    - 历史查询：使用FilterLogs一次性获取过去的事件
	//    - 实时监听：使用SubscribeFilterLogs持续监听新事件
	//
	// 2. 事件数据结构：
	//    - Topics[0]：事件签名的Keccak256哈希
	//    - Topics[1:]：indexed参数的值
	//    - Data：non-indexed参数的ABI编码数据
	//
	// 3. ABI解析：
	//    - UnpackIntoInterface只解析Data字段
	//    - indexed参数需要从Topics中手动提取
	//
	// 4. 查询优化：
	//    - 使用合适的区块范围避免超时
	//    - 使用Topics过滤器减少不必要的数据
	//    - 考虑分批查询大范围的历史数据
	//
	// 5. 性能考虑：
	//    - 历史查询：一次性查询大量数据，注意区块范围不要太大
	//    - 实时监听：持续运行，注意内存管理和错误处理
}