# Solidity 开发参考手册

本文档整理了 Solidity 开发中常用的 API、合约模板和最佳实践，方便快速查找和编写代码。

## 目录

- [基础语法](#基础语法)
- [数据类型](#数据类型)
- [全局变量和函数](#全局变量和函数)
- [常用合约模板](#常用合约模板)
- [安全模式](#安全模式)
- [Gas 优化](#gas-优化)
- [事件和日志](#事件和日志)
- [错误处理](#错误处理)

---

## 基础语法

### 合约结构模板

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// 导入其他合约或库
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title 合约标题
 * @dev 合约描述
 * @author 作者名称
 */
contract MyContract {
    // 状态变量
    uint256 public totalSupply;
    mapping(address => uint256) public balances;
    
    // 事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    
    // 修饰符
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }
    
    // 构造函数
    constructor(uint256 _initialSupply) {
        totalSupply = _initialSupply;
        balances[msg.sender] = _initialSupply;
    }
    
    // 函数
    function transfer(address to, uint256 amount) public returns (bool) {
        // 函数实现
        return true;
    }
}
```

### 函数可见性

```solidity
// public - 内部和外部都可调用
function publicFunction() public {}

// external - 只能从外部调用
function externalFunction() external {}

// internal - 只能从内部调用（包括继承合约）
function internalFunction() internal {}

// private - 只能从当前合约调用
function privateFunction() private {}
```

### 函数修饰符

```solidity
// view - 不修改状态
function getValue() public view returns (uint256) {
    return value;
}

// pure - 不读取也不修改状态
function add(uint256 a, uint256 b) public pure returns (uint256) {
    return a + b;
}

// payable - 可以接收以太币
function deposit() public payable {
    // 处理以太币
}
```

---

## 数据类型

### 基础类型

```solidity
// 布尔类型
bool public isActive = true;

// 整数类型
uint256 public largeNumber;  // 0 到 2^256 - 1
int256 public signedNumber;  // -2^255 到 2^255 - 1
uint8 public smallNumber;    // 0 到 255

// 地址类型
address public owner;
address payable public recipient;  // 可接收以太币的地址

// 字节类型
bytes32 public hash;
bytes public data;

// 字符串
string public name;

// 枚举
enum Status { Pending, Active, Inactive }
Status public currentStatus;
```

### 复合类型

```solidity
// 数组
uint256[] public dynamicArray;
uint256[10] public fixedArray;

// 映射
mapping(address => uint256) public balances;
mapping(address => mapping(address => uint256)) public allowances;

// 结构体
struct User {
    string name;
    uint256 age;
    bool isActive;
}
User[] public users;
mapping(address => User) public userInfo;
```

---

## 全局变量和函数

### 区块和交易属性

```solidity
// 区块信息
block.timestamp     // 当前区块时间戳
block.number        // 当前区块号
block.difficulty    // 当前区块难度
block.gaslimit      // 当前区块 gas 限制
block.coinbase      // 当前区块矿工地址

// 交易信息
msg.sender          // 消息发送者地址
msg.value           // 发送的以太币数量（wei）
msg.data            // 完整的调用数据
msg.sig             // 调用数据的前四个字节（函数标识符）

// 交易属性
tx.origin           // 交易发起者地址
tx.gasprice         // 交易 gas 价格
```

### 地址相关函数

```solidity
// 地址余额
address(this).balance           // 当前合约余额
owner.balance                   // 指定地址余额

// 转账
payable(recipient).transfer(amount);    // 转账（失败时回滚）
payable(recipient).send(amount);        // 转账（返回布尔值）

// 低级调用
(bool success, bytes memory data) = target.call{value: amount}(abi.encodeWithSignature("function()"));
```

### 加密函数

```solidity
// 哈希函数
keccak256(abi.encodePacked(data))       // Keccak-256 哈希
sha256(abi.encodePacked(data))          // SHA-256 哈希
ripemd160(abi.encodePacked(data))       // RIPEMD-160 哈希

// 签名验证
ecrecover(hash, v, r, s)                // 从签名恢复地址
```

### ABI 编码函数

```solidity
// ABI 编码
abi.encode(param1, param2)              // 标准 ABI 编码
abi.encodePacked(param1, param2)        // 紧密打包编码
abi.encodeWithSignature("func(uint256)", param)  // 带函数签名编码
abi.encodeWithSelector(selector, param) // 带选择器编码

// ABI 解码
(uint256 a, string memory b) = abi.decode(data, (uint256, string));
```

---

## 常用合约模板

### 1. ERC20 代币合约

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

contract ERC20Token is IERC20 {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    
    uint256 private _totalSupply;
    string public name;
    string public symbol;
    uint8 public decimals;
    
    constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _totalSupply) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        _totalSupply = _totalSupply * 10**_decimals;
        _balances[msg.sender] = _totalSupply;
        emit Transfer(address(0), msg.sender, _totalSupply);
    }
    
    function totalSupply() public view override returns (uint256) {
        return _totalSupply;
    }
    
    function balanceOf(address account) public view override returns (uint256) {
        return _balances[account];
    }
    
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }
    
    function allowance(address owner, address spender) public view override returns (uint256) {
        return _allowances[owner][spender];
    }
    
    function approve(address spender, uint256 amount) public override returns (bool) {
        _approve(msg.sender, spender, amount);
        return true;
    }
    
    function transferFrom(address sender, address recipient, uint256 amount) public override returns (bool) {
        uint256 currentAllowance = _allowances[sender][msg.sender];
        require(currentAllowance >= amount, "ERC20: transfer amount exceeds allowance");
        
        _transfer(sender, recipient, amount);
        _approve(sender, msg.sender, currentAllowance - amount);
        
        return true;
    }
    
    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "ERC20: transfer from the zero address");
        require(recipient != address(0), "ERC20: transfer to the zero address");
        require(_balances[sender] >= amount, "ERC20: transfer amount exceeds balance");
        
        _balances[sender] -= amount;
        _balances[recipient] += amount;
        emit Transfer(sender, recipient, amount);
    }
    
    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");
        
        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }
}
```

### 2. 访问控制合约

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AccessControl {
    address public owner;
    mapping(address => bool) public admins;
    mapping(address => bool) public operators;
    
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event AdminAdded(address indexed admin);
    event AdminRemoved(address indexed admin);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }
    
    modifier onlyAdmin() {
        require(admins[msg.sender] || msg.sender == owner, "Not an admin");
        _;
    }
    
    modifier onlyOperator() {
        require(operators[msg.sender] || admins[msg.sender] || msg.sender == owner, "Not an operator");
        _;
    }
    
    constructor() {
        owner = msg.sender;
        admins[msg.sender] = true;
    }
    
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner is the zero address");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }
    
    function addAdmin(address admin) public onlyOwner {
        admins[admin] = true;
        emit AdminAdded(admin);
    }
    
    function removeAdmin(address admin) public onlyOwner {
        admins[admin] = false;
        emit AdminRemoved(admin);
    }
    
    function addOperator(address operator) public onlyAdmin {
        operators[operator] = true;
    }
    
    function removeOperator(address operator) public onlyAdmin {
        operators[operator] = false;
    }
}
```

### 3. 多签钱包合约

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MultiSigWallet {
    event Deposit(address indexed sender, uint256 amount, uint256 balance);
    event SubmitTransaction(address indexed owner, uint256 indexed txIndex, address indexed to, uint256 value, bytes data);
    event ConfirmTransaction(address indexed owner, uint256 indexed txIndex);
    event RevokeConfirmation(address indexed owner, uint256 indexed txIndex);
    event ExecuteTransaction(address indexed owner, uint256 indexed txIndex);
    
    address[] public owners;
    mapping(address => bool) public isOwner;
    uint256 public numConfirmationsRequired;
    
    struct Transaction {
        address to;
        uint256 value;
        bytes data;
        bool executed;
        uint256 numConfirmations;
    }
    
    mapping(uint256 => mapping(address => bool)) public isConfirmed;
    Transaction[] public transactions;
    
    modifier onlyOwner() {
        require(isOwner[msg.sender], "not owner");
        _;
    }
    
    modifier txExists(uint256 _txIndex) {
        require(_txIndex < transactions.length, "tx does not exist");
        _;
    }
    
    modifier notExecuted(uint256 _txIndex) {
        require(!transactions[_txIndex].executed, "tx already executed");
        _;
    }
    
    modifier notConfirmed(uint256 _txIndex) {
        require(!isConfirmed[_txIndex][msg.sender], "tx already confirmed");
        _;
    }
    
    constructor(address[] memory _owners, uint256 _numConfirmationsRequired) {
        require(_owners.length > 0, "owners required");
        require(_numConfirmationsRequired > 0 && _numConfirmationsRequired <= _owners.length, "invalid number of required confirmations");
        
        for (uint256 i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            require(owner != address(0), "invalid owner");
            require(!isOwner[owner], "owner not unique");
            
            isOwner[owner] = true;
            owners.push(owner);
        }
        
        numConfirmationsRequired = _numConfirmationsRequired;
    }
    
    receive() external payable {
        emit Deposit(msg.sender, msg.value, address(this).balance);
    }
    
    function submitTransaction(address _to, uint256 _value, bytes memory _data) public onlyOwner {
        uint256 txIndex = transactions.length;
        
        transactions.push(Transaction({
            to: _to,
            value: _value,
            data: _data,
            executed: false,
            numConfirmations: 0
        }));
        
        emit SubmitTransaction(msg.sender, txIndex, _to, _value, _data);
    }
    
    function confirmTransaction(uint256 _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) notConfirmed(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations += 1;
        isConfirmed[_txIndex][msg.sender] = true;
        
        emit ConfirmTransaction(msg.sender, _txIndex);
    }
    
    function executeTransaction(uint256 _txIndex) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        
        require(transaction.numConfirmations >= numConfirmationsRequired, "cannot execute tx");
        
        transaction.executed = true;
        
        (bool success, ) = transaction.to.call{value: transaction.value}(transaction.data);
        require(success, "tx failed");
        
        emit ExecuteTransaction(msg.sender, _txIndex);
    }
}
```

### 4. 时间锁合约

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TimeLock {
    event QueueTransaction(bytes32 indexed txHash, address indexed target, uint256 value, string signature, bytes data, uint256 eta);
    event ExecuteTransaction(bytes32 indexed txHash, address indexed target, uint256 value, string signature, bytes data, uint256 eta);
    event CancelTransaction(bytes32 indexed txHash, address indexed target, uint256 value, string signature, bytes data, uint256 eta);
    
    uint256 public constant GRACE_PERIOD = 14 days;
    uint256 public constant MINIMUM_DELAY = 2 days;
    uint256 public constant MAXIMUM_DELAY = 30 days;
    
    address public admin;
    uint256 public delay;
    mapping(bytes32 => bool) public queuedTransactions;
    
    modifier onlyAdmin() {
        require(msg.sender == admin, "TimeLock: Call must come from admin");
        _;
    }
    
    constructor(address _admin, uint256 _delay) {
        require(_delay >= MINIMUM_DELAY, "TimeLock: Delay must exceed minimum delay");
        require(_delay <= MAXIMUM_DELAY, "TimeLock: Delay must not exceed maximum delay");
        
        admin = _admin;
        delay = _delay;
    }
    
    function queueTransaction(address target, uint256 value, string memory signature, bytes memory data, uint256 eta) public onlyAdmin returns (bytes32) {
        require(eta >= getBlockTimestamp() + delay, "TimeLock: Estimated execution block must satisfy delay");
        
        bytes32 txHash = keccak256(abi.encode(target, value, signature, data, eta));
        queuedTransactions[txHash] = true;
        
        emit QueueTransaction(txHash, target, value, signature, data, eta);
        return txHash;
    }
    
    function cancelTransaction(address target, uint256 value, string memory signature, bytes memory data, uint256 eta) public onlyAdmin {
        bytes32 txHash = keccak256(abi.encode(target, value, signature, data, eta));
        queuedTransactions[txHash] = false;
        
        emit CancelTransaction(txHash, target, value, signature, data, eta);
    }
    
    function executeTransaction(address target, uint256 value, string memory signature, bytes memory data, uint256 eta) public payable onlyAdmin returns (bytes memory) {
        bytes32 txHash = keccak256(abi.encode(target, value, signature, data, eta));
        require(queuedTransactions[txHash], "TimeLock: Transaction hasn't been queued");
        require(getBlockTimestamp() >= eta, "TimeLock: Transaction hasn't surpassed time lock");
        require(getBlockTimestamp() <= eta + GRACE_PERIOD, "TimeLock: Transaction is stale");
        
        queuedTransactions[txHash] = false;
        
        bytes memory callData;
        if (bytes(signature).length == 0) {
            callData = data;
        } else {
            callData = abi.encodePacked(bytes4(keccak256(bytes(signature))), data);
        }
        
        (bool success, bytes memory returnData) = target.call{value: value}(callData);
        require(success, "TimeLock: Transaction execution reverted");
        
        emit ExecuteTransaction(txHash, target, value, signature, data, eta);
        
        return returnData;
    }
    
    function getBlockTimestamp() internal view returns (uint256) {
        return block.timestamp;
    }
}
```

---

## 安全模式

### 重入攻击防护

```solidity
// 使用 ReentrancyGuard
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract MyContract is ReentrancyGuard {
    function withdraw() public nonReentrant {
        // 安全的提取逻辑
    }
}

// 手动实现重入保护
contract ManualReentrancyGuard {
    bool private _notEntered = true;
    
    modifier nonReentrant() {
        require(_notEntered, "ReentrancyGuard: reentrant call");
        _notEntered = false;
        _;
        _notEntered = true;
    }
}
```

### 检查-效果-交互模式

```solidity
function withdraw(uint256 amount) public {
    // 1. 检查
    require(balances[msg.sender] >= amount, "Insufficient balance");
    
    // 2. 效果
    balances[msg.sender] -= amount;
    
    // 3. 交互
    payable(msg.sender).transfer(amount);
}
```

### 安全的随机数生成

```solidity
// 不安全的方式（不要使用）
// uint256 random = uint256(keccak256(abi.encodePacked(block.timestamp, block.difficulty)));

// 使用 Chainlink VRF（推荐）
import "@chainlink/contracts/src/v0.8/VRFConsumerBase.sol";

contract RandomNumberConsumer is VRFConsumerBase {
    bytes32 internal keyHash;
    uint256 internal fee;
    uint256 public randomResult;
    
    constructor() VRFConsumerBase(
        0xb3dCcb4Cf7a26f6cf6B120Cf5A73875B7BBc655C, // VRF Coordinator
        0x01BE23585060835E02B77ef475b0Cc51aA1e0709  // LINK Token
    ) {
        keyHash = 0x2ed0feb3e7fd2022120aa84fab1945545a9f2ffc9076fd6156fa96eaff4c1311;
        fee = 0.1 * 10 ** 18; // 0.1 LINK
    }
    
    function getRandomNumber() public returns (bytes32 requestId) {
        require(LINK.balanceOf(address(this)) >= fee, "Not enough LINK");
        return requestRandomness(keyHash, fee);
    }
    
    function fulfillRandomness(bytes32 requestId, uint256 randomness) internal override {
        randomResult = randomness;
    }
}
```

---

## Gas 优化

### 存储优化

```solidity
// 打包结构体以节省存储
struct OptimizedStruct {
    uint128 value1;  // 16 bytes
    uint128 value2;  // 16 bytes
    // 总共 32 bytes，占用一个存储槽
}

struct UnoptimizedStruct {
    uint256 value1;  // 32 bytes
    uint128 value2;  // 16 bytes + 16 bytes padding
    // 总共 64 bytes，占用两个存储槽
}

// 使用常量和不可变变量
uint256 public constant RATE = 100;  // 编译时常量
uint256 public immutable deployTime; // 部署时设置

constructor() {
    deployTime = block.timestamp;
}
```

### 循环优化

```solidity
// 优化前
function inefficientLoop(uint256[] memory arr) public {
    for (uint256 i = 0; i < arr.length; i++) {
        // 每次都读取 arr.length
    }
}

// 优化后
function efficientLoop(uint256[] memory arr) public {
    uint256 length = arr.length;
    for (uint256 i = 0; i < length; i++) {
        // 只读取一次 length
    }
}

// 使用 unchecked 块（Solidity 0.8+）
function uncheckedLoop(uint256[] memory arr) public {
    uint256 length = arr.length;
    for (uint256 i = 0; i < length;) {
        // 循环体
        unchecked {
            i++;
        }
    }
}
```

---

## 事件和日志

### 事件定义和使用

```solidity
// 事件定义
event Transfer(address indexed from, address indexed to, uint256 value);
event Approval(address indexed owner, address indexed spender, uint256 value);
event StateChanged(uint256 indexed id, string oldState, string newState);

// 触发事件
function transfer(address to, uint256 amount) public {
    // 转账逻辑
    emit Transfer(msg.sender, to, amount);
}

// 带多个索引的事件
event ComplexEvent(
    address indexed user,
    uint256 indexed tokenId,
    bytes32 indexed category,
    string data,
    uint256 timestamp
);
```

### 日志查询优化

```solidity
// 使用索引参数进行高效查询
event UserAction(
    address indexed user,      // 可按用户查询
    uint256 indexed actionType, // 可按动作类型查询
    bytes32 indexed category,   // 可按分类查询
    string data,               // 不索引的数据
    uint256 timestamp          // 不索引的时间戳
);

// 最多只能有 3 个索引参数
```

---

## 错误处理

### 自定义错误（Solidity 0.8.4+）

```solidity
// 定义自定义错误
error InsufficientBalance(uint256 available, uint256 required);
error Unauthorized(address caller);
error InvalidAddress();

contract ErrorHandling {
    mapping(address => uint256) public balances;
    address public owner;
    
    modifier onlyOwner() {
        if (msg.sender != owner) {
            revert Unauthorized(msg.sender);
        }
        _;
    }
    
    function transfer(address to, uint256 amount) public {
        if (to == address(0)) {
            revert InvalidAddress();
        }
        
        if (balances[msg.sender] < amount) {
            revert InsufficientBalance(balances[msg.sender], amount);
        }
        
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }
}
```

### 传统错误处理

```solidity
// require 语句
require(condition, "Error message");
require(balances[msg.sender] >= amount, "Insufficient balance");

// assert 语句（用于内部错误）
assert(totalSupply >= balances[msg.sender]);

// revert 语句
if (condition) {
    revert("Error message");
}
```

---

## 常用库和工具

### OpenZeppelin 合约

```solidity
// 访问控制
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

// 安全
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

// 代币标准
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";

// 工具
import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
```

### 数学库使用

```solidity
// Solidity 0.8+ 内置溢出检查
uint256 result = a + b;  // 自动检查溢出

// 使用 unchecked 跳过检查（谨慎使用）
unchecked {
    uint256 result = a + b;
}

// OpenZeppelin SafeMath（0.8 之前版本）
using SafeMath for uint256;
uint256 result = a.add(b);
```

---

## 最佳实践总结

1. **安全第一**：使用 OpenZeppelin 库，遵循安全模式
2. **Gas 优化**：合理使用存储，优化循环和计算
3. **代码清晰**：使用有意义的变量名和函数名
4. **文档完整**：使用 NatSpec 注释
5. **测试充分**：编写全面的测试用例
6. **版本管理**：明确指定 Solidity 版本
7. **事件记录**：重要操作都要触发事件
8. **错误处理**：使用自定义错误提高 Gas 效率
9. **权限控制**：实现适当的访问控制机制
10. **升级策略**：考虑合约升级的需求

---

## 开发工具推荐

- **Hardhat**：开发框架和测试环境
- **Foundry**：快速的开发工具链
- **Remix**：在线 IDE
- **Truffle**：经典开发框架
- **OpenZeppelin Wizard**：合约生成器
- **Slither**：静态分析工具
- **MythX**：安全分析平台

这个参考手册涵盖了 Solidity 开发的核心内容，可以作为日常开发的快速查询工具。