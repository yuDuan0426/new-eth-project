// 以太坊钱包创建工具
// 本程序演示如何生成以太坊钱包的私钥、公钥和地址
// 包含完整的密钥生成和地址推导流程

package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func main() {
	fmt.Println("=== 以太坊钱包创建工具 ===")
	fmt.Println("本工具演示如何生成新的以太坊钱包\n")

	// ===== 第1步：生成私钥 =====
	// 使用椭圆曲线数字签名算法(ECDSA)生成随机私钥
	// 以太坊使用secp256k1椭圆曲线，这与比特币相同
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("生成私钥失败:", err)
	}
	fmt.Println("✓ 成功生成私钥")

	// ===== 第2步：导出私钥 =====
	// 将私钥转换为字节数组，然后编码为十六进制字符串
	// 私钥是32字节(256位)的随机数
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// 去掉'0x'前缀，只保留十六进制字符串
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Printf("私钥 (Hex): %s\n", privateKeyHex)
	fmt.Printf("私钥长度: %d 字节\n\n", len(privateKeyBytes))

	// ===== 第3步：从私钥推导公钥 =====
	// 从私钥计算对应的公钥
	// 公钥是椭圆曲线上的一个点，由私钥通过椭圆曲线乘法得到
	publicKey := privateKey.Public()
	// 将接口类型转换为具体的ECDSA公钥类型
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型转换失败: 无法转换为*ecdsa.PublicKey类型")
	}
	fmt.Println("✓ 成功推导公钥")

	// ===== 第4步：导出公钥 =====
	// 将公钥转换为字节数组
	// 公钥包含x和y坐标，总共65字节(1字节前缀+32字节x+32字节y)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// 去掉'0x04'前缀，只保留坐标数据
	publicKeyHex := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Printf("公钥 (Hex): %s\n", publicKeyHex)
	fmt.Printf("公钥长度: %d 字节\n\n", len(publicKeyBytes))

	// ===== 第5步：从公钥推导以太坊地址 =====
	// 以太坊地址是公钥的Keccak256哈希值的后20字节
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("以太坊地址: %s\n\n", address)

	// ===== 第6步：手动验证地址推导过程 =====
	fmt.Println("=== 地址推导验证 ===")
	// 创建Keccak256哈希器
	hash := sha3.NewLegacyKeccak256()
	// 对公钥进行哈希(去掉第一个字节的前缀)
	// 公钥格式：0x04 + 32字节x坐标 + 32字节y坐标
	// 地址计算只使用x和y坐标，不包括0x04前缀
	hash.Write(publicKeyBytes[1:])

	// 获取完整的32字节哈希值
	fullHash := hash.Sum(nil)
	fullHashHex := hexutil.Encode(fullHash[:])
	fmt.Printf("公钥Keccak256哈希 (完整32字节): %s\n", fullHashHex)

	// 取哈希值的后20字节作为以太坊地址
	// 原长32字节，截去前12字节，保留后20字节
	addressFromHash := hexutil.Encode(fullHash[12:])
	fmt.Printf("地址 (后20字节): %s\n", addressFromHash)

	// ===== 第7步：验证结果一致性 =====
	fmt.Println("\n=== 验证结果 ===")
	if address == addressFromHash {
		fmt.Println("✓ 地址推导验证成功！两种方法得到相同结果")
	} else {
		fmt.Println("✗ 地址推导验证失败！结果不一致")
	}

	fmt.Println("\n=== 钱包信息汇总 ===")
	fmt.Printf("私钥: %s\n", privateKeyHex)
	fmt.Printf("公钥: %s\n", publicKeyHex)
	fmt.Printf("地址: %s\n", address)

	fmt.Println("\n=== 重要提醒 ===")
	fmt.Println("1. 私钥是您钱包的唯一凭证，请务必安全保管")
	fmt.Println("2. 任何人获得私钥都可以控制您的资产")
	fmt.Println("3. 建议使用硬件钱包或安全的密钥管理工具")
	fmt.Println("4. 不要在生产环境中使用此工具生成的钱包")
	fmt.Println("5. 这是演示代码，实际使用请采用更安全的方式")

	// ===== 技术说明 =====
	// 1. 以太坊钱包生成流程：
	//    - 生成256位随机私钥 → 计算公钥 → 计算地址
	//
	// 2. 关键概念：
	//    - 私钥：32字节随机数，控制钱包的唯一凭证
	//    - 公钥：椭圆曲线上的点，由私钥推导得出
	//    - 地址：公钥Keccak256哈希的后20字节
	//
	// 3. 椭圆曲线密码学：
	//    - 以太坊使用secp256k1曲线
	//    - 公钥 = 私钥 × G (G是曲线的生成点)
	//    - 这是单向函数，从公钥无法推导私钥
	//
	// 4. 地址生成算法：
	//    - 取公钥的x和y坐标(去掉0x04前缀)
	//    - 计算Keccak256哈希
	//    - 取哈希值的后20字节
	//    - 添加0x前缀得到最终地址
	//
	// 5. 安全注意事项：
	//    - 私钥必须真正随机生成
	//    - 私钥泄露等于资产丢失
	//    - 建议使用硬件随机数生成器
	//    - 生产环境使用专业的密钥管理方案
}