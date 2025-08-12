package main

import (
	"context"      // 上下文管理，用于控制请求的生命周期
	"crypto/ecdsa" // 椭圆曲线数字签名算法
	"fmt"          // 格式化输入输出
	"log"          // 日志记录
	"math/big"     // 大数运算，处理代币数量等大整数

	"golang.org/x/crypto/sha3" // SHA3哈希算法

	"github.com/ethereum/go-ethereum"                // 以太坊核心库
	"github.com/ethereum/go-ethereum/common"         // 以太坊通用工具
	"github.com/ethereum/go-ethereum/common/hexutil" // 十六进制工具
	"github.com/ethereum/go-ethereum/core/types"     // 以太坊交易类型
	"github.com/ethereum/go-ethereum/crypto"         // 以太坊加密工具
	"github.com/ethereum/go-ethereum/ethclient"      // 以太坊客户端
)

// main函数 - ERC20代币转账
// 功能：发送ERC20代币从一个地址到另一个地址
// 这是一个完整的代币转账实现，包括智能合约调用
func main() {
	// 步骤1：连接以太坊网络
	// 使用Alchemy提供的Sepolia测试网络节点
	// 注意：需要将URL中的API_KEY替换为您的实际密钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal("连接以太坊客户端失败:", err)
	}

	// 步骤2：加载发送方私钥
	// 注意：实际使用时需要将"账户私钥"替换为真实的私钥
	// 私钥格式：64位十六进制字符串（不包含0x前缀）
	privateKey, err := crypto.HexToECDSA("账户私钥")
	if err != nil {
		log.Fatal("私钥解析失败:", err)
	}

	// 步骤3：从私钥推导公钥和地址
	// 获取公钥
	publicKey := privateKey.Public()
	// 将公钥转换为ECDSA公钥类型
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型转换失败: 无法转换为*ecdsa.PublicKey类型")
	}

	// 从公钥推导出发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("发送方地址: %s\n", fromAddress.Hex())

	// 步骤4：获取交易nonce
	// nonce是账户的交易序号，防止重复交易
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取nonce失败:", err)
	}
	fmt.Printf("当前nonce: %d\n", nonce)

	// 步骤5：设置交易参数
	// 代币转账不需要发送ETH，所以value设为0
	value := big.NewInt(0) // 0 ETH（因为我们转账的是代币，不是ETH）

	// 获取当前网络建议的Gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取Gas价格失败:", err)
	}
	fmt.Printf("Gas价格: %s wei\n", gasPrice.String())

	// 步骤6：设置转账目标和代币合约地址
	// 接收代币的地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	// ERC20代币合约地址
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	fmt.Printf("接收地址: %s\n", toAddress.Hex())
	fmt.Printf("代币合约地址: %s\n", tokenAddress.Hex())

	// 步骤7：构建智能合约调用数据
	// ERC20的transfer函数签名：transfer(address,uint256)
	transferFnSignature := []byte("transfer(address,uint256)")

	// 计算函数选择器（前4字节的Keccak256哈希）
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("函数选择器: %s\n", hexutil.Encode(methodID)) // 0xa9059cbb

	// 将接收地址填充为32字节（左填充0）
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Printf("填充后的地址: %s\n", hexutil.Encode(paddedAddress))

	// 设置转账数量（这里是1000个代币，假设代币有18位小数）
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 * 10^18 = 1000个代币
	// 将数量填充为32字节
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("填充后的数量: %s\n", hexutil.Encode(paddedAmount))

	// 组装完整的调用数据：函数选择器 + 参数
	var data []byte
	data = append(data, methodID...)      // 函数选择器
	data = append(data, paddedAddress...) // 接收地址参数
	data = append(data, paddedAmount...)  // 转账数量参数

	// 步骤8：估算Gas消耗
	// 注意：这里To地址应该是代币合约地址，不是接收方地址
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress, // 修正：应该是代币合约地址
		Data: data,
	})
	if err != nil {
		log.Fatal("估算Gas失败:", err)
	}
	fmt.Printf("估算Gas消耗: %d\n", gasLimit)

	// 步骤9：创建交易
	// 注意：To地址是代币合约地址，不是接收方地址
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// 步骤10：获取网络ID并签名交易
	// 获取链ID（网络标识符）
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取网络ID失败:", err)
	}
	fmt.Printf("网络ID: %s\n", chainID.String())

	// 使用EIP155签名器签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("交易签名失败:", err)
	}

	// 步骤11：广播交易到网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("发送交易失败:", err)
	}

	// 输出交易结果
	fmt.Printf("\n=== 交易发送成功 ===\n")
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("接收方: %s\n", toAddress.Hex())
	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("转账数量: %s (最小单位)\n", amount.String())

	// 小白说明：
	// 1. ERC20代币转账实际上是调用智能合约的transfer函数
	// 2. 函数选择器是函数签名的Keccak256哈希的前4字节
	// 3. 智能合约的参数需要按照ABI规范进行编码（32字节对齐）
	// 4. 交易的To地址是代币合约地址，不是接收方地址
	// 5. value为0因为我们转账的是代币，不是ETH
	// 6. Gas费用仍然用ETH支付，即使转账的是代币
}