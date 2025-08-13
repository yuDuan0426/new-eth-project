# Solidity 智能合约

本目录包含项目的 Solidity 智能合约文件。

## 合约列表

### SimpleStorage.sol
一个简单的存储合约，演示基本的状态变量操作：
- `set(uint256)` - 设置存储值
- `get()` - 获取存储值
- `increment()` - 增加存储值
- `reset()` - 重置存储值为0

### MyToken.sol
一个完整的 ERC20 代币合约实现，包含：
- 标准 ERC20 功能（转账、授权等）
- 铸币功能（仅限所有者）
- 销毁功能
- 所有者权限管理

## 编译合约

### 方法1：使用 Solidity 编译器 (solc)

```bash
# 安装 Solidity 编译器
npm install -g solc

# 编译合约
solc --bin --abi contracts/SimpleStorage.sol -o build/
solc --bin --abi contracts/MyToken.sol -o build/
```

### 方法2：使用 Hardhat

```bash
# 初始化 Hardhat 项目
npx hardhat init

# 编译合约
npx hardhat compile
```

### 方法3：使用 Foundry

```bash
# 安装 Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# 初始化项目
forge init

# 编译合约
forge build
```

## 生成 Go 绑定

使用 `abigen` 工具生成 Go 语言绑定：

```bash
# 安装 abigen（如果还没有安装）
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 生成 SimpleStorage 绑定
abigen --abi=build/SimpleStorage.abi --bin=build/SimpleStorage.bin --pkg=contracts --out=pkg/contracts/SimpleStorage.go

# 生成 MyToken 绑定
abigen --abi=build/MyToken.abi --bin=build/MyToken.bin --pkg=contracts --out=pkg/contracts/MyToken.go
```

## 部署合约

### 使用项目中的部署工具

```bash
# 部署 SimpleStorage 合约
go run cmd/deploy-contract/main.go
```

### 使用 Hardhat 部署

```javascript
// scripts/deploy.js
const { ethers } = require("hardhat");

async function main() {
  // 部署 SimpleStorage
  const SimpleStorage = await ethers.getContractFactory("SimpleStorage");
  const simpleStorage = await SimpleStorage.deploy();
  await simpleStorage.deployed();
  console.log("SimpleStorage deployed to:", simpleStorage.address);

  // 部署 MyToken
  const MyToken = await ethers.getContractFactory("MyToken");
  const myToken = await MyToken.deploy(
    "My Token",    // name
    "MTK",         // symbol
    18,            // decimals
    1000000        // initial supply
  );
  await myToken.deployed();
  console.log("MyToken deployed to:", myToken.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
```

## 测试合约

### Solidity 测试示例

```solidity
// test/SimpleStorage.test.sol
pragma solidity ^0.8.0;

import "../contracts/SimpleStorage.sol";

contract SimpleStorageTest {
    SimpleStorage simpleStorage;

    function setUp() public {
        simpleStorage = new SimpleStorage();
    }

    function testSet() public {
        simpleStorage.set(42);
        assert(simpleStorage.get() == 42);
    }

    function testIncrement() public {
        simpleStorage.set(10);
        simpleStorage.increment();
        assert(simpleStorage.get() == 11);
    }
}
```

### JavaScript 测试示例

```javascript
// test/SimpleStorage.test.js
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("SimpleStorage", function () {
  let simpleStorage;

  beforeEach(async function () {
    const SimpleStorage = await ethers.getContractFactory("SimpleStorage");
    simpleStorage = await SimpleStorage.deploy();
    await simpleStorage.deployed();
  });

  it("Should set and get value", async function () {
    await simpleStorage.set(42);
    expect(await simpleStorage.get()).to.equal(42);
  });

  it("Should increment value", async function () {
    await simpleStorage.set(10);
    await simpleStorage.increment();
    expect(await simpleStorage.get()).to.equal(11);
  });
});
```

## 在 Trae IDE 中使用 Solidity

Trae IDE 支持 Solidity 语言的语法高亮和基本编辑功能：

1. **语法高亮**：`.sol` 文件会自动识别并提供语法高亮
2. **代码补全**：支持基本的 Solidity 关键字补全
3. **错误检测**：可以检测基本的语法错误
4. **文件管理**：可以创建、编辑和管理 Solidity 文件

### 推荐的开发流程

1. 在 Trae IDE 中编写 Solidity 合约
2. 使用外部工具（如 Hardhat 或 Foundry）编译合约
3. 使用 `abigen` 生成 Go 绑定
4. 在 Go 代码中使用生成的绑定与合约交互

## 最佳实践

1. **版本管理**：始终指定 Solidity 版本
2. **许可证**：添加 SPDX 许可证标识符
3. **文档**：使用 NatSpec 注释文档化合约
4. **安全性**：使用 `require` 进行输入验证
5. **Gas 优化**：注意 Gas 消耗，优化合约代码
6. **测试**：编写全面的测试用例
7. **审计**：在主网部署前进行安全审计

## 相关资源

- [Solidity 官方文档](https://docs.soliditylang.org/)
- [OpenZeppelin 合约库](https://openzeppelin.com/contracts/)
- [Hardhat 开发框架](https://hardhat.org/)
- [Foundry 工具链](https://getfoundry.sh/)
- [Remix IDE](https://remix.ethereum.org/)