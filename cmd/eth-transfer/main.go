// 以太坊转账工具
// 本程序演示如何在以太坊网络上进行ETH转账
// 包含完整的交易创建、签名和发送流程

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 以太坊转账工具 ===")
	fmt.Println("本工具演示如何进行ETH转账，包括交易创建、签名和发送\n")

	// ===== 第1步：连接以太坊网络 =====
	// 注意：这里使用的是Rinkeby测试网（已废弃），建议改为Sepolia测试网
	// 生产环境请使用主网端点
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✓ 成功连接到以太坊网络")

	// ===== 第2步：加载私钥 =====
	// 从十六进制字符串加载私钥
	// 警告：这是示例私钥，实际使用时请使用安全的私钥管理方式
	// 生产环境中绝不要在代码中硬编码私钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal("加载私钥失败:", err)
	}
	fmt.Println("✓ 成功加载私钥")

	// ===== 第3步：从私钥推导公钥和地址 =====
	// 获取公钥
	publicKey := privateKey.Public()
	// 将公钥转换为ECDSA公钥类型
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型转换失败: 无法转换为*ecdsa.PublicKey类型")
	}

	// 从公钥推导以太坊地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("✓ 发送方地址: %s\n", fromAddress.Hex())

	// ===== 第4步：获取账户nonce =====
	// nonce是账户发送的交易序号，用于防止重放攻击
	// PendingNonceAt获取包含待处理交易的最新nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取nonce失败:", err)
	}
	fmt.Printf("✓ 当前nonce: %d\n", nonce)

	// ===== 第5步：设置交易参数 =====
	// 转账金额：1 ETH = 10^18 Wei
	value := big.NewInt(1000000000000000000) // 1 ETH in wei
	fmt.Printf("转账金额: %s Wei (1 ETH)\n", value.String())

	// Gas限制：ETH转账的标准Gas限制是21000
	gasLimit := uint64(21000)
	fmt.Printf("Gas限制: %d\n", gasLimit)

	// 获取建议的Gas价格
	// Gas价格决定了交易的优先级，价格越高越容易被矿工打包
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取Gas价格失败:", err)
	}
	fmt.Printf("建议Gas价格: %s Wei\n", gasPrice.String())

	// 接收方地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	fmt.Printf("✓ 接收方地址: %s\n\n", toAddress.Hex())

	// ===== 第6步：创建交易 =====
	// 对于ETH转账，data字段为空
	var data []byte
	// 创建新的交易对象
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	fmt.Println("✓ 交易对象创建成功")

	// ===== 第7步：获取网络ID并签名交易 =====
	// 获取网络ID（链ID），用于EIP155签名防止重放攻击
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取网络ID失败:", err)
	}
	fmt.Printf("✓ 网络ID: %s\n", chainID.String())

	// 使用EIP155签名器对交易进行签名
	// EIP155是以太坊改进提案，用于防止跨链重放攻击
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("交易签名失败:", err)
	}
	fmt.Println("✓ 交易签名成功")

	// ===== 第8步：发送交易到网络 =====
	// 将签名后的交易广播到以太坊网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败:", err)
	}

	// ===== 第9步：显示交易结果 =====
	fmt.Println("\n=== 交易发送成功 ===")
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Println("\n=== 交易详情 ===")
	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("接收方: %s\n", toAddress.Hex())
	fmt.Printf("金额: %s Wei (1 ETH)\n", value.String())
	fmt.Printf("Gas限制: %d\n", gasLimit)
	fmt.Printf("Gas价格: %s Wei\n", gasPrice.String())
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("网络ID: %s\n", chainID.String())

	fmt.Println("\n=== 重要提醒 ===")
	fmt.Println("1. 这是示例代码，使用的是测试网络")
	fmt.Println("2. 实际使用时请替换为您自己的私钥和接收地址")
	fmt.Println("3. 生产环境中请使用安全的私钥管理方式")
	fmt.Println("4. 建议先在测试网验证后再在主网使用")
	fmt.Println("5. 请确保发送方账户有足够的ETH支付Gas费用")

	// ===== 技术说明 =====
	// 1. 以太坊转账流程：
	//    - 连接网络 → 加载私钥 → 获取nonce → 设置参数 → 创建交易 → 签名 → 发送
	//
	// 2. 关键概念：
	//    - Wei: 以太坊最小单位，1 ETH = 10^18 Wei
	//    - Nonce: 账户交易序号，防止重放攻击
	//    - Gas: 计算资源的度量单位
	//    - Gas Price: 每单位Gas的价格，影响交易优先级
	//
	// 3. 安全注意事项：
	//    - 私钥管理：使用硬件钱包或安全的密钥管理服务
	//    - 网络选择：测试网用于开发，主网用于生产
	//    - Gas设置：过低可能导致交易失败，过高浪费资金
	//    - 地址验证：确保接收地址正确，交易不可逆
	//
	// 4. EIP155签名：
	//    - 包含链ID的签名方案，防止跨链重放攻击
	//    - 确保交易只能在指定的区块链网络上执行
	//
	// 5. 交易状态：
	//    - 发送成功只表示交易进入内存池
	//    - 需要等待矿工打包确认才算真正完成
	//    - 可以通过交易哈希查询确认状态
}