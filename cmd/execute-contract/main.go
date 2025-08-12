package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 合约地址常量
// 这是一个示例合约地址，你需要替换为实际部署的合约地址
const (
	contractAddr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
)

func main() {
	// ===== 第1步：连接以太坊网络 =====
	// 连接到以太坊节点，这里需要替换为实际的节点URL
	// 可以使用Infura、Alchemy等服务提供的节点
	client, err := ethclient.Dial("<execution-layer-endpoint-url>")
	if err != nil {
		log.Fatal("连接以太坊网络失败:", err)
	}
	fmt.Println("✅ 成功连接到以太坊网络")

	// ===== 第2步：加载私钥 =====
	// 从十六进制字符串加载私钥（用于签名交易）
	// 注意：私钥不要包含"0x"前缀，且要妥善保管，不要泄露
	privateKey, err := crypto.HexToECDSA("<your private key>")
	if err != nil {
		log.Fatal("加载私钥失败:", err)
	}

	// ===== 第3步：获取公钥地址 =====
	// 从私钥推导出公钥，再从公钥推导出以太坊地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("转换公钥类型失败")
	}
	// 从公钥生成以太坊地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("发送方地址: %s\n", fromAddress.Hex())

	// ===== 第4步：获取账户nonce =====
	// nonce是账户发送交易的序号，防止重放攻击
	// PendingNonceAt获取包含pending交易的最新nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取nonce失败:", err)
	}
	fmt.Printf("当前nonce: %d\n", nonce)

	// ===== 第5步：估算Gas价格 =====
	// Gas价格决定了交易的优先级，价格越高越容易被矿工打包
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取Gas价格失败:", err)
	}
	fmt.Printf("建议Gas价格: %s wei\n", gasPrice.String())

	// ===== 第6步：准备合约ABI和交易数据 =====
	// ABI（Application Binary Interface）定义了如何与合约交互
	// 这里直接在代码中定义ABI，实际项目中通常从文件加载
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatal("解析合约ABI失败:", err)
	}

	// ===== 第7步：准备调用数据 =====
	// 定义要调用的合约方法名
	methodName := "setItem"
	// 准备方法参数：key和value都是bytes32类型
	var key [32]byte
	var value [32]byte

	// 将字符串转换为bytes32格式
	// copy函数会将字符串的字节复制到固定长度的数组中
	copy(key[:], []byte("demo_save_key_use_abi"))
	copy(value[:], []byte("demo_save_value_use_abi_11111"))

	// 使用ABI编码方法调用数据
	// Pack方法将方法名和参数编码为交易的input数据
	input, err := contractABI.Pack(methodName, key, value)
	if err != nil {
		log.Fatal("编码交易数据失败:", err)
	}

	// ===== 第8步：创建和签名交易 =====
	// 设置链ID（Sepolia测试网的链ID是11155111）
	chainID := big.NewInt(int64(11155111))

	// 创建交易对象
	// 参数：nonce, 目标地址, 转账金额(0表示不转ETH), Gas限制, Gas价格, 交易数据
	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddr), big.NewInt(0), 300000, gasPrice, input)

	// 使用EIP155签名方法签名交易（防止重放攻击）
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("签名交易失败:", err)
	}

	// ===== 第9步：发送交易 =====
	// 将签名后的交易广播到网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败:", err)
	}
	fmt.Printf("✅ 交易已发送，交易哈希: %s\n", signedTx.Hash().Hex())

	// ===== 第10步：等待交易确认 =====
	// 等待交易被矿工打包并获取交易收据
	fmt.Println("⏳ 等待交易确认...")
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal("等待交易确认失败:", err)
	}
	fmt.Println("✅ 交易已确认")

	// ===== 第11步：查询合约状态（只读调用） =====
	// 调用合约的items方法查询刚刚设置的值
	fmt.Println("🔍 查询刚刚设置的值...")

	// 编码查询调用的数据
	callInput, err := contractABI.Pack("items", key)
	if err != nil {
		log.Fatal("编码查询数据失败:", err)
	}

	// 创建调用消息（只读调用，不需要Gas费用）
	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	// ===== 第12步：执行只读调用并解析结果 =====
	// CallContract执行只读调用，不会改变区块链状态
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal("调用合约失败:", err)
	}

	// 解析返回的数据
	var unpacked [32]byte
	err = contractABI.UnpackIntoInterface(&unpacked, "items", result)
	if err != nil {
		log.Fatal("解析返回数据失败:", err)
	}

	// ===== 第13步：验证结果 =====
	// 比较查询到的值是否与设置的值相同
	isEqual := unpacked == value
	fmt.Printf("📊 查询结果验证: %t\n", isEqual)
	if isEqual {
		fmt.Println("🎉 合约调用成功！设置的值已正确保存")
	} else {
		fmt.Println("❌ 验证失败，查询到的值与设置的值不匹配")
	}

	fmt.Println("\n===== 程序执行完成 =====")
}


// waitForReceipt 等待交易收据的辅助函数
// 这个函数会持续查询交易状态，直到交易被确认或出现错误
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		// 尝试获取交易收据
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			// 成功获取收据，交易已确认
			return receipt, nil
		}
		if err != ethereum.NotFound {
			// 出现其他错误（非"未找到"错误）
			return nil, err
		}
		// 交易还未被确认，等待1秒后重试
		time.Sleep(1 * time.Second)
	}
}