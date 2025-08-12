// 以太坊交易回执查询工具
// 本程序演示如何查询以太坊交易回执信息
// 包含批量查询区块回执和单个交易回执查询两种方式

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	fmt.Println("=== 以太坊交易回执查询工具 ===")
	fmt.Println("本工具演示如何查询交易回执信息，包括批量和单个查询\n")

	// ===== 第1步：连接以太坊网络 =====
	// 连接到以太坊Sepolia测试网
	// 注意：需要将<API_KEY>替换为实际的Alchemy或Infura API密钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/<API_KEY>")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✓ 成功连接到以太坊Sepolia测试网")

	// ===== 第2步：设置查询参数 =====
	// 指定要查询的区块号和区块哈希
	blockNumber := big.NewInt(5671744)
	blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	fmt.Printf("查询区块号: %s\n", blockNumber.String())
	fmt.Printf("查询区块哈希: %s\n\n", blockHash.Hex())

	// ===== 第3步：按区块哈希批量查询交易回执 =====
	fmt.Println("=== 按区块哈希批量查询交易回执 ===")
	// BlockReceipts方法可以一次性获取整个区块中所有交易的回执
	// 这比逐个查询交易回执更高效
	receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		log.Fatal("按区块哈希查询回执失败:", err)
	}
	fmt.Printf("✓ 通过区块哈希获取到 %d 个交易回执\n", len(receiptByHash))

	// ===== 第4步：按区块号批量查询交易回执 =====
	fmt.Println("\n=== 按区块号批量查询交易回执 ===")
	// 同样的功能，但使用区块号而不是区块哈希
	receiptsByNum, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	if err != nil {
		log.Fatal("按区块号查询回执失败:", err)
	}
	fmt.Printf("✓ 通过区块号获取到 %d 个交易回执\n", len(receiptsByNum))

	// ===== 第5步：验证两种查询方式的结果一致性 =====
	fmt.Println("\n=== 验证查询结果一致性 ===")
	// 比较两种方式获取的第一个回执是否相同
	if len(receiptByHash) > 0 && len(receiptsByNum) > 0 {
		isEqual := receiptByHash[0].TxHash == receiptsByNum[0].TxHash &&
			receiptByHash[0].Status == receiptsByNum[0].Status &&
			receiptByHash[0].TransactionIndex == receiptsByNum[0].TransactionIndex
		fmt.Printf("两种查询方式结果一致: %t\n", isEqual)
	} else {
		fmt.Println("无法比较：其中一种方式未返回回执")
	}

	// ===== 第6步：分析第一个交易回执的详细信息 =====
	fmt.Println("\n=== 分析第一个交易回执详情 ===")
	if len(receiptByHash) > 0 {
		firstReceipt := receiptByHash[0]
		
		// 交易执行状态
		if firstReceipt.Status == 1 {
			fmt.Println("交易状态: 成功执行")
		} else {
			fmt.Println("交易状态: 执行失败")
		}
		
		// 事件日志信息
		fmt.Printf("事件日志数量: %d\n", len(firstReceipt.Logs))
		if len(firstReceipt.Logs) == 0 {
			fmt.Println("事件日志: 无（普通转账交易）")
		} else {
			fmt.Println("事件日志: 包含智能合约事件")
		}
		
		// 交易基本信息
		fmt.Printf("交易哈希: %s\n", firstReceipt.TxHash.Hex())
		fmt.Printf("交易在区块中的索引: %d\n", firstReceipt.TransactionIndex)
		fmt.Printf("Gas消耗量: %d\n", firstReceipt.GasUsed)
		fmt.Printf("累计Gas消耗: %d\n", firstReceipt.CumulativeGasUsed)
		
		// 合约地址信息
		if firstReceipt.ContractAddress == (common.Address{}) {
			fmt.Println("合约地址: 无（非合约创建交易）")
		} else {
			fmt.Printf("合约地址: %s（合约创建交易）\n", firstReceipt.ContractAddress.Hex())
		}
		
		// 区块信息
		fmt.Printf("所在区块哈希: %s\n", firstReceipt.BlockHash.Hex())
		fmt.Printf("所在区块号: %d\n", firstReceipt.BlockNumber.Uint64())
	} else {
		fmt.Println("该区块中没有交易")
	}

	// ===== 第7步：单独查询指定交易的回执 =====
	fmt.Println("\n=== 单独查询指定交易回执 ===")
	// 指定要查询的交易哈希
	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
	fmt.Printf("查询交易哈希: %s\n", txHash.Hex())
	
	// 使用TransactionReceipt方法查询单个交易的回执
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal("查询单个交易回执失败:", err)
	}
	fmt.Println("✓ 成功获取交易回执")

	// ===== 第8步：显示单个交易回执的详细信息 =====
	fmt.Println("\n=== 单个交易回执详细信息 ===")
	
	// 交易执行状态
	if receipt.Status == 1 {
		fmt.Println("交易状态: 成功执行")
	} else {
		fmt.Println("交易状态: 执行失败")
	}
	
	// 事件日志
	fmt.Printf("事件日志数量: %d\n", len(receipt.Logs))
	if len(receipt.Logs) == 0 {
		fmt.Println("事件日志: 无（普通转账交易）")
	} else {
		fmt.Println("事件日志详情:")
		for i, eventLog := range receipt.Logs {
			fmt.Printf("  日志 #%d: 地址=%s, 主题数量=%d\n", i+1, eventLog.Address.Hex(), len(eventLog.Topics))
		}
	}
	
	// 交易标识信息
	fmt.Printf("交易哈希: %s\n", receipt.TxHash.Hex())
	fmt.Printf("交易索引: %d\n", receipt.TransactionIndex)
	
	// Gas相关信息
	fmt.Printf("Gas消耗量: %d\n", receipt.GasUsed)
	fmt.Printf("累计Gas消耗: %d\n", receipt.CumulativeGasUsed)
	fmt.Printf("有效Gas价格: %s Wei\n", receipt.EffectiveGasPrice.String())
	
	// 合约相关信息
	if receipt.ContractAddress == (common.Address{}) {
		fmt.Println("合约地址: 无（非合约创建交易）")
	} else {
		fmt.Printf("新创建的合约地址: %s\n", receipt.ContractAddress.Hex())
	}
	
	// 区块位置信息
	fmt.Printf("所在区块哈希: %s\n", receipt.BlockHash.Hex())
	fmt.Printf("所在区块号: %d\n", receipt.BlockNumber.Uint64())
	
	// 交易类型（EIP-2718）
	fmt.Printf("交易类型: %d\n", receipt.Type)

	fmt.Printf("\n=== 验证信息 ===\n")
	fmt.Println("✓ 批量回执查询完成")
	fmt.Println("✓ 单个回执查询完成")
	fmt.Println("✓ 回执信息解析完成")

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 替换<API_KEY>为实际的Alchemy或Infura密钥")
	fmt.Println("2. 调整区块号、区块哈希和交易哈希为要查询的实际值")
	fmt.Println("3. 批量查询适合分析整个区块的交易执行情况")
	fmt.Println("4. 单个查询适合验证特定交易的执行结果")
	fmt.Println("5. 回执中的Status字段是判断交易成功与否的关键")

	// ===== 技术说明 =====
	// 1. 交易回执(Receipt)包含的关键信息：
	//    - Status: 交易执行状态（1=成功，0=失败）
	//    - GasUsed: 实际消耗的Gas数量
	//    - Logs: 智能合约产生的事件日志
	//    - ContractAddress: 新创建的合约地址（仅合约创建交易）
	//
	// 2. 批量查询 vs 单个查询：
	//    - BlockReceipts: 一次获取整个区块的所有回执，效率高
	//    - TransactionReceipt: 查询单个交易回执，精确定位
	//
	// 3. Gas相关字段：
	//    - GasUsed: 该交易实际消耗的Gas
	//    - CumulativeGasUsed: 到该交易为止区块中累计消耗的Gas
	//    - EffectiveGasPrice: 实际生效的Gas价格（EIP-1559后引入）
	//
	// 4. 事件日志：
	//    - 智能合约执行时产生的事件记录
	//    - 包含事件的地址、主题(Topics)和数据
	//    - 普通转账交易通常没有事件日志
	//
	// 5. 交易类型：
	//    - Type 0: Legacy交易
	//    - Type 1: EIP-2930 (访问列表交易)
	//    - Type 2: EIP-1559 (动态费用交易)
}