package common

// Config 项目配置结构
type Config struct {
	RPCURL   string // 以太坊节点RPC地址
	GasLimit uint64 // Gas限制
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		RPCURL:   "https://mainnet.infura.io/v3/YOUR_PROJECT_ID", // 请替换为实际的Infura项目ID
		GasLimit: 21000,
	}
}