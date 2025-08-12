package common

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NewEthClient 创建以太坊客户端连接
func NewEthClient(rpcURL string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcURL)
}

// IsValidAddress 验证以太坊地址是否有效
func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

// HexToAddress 将十六进制字符串转换为Address对象
func HexToAddress(hex string) common.Address {
	return common.HexToAddress(hex)
}

// HexToPrivateKey 将十六进制字符串转换为私钥
func HexToPrivateKey(hex string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hex)
}

// PrivateKeyToAddress 从私钥推导出地址
func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

// EtherToWei 将以太币转换为Wei
func EtherToWei(ether *big.Float) *big.Int {
	wei := new(big.Float)
	wei.Mul(ether, big.NewFloat(1e18))
	result, _ := wei.Int(nil)
	return result
}

// WeiToEther 将Wei转换为以太币
func WeiToEther(wei *big.Int) *big.Float {
	ether := new(big.Float)
	ether.SetInt(wei)
	ether.Quo(ether, big.NewFloat(1e18))
	return ether
}