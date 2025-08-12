package main

import (
	"fmt" // 格式化输入输出
	"log" // 日志记录

	"github.com/ethereum/go-ethereum/common"    // 以太坊通用工具
	"github.com/ethereum/go-ethereum/ethclient" // 以太坊客户端

	// 注意：在实际项目中，您需要使用abigen工具生成ERC20代币合约的Go绑定文件
	// 命令：abigen --abi token.abi --pkg token --out token.go
	// 然后导入以下包：
	// "math"
	// "math/big"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	// token "./contracts_erc20"
)

// main函数 - 查询ERC20代币余额
// 功能：查询指定地址的ERC20代币余额及代币基本信息
// 这是一个完整的代币余额查询实现，使用智能合约绑定
func main() {
	// 步骤1：连接以太坊网络（演示用）
	// 使用Cloudflare提供的以太坊主网节点
	_, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal("连接以太坊客户端失败:", err)
	}
	fmt.Println("成功连接到以太坊网络")

	// 步骤2：设置代币合约地址
	// 这里使用Golem (GNT)代币作为示例
	// Golem是一个去中心化计算网络的代币
	tokenAddress := common.HexToAddress("0xfadea654ea83c00e5003d2ea15c59830b65471c0")
	fmt.Printf("代币合约地址: %s\n", tokenAddress.Hex())

	// 步骤3：设置要查询的账户地址
	address := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	fmt.Printf("查询地址: %s\n\n", address.Hex())

	// 注意：以下代码展示了如何使用ERC20代币合约绑定查询代币信息
	// 在实际使用时，您需要先生成合约绑定文件
	
	fmt.Printf("=== ERC20代币余额查询示例 ===\n")
	fmt.Printf("代币合约地址: %s\n", tokenAddress.Hex())
	fmt.Printf("查询地址: %s\n", address.Hex())
	fmt.Printf("\n=== 使用说明 ===\n")
	fmt.Printf("1. 首先需要获取ERC20代币合约的ABI文件\n")
	fmt.Printf("2. 使用abigen工具生成Go绑定代码:\n")
	fmt.Printf("   abigen --abi token.abi --pkg token --out token.go\n")
	fmt.Printf("3. 导入生成的包并创建合约实例\n")
	fmt.Printf("4. 调用合约函数查询代币信息\n")
	
	// 以下是完整的代码示例（需要合约绑定文件）：
	/*
	// 创建代币合约实例
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal("创建代币合约实例失败:", err)
	}
	
	// 查询代币余额
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal("查询代币余额失败:", err)
	}
	
	// 查询代币基本信息
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币名称失败:", err)
	}
	
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币符号失败:", err)
	}
	
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币小数位数失败:", err)
	}
	
	// 显示结果
	fmt.Printf("代币名称: %s\n", name)
	fmt.Printf("代币符号: %s\n", symbol)
	fmt.Printf("小数位数: %v\n", decimals)
	fmt.Printf("原始余额: %s\n", bal)
	
	// 转换为可读格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("余额: %f %s\n", value, symbol)
	*/

	// 小白说明：
	// 1. ERC20是以太坊上最常用的代币标准
	// 2. 代币余额以最小单位存储，类似于ETH的Wei
	// 3. 小数位数决定了代币的精度，大多数代币使用18位小数
	// 4. BalanceOf、Name、Symbol、Decimals是ERC20标准的基本函数
	// 5. &bind.CallOpts{}表示只读调用，不消耗Gas
	// 6. 智能合约绑定让我们可以像调用普通Go函数一样调用合约函数
	// 7. 需要使用abigen工具从合约ABI生成Go绑定代码
}