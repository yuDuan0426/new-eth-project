package main

import (
	"context"      // 上下文管理，用于控制请求的生命周期
	"crypto/ecdsa" // 椭圆曲线数字签名算法
	"encoding/hex" // 十六进制编码解码
	"fmt"          // 格式化输入输出
	"log"          // 日志记录
	"math/big"     // 大整数运算
	"time"         // 时间处理

	"github.com/ethereum/go-ethereum"            // 以太坊核心库
	"github.com/ethereum/go-ethereum/common"     // 以太坊通用类型
	"github.com/ethereum/go-ethereum/core/types" // 以太坊核心类型
	"github.com/ethereum/go-ethereum/crypto"     // 以太坊加密工具
	"github.com/ethereum/go-ethereum/ethclient"  // 以太坊客户端
)

// 合约字节码（示例 - 一个简单的存储合约）
// 在实际使用中，您需要替换为您自己的合约字节码
const contractBytecode = "608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506102db806100606000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906101a3565b60405180910390f35b610073600480360381019061006e919061014f565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610270565b92915050565b6000602082840312156100b3576100b261026b565b600080fd5b60006100c384610088565b9050919050565b6000819050919050565b6100dd816100ca565b82525050565b60006020820190506100f860008301846100d4565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b61014081610270565b811461014b57600080fd5b50565b60006020828403121561016457610163610137565b600080fd5b600061017484610088565b9050919050565b610184816100ca565b811461018f57600080fd5b50565b6000813590506101a18161017b565b92915050565b60006020820190506101bc60008301846100d4565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806101f957607f821691505b6020821081141561020d5761020c6101c2565b50565b6000819050919050565b61022381610210565b82525050565b600060208201905061023e600083018461021a565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061027f82610210565b915061028a83610210565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156102bf576102be610244565b5b828201905092915050565b6102d381610210565b81146102de57600080fd5b5056fea2646970667358221220c7f729d1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c1a1c164736f6c63430008070033"

// main函数 - 部署智能合约到以太坊网络
// 功能：将智能合约字节码部署到以太坊网络，并返回合约地址
// 这是一个完整的合约部署实现
func main() {
	// 步骤1：连接到以太坊网络
	// 使用Goerli测试网络作为示例（您需要替换为您的Infura项目ID）
	// 注意：请将YOUR-PROJECT-ID替换为您的实际Infura项目ID
	client, err := ethclient.Dial("https://goerli.infura.io/v3/YOUR-PROJECT-ID")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("成功连接到以太坊网络")

	// 步骤2：创建私钥
	// 在实际应用中，您应该使用更安全的方式来管理私钥
	// 注意：请将YOUR-PRIVATE-KEY-HERE替换为您的实际私钥（不包含0x前缀）
	privateKey, err := crypto.HexToECDSA("YOUR-PRIVATE-KEY-HERE")
	if err != nil {
		log.Fatal("创建私钥失败:", err)
	}
	fmt.Println("私钥加载成功")

	// 步骤3：从私钥推导公钥和地址
	// 公钥是从私钥计算得出的，地址是从公钥计算得出的
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型转换失败")
	}

	// 从公钥推导出以太坊地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("部署账户地址: %s\n", fromAddress.Hex())

	// 步骤4：获取nonce
	// nonce是账户发送交易的序号，防止重放攻击
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取nonce失败:", err)
	}
	fmt.Printf("当前nonce: %d\n", nonce)

	// 步骤5：获取建议的gas价格
	// Gas价格决定了交易的优先级，价格越高越容易被矿工打包
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取gas价格失败:", err)
	}
	fmt.Printf("建议gas价格: %s wei\n", gasPrice.String())

	// 步骤6：解码合约字节码
	// 合约字节码是编译后的智能合约代码
	data, err := hex.DecodeString(contractBytecode)
	if err != nil {
		log.Fatal("解码合约字节码失败:", err)
	}
	fmt.Printf("合约字节码长度: %d bytes\n", len(data))

	// 步骤7：创建合约部署交易
	// NewContractCreation创建一个合约部署交易
	// 参数：nonce, value(发送的ETH数量), gasLimit, gasPrice, 合约字节码
	tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)
	fmt.Println("合约部署交易创建成功")

	// 步骤8：获取网络ID并签名交易
	// 网络ID用于防止跨链重放攻击
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取网络ID失败:", err)
	}
	fmt.Printf("网络ID: %s\n", chainID.String())

	// 使用EIP155签名标准签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("签名交易失败:", err)
	}
	fmt.Println("交易签名成功")

	// 步骤9：发送交易到网络
	// 将签名后的交易广播到以太坊网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败:", err)
	}

	fmt.Printf("交易已发送，交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Println("等待交易被挖矿...")

	// 步骤10：等待交易被挖矿并获取收据
	// 交易收据包含合约地址等重要信息
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal("等待交易收据失败:", err)
	}

	// 步骤11：显示部署结果
	if receipt.Status == 1 {
		fmt.Println("\n=== 合约部署成功! ===")
		fmt.Printf("合约地址: %s\n", receipt.ContractAddress.Hex())
		fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("Gas使用量: %d\n", receipt.GasUsed)
		fmt.Printf("交易哈希: %s\n", receipt.TxHash.Hex())
	} else {
		fmt.Println("合约部署失败")
	}
}

// waitForReceipt 等待交易被挖矿并返回交易收据
// 参数：client - 以太坊客户端，txHash - 交易哈希
// 返回：交易收据和可能的错误
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	fmt.Println("正在等待交易确认...")

	// 无限循环直到交易被确认
	for {
		// 尝试获取交易收据
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			// 交易已被确认，返回收据
			return receipt, nil
		}

		// 如果错误不是"未找到"，说明出现了其他问题
		if err != ethereum.NotFound {
			return nil, err
		}

		// 等待1秒后再次查询
		// 避免过于频繁的查询
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}

	// 小白说明：
	// 1. 合约部署是一种特殊的交易，to地址为空，data字段包含合约字节码
	// 2. 合约地址是根据部署者地址和nonce计算得出的
	// 3. Gas限制需要足够大以容纳合约部署，通常比普通转账需要更多Gas
	// 4. 交易收据中的Status字段：1表示成功，0表示失败
	// 5. 合约部署成功后，可以通过合约地址与合约交互
	// 6. 私钥管理非常重要，在生产环境中应使用硬件钱包或密钥管理服务
	// 7. 测试网络用于开发和测试，主网用于正式部署
	// 8. 合约字节码是Solidity编译后的结果，包含合约的所有逻辑
}
