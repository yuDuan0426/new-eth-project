# 以太坊开发工作流程指南

本文档以伪代码形式记录以太坊开发中各个模块的标准流程，便于理解和实现。

## 目录

1. [钱包创建流程](#钱包创建流程)
2. [余额查询流程](#余额查询流程)
3. [ETH转账流程](#eth转账流程)
4. [代币转账流程](#代币转账流程)
5. [区块查询流程](#区块查询流程)
6. [交易查询流程](#交易查询流程)
7. [交易回执查询流程](#交易回执查询流程)
8. [智能合约部署流程](#智能合约部署流程)
9. [智能合约调用流程](#智能合约调用流程)
10. [智能合约加载流程](#智能合约加载流程)
11. [事件监听流程](#事件监听流程)
12. [历史事件查询流程](#历史事件查询流程)
13. [区块订阅流程](#区块订阅流程)

---

## 钱包创建流程

### 流程概述
生成新的以太坊钱包，包括私钥、公钥和地址的完整推导过程。

### 伪代码

```pseudocode
FUNCTION createWallet():
    // 第1步：生成私钥
    privateKey = generateRandomPrivateKey()
    IF privateKey == NULL:
        THROW "私钥生成失败"
    
    // 第2步：导出私钥
    privateKeyBytes = convertToBytes(privateKey)
    privateKeyHex = bytesToHex(privateKeyBytes)
    
    // 第3步：从私钥推导公钥
    publicKey = derivePublicKey(privateKey)
    IF publicKey == NULL:
        THROW "公钥推导失败"
    
    // 第4步：导出公钥
    publicKeyBytes = convertToBytes(publicKey)
    publicKeyHex = bytesToHex(publicKeyBytes)
    
    // 第5步：从公钥生成地址
    address = deriveAddress(publicKey)
    
    // 第6步：验证地址推导
    manualAddress = keccak256Hash(publicKeyBytes[1:])[12:]
    IF address != manualAddress:
        THROW "地址推导验证失败"
    
    // 第7步：返回钱包信息
    RETURN {
        privateKey: privateKeyHex,
        publicKey: publicKeyHex,
        address: address.hex(),
        verified: true
    }
END FUNCTION
```

### 关键步骤
1. **私钥生成**：使用密码学安全的随机数生成器
2. **公钥推导**：通过椭圆曲线数字签名算法(ECDSA)
3. **地址生成**：对公钥进行Keccak256哈希，取后20字节
4. **验证检查**：确保推导过程的正确性

---

## 余额查询流程

### 流程概述
查询指定地址的ETH余额，支持指定区块号查询历史余额。

### 伪代码

```pseudocode
FUNCTION queryBalance(address, blockNumber = "latest"):
    // 第1步：连接以太坊网络
    client = connectToEthereum(RPC_ENDPOINT)
    IF client == NULL:
        THROW "网络连接失败"
    
    // 第2步：验证地址格式
    IF NOT isValidAddress(address):
        THROW "无效的地址格式"
    
    // 第3步：查询余额
    balanceWei = client.getBalance(address, blockNumber)
    IF balanceWei == NULL:
        THROW "余额查询失败"
    
    // 第4步：单位转换
    balanceEth = weiToEth(balanceWei)
    
    // 第5步：返回结果
    RETURN {
        address: address,
        balanceWei: balanceWei.toString(),
        balanceEth: balanceEth.toString(),
        blockNumber: blockNumber
    }
END FUNCTION

// 辅助函数：Wei转ETH
FUNCTION weiToEth(wei):
    RETURN wei / (10^18)
END FUNCTION
```

### 关键步骤
1. **网络连接**：建立与以太坊节点的连接
2. **地址验证**：确保地址格式正确
3. **余额查询**：调用RPC接口获取余额
4. **单位转换**：将Wei转换为ETH便于阅读

---

## ETH转账流程

### 流程概述
发送ETH从一个地址到另一个地址，包括交易创建、签名和广播。

### 伪代码

```pseudocode
FUNCTION transferETH(fromPrivateKey, toAddress, amount, gasPrice, gasLimit):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：从私钥推导发送方地址
    fromAddress = deriveAddressFromPrivateKey(fromPrivateKey)
    
    // 第3步：获取nonce
    nonce = client.getPendingNonce(fromAddress)
    
    // 第4步：获取链ID
    chainID = client.getChainID()
    
    // 第5步：转换金额单位
    amountWei = ethToWei(amount)
    
    // 第6步：创建交易
    transaction = createTransaction({
        nonce: nonce,
        to: toAddress,
        value: amountWei,
        gasLimit: gasLimit,
        gasPrice: gasPrice,
        data: NULL
    })
    
    // 第7步：签名交易
    signedTransaction = signTransaction(transaction, fromPrivateKey, chainID)
    
    // 第8步：发送交易
    txHash = client.sendTransaction(signedTransaction)
    
    // 第9步：等待确认（可选）
    receipt = waitForReceipt(client, txHash)
    
    RETURN {
        txHash: txHash,
        from: fromAddress,
        to: toAddress,
        amount: amount,
        status: receipt.status
    }
END FUNCTION

// 辅助函数：等待交易确认
FUNCTION waitForReceipt(client, txHash):
    WHILE true:
        receipt = client.getTransactionReceipt(txHash)
        IF receipt != NULL:
            RETURN receipt
        SLEEP(1000) // 等待1秒
    END WHILE
END FUNCTION
```

### 关键步骤
1. **Nonce获取**：防止重放攻击
2. **交易创建**：设置转账参数
3. **交易签名**：使用私钥和链ID签名
4. **交易广播**：发送到网络
5. **确认等待**：等待矿工打包确认

---

## 代币转账流程

### 流程概述
发送ERC20代币，需要调用代币合约的transfer函数。

### 伪代码

```pseudocode
FUNCTION transferToken(fromPrivateKey, tokenAddress, toAddress, amount):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：获取发送方地址
    fromAddress = deriveAddressFromPrivateKey(fromPrivateKey)
    
    // 第3步：获取代币精度
    decimals = getTokenDecimals(client, tokenAddress)
    
    // 第4步：转换代币数量
    tokenAmount = amount * (10^decimals)
    
    // 第5步：编码transfer函数调用
    functionSignature = "transfer(address,uint256)"
    methodID = keccak256(functionSignature)[0:4]
    encodedParams = encodeParameters(toAddress, tokenAmount)
    callData = methodID + encodedParams
    
    // 第6步：获取交易参数
    nonce = client.getPendingNonce(fromAddress)
    gasPrice = client.suggestGasPrice()
    
    // 第7步：估算Gas
    gasLimit = client.estimateGas({
        from: fromAddress,
        to: tokenAddress,
        data: callData
    })
    
    // 第8步：创建交易
    transaction = createTransaction({
        nonce: nonce,
        to: tokenAddress,
        value: 0, // 代币转账不发送ETH
        gasLimit: gasLimit,
        gasPrice: gasPrice,
        data: callData
    })
    
    // 第9步：签名并发送
    chainID = client.getChainID()
    signedTx = signTransaction(transaction, fromPrivateKey, chainID)
    txHash = client.sendTransaction(signedTx)
    
    RETURN {
        txHash: txHash,
        tokenAddress: tokenAddress,
        from: fromAddress,
        to: toAddress,
        amount: amount
    }
END FUNCTION

// 辅助函数：获取代币精度
FUNCTION getTokenDecimals(client, tokenAddress):
    functionSig = "decimals()"
    methodID = keccak256(functionSig)[0:4]
    result = client.call({
        to: tokenAddress,
        data: methodID
    })
    RETURN decodeUint256(result)
END FUNCTION
```

### 关键步骤
1. **代币精度获取**：确定代币的小数位数
2. **函数调用编码**：编码transfer函数调用数据
3. **Gas估算**：预估交易所需Gas
4. **合约调用**：发送交易到代币合约

---

## 区块查询流程

### 流程概述
查询区块信息，支持按区块号和区块哈希查询。

### 伪代码

```pseudocode
FUNCTION queryBlock(blockIdentifier, includeTransactions = false):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：确定查询方式
    IF isNumber(blockIdentifier):
        // 按区块号查询
        IF includeTransactions:
            block = client.getBlockByNumber(blockIdentifier, true)
        ELSE:
            header = client.getBlockHeaderByNumber(blockIdentifier)
        END IF
    ELSE IF isHash(blockIdentifier):
        // 按区块哈希查询
        block = client.getBlockByHash(blockIdentifier, includeTransactions)
    ELSE:
        THROW "无效的区块标识符"
    END IF
    
    // 第3步：验证查询结果
    IF block == NULL AND header == NULL:
        THROW "区块不存在"
    END IF
    
    // 第4步：提取区块信息
    IF includeTransactions:
        result = {
            number: block.number,
            hash: block.hash,
            timestamp: block.timestamp,
            difficulty: block.difficulty,
            gasLimit: block.gasLimit,
            gasUsed: block.gasUsed,
            transactionCount: block.transactions.length,
            transactions: block.transactions,
            size: block.size
        }
    ELSE:
        result = {
            number: header.number,
            hash: header.hash,
            timestamp: header.timestamp,
            difficulty: header.difficulty,
            gasLimit: header.gasLimit,
            gasUsed: header.gasUsed
        }
    END IF
    
    // 第5步：数据一致性验证（可选）
    IF includeTransactions:
        txCount = client.getTransactionCount(block.hash)
        IF txCount != block.transactions.length:
            THROW "交易数量不一致"
        END IF
    END IF
    
    RETURN result
END FUNCTION
```

### 关键步骤
1. **查询方式判断**：区分按号码还是哈希查询
2. **数据获取**：选择区块头或完整区块
3. **结果验证**：确保数据完整性
4. **信息提取**：格式化返回数据

---

## 交易查询流程

### 流程概述
查询交易详细信息，包括交易状态和发送方地址恢复。

### 伪代码

```pseudocode
FUNCTION queryTransaction(txHash):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：获取链ID
    chainID = client.getChainID()
    
    // 第3步：查询交易
    transaction, isPending = client.getTransactionByHash(txHash)
    IF transaction == NULL:
        THROW "交易不存在"
    END IF
    
    // 第4步：恢复发送方地址
    signer = createEIP155Signer(chainID)
    senderAddress = recoverSender(signer, transaction)
    
    // 第5步：获取交易回执（如果已确认）
    receipt = NULL
    IF NOT isPending:
        receipt = client.getTransactionReceipt(txHash)
    END IF
    
    // 第6步：组装结果
    result = {
        hash: transaction.hash,
        from: senderAddress,
        to: transaction.to,
        value: transaction.value,
        gas: transaction.gas,
        gasPrice: transaction.gasPrice,
        nonce: transaction.nonce,
        data: transaction.data,
        isPending: isPending
    }
    
    // 第7步：添加回执信息（如果有）
    IF receipt != NULL:
        result.status = receipt.status
        result.gasUsed = receipt.gasUsed
        result.blockNumber = receipt.blockNumber
        result.blockHash = receipt.blockHash
        result.transactionIndex = receipt.transactionIndex
    END IF
    
    RETURN result
END FUNCTION

// 多种查询方式
FUNCTION queryTransactionsByBlock(blockIdentifier):
    // 第1步：获取区块
    block = client.getBlockByNumber(blockIdentifier, true)
    
    // 第2步：遍历交易
    transactions = []
    FOR each tx IN block.transactions:
        txInfo = queryTransaction(tx.hash)
        transactions.append(txInfo)
    END FOR
    
    RETURN transactions
END FUNCTION
```

### 关键步骤
1. **交易获取**：通过哈希查询交易
2. **地址恢复**：从签名中恢复发送方地址
3. **状态检查**：判断交易是否已确认
4. **回执获取**：获取执行结果

---

## 交易回执查询流程

### 流程概述
查询交易执行回执，包含Gas使用、状态和事件日志。

### 伪代码

```pseudocode
FUNCTION queryTransactionReceipt(txHash):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：查询回执
    receipt = client.getTransactionReceipt(txHash)
    IF receipt == NULL:
        THROW "交易回执不存在（可能未确认）"
    END IF
    
    // 第3步：解析回执信息
    result = {
        transactionHash: receipt.transactionHash,
        transactionIndex: receipt.transactionIndex,
        blockHash: receipt.blockHash,
        blockNumber: receipt.blockNumber,
        from: receipt.from,
        to: receipt.to,
        gasUsed: receipt.gasUsed,
        cumulativeGasUsed: receipt.cumulativeGasUsed,
        status: receipt.status, // 1=成功, 0=失败
        logs: [],
        contractAddress: receipt.contractAddress // 合约创建时有值
    }
    
    // 第4步：解析事件日志
    FOR each log IN receipt.logs:
        logInfo = {
            address: log.address,
            topics: log.topics,
            data: log.data,
            blockNumber: log.blockNumber,
            transactionHash: log.transactionHash,
            logIndex: log.logIndex
        }
        result.logs.append(logInfo)
    END FOR
    
    // 第5步：状态判断
    IF result.status == 1:
        result.statusText = "成功"
    ELSE:
        result.statusText = "失败"
    END IF
    
    RETURN result
END FUNCTION

// 批量查询区块回执
FUNCTION queryBlockReceipts(blockNumber):
    // 第1步：获取区块交易
    block = client.getBlockByNumber(blockNumber, true)
    
    // 第2步：批量查询回执
    receipts = []
    FOR each tx IN block.transactions:
        receipt = queryTransactionReceipt(tx.hash)
        receipts.append(receipt)
    END FOR
    
    // 第3步：统计信息
    totalGasUsed = 0
    successCount = 0
    FOR each receipt IN receipts:
        totalGasUsed += receipt.gasUsed
        IF receipt.status == 1:
            successCount += 1
        END IF
    END FOR
    
    RETURN {
        blockNumber: blockNumber,
        totalTransactions: receipts.length,
        successfulTransactions: successCount,
        totalGasUsed: totalGasUsed,
        receipts: receipts
    }
END FUNCTION
```

### 关键步骤
1. **回执查询**：获取交易执行结果
2. **状态解析**：判断交易成功或失败
3. **日志解析**：提取事件日志信息
4. **统计分析**：计算Gas使用等统计数据

---

## 智能合约部署流程

### 流程概述
部署智能合约到区块链，获取合约地址。

### 伪代码

```pseudocode
FUNCTION deployContract(privateKey, contractBytecode, constructorParams, gasLimit, gasPrice):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：获取部署者地址
    deployerAddress = deriveAddressFromPrivateKey(privateKey)
    
    // 第3步：获取nonce
    nonce = client.getPendingNonce(deployerAddress)
    
    // 第4步：编码构造函数参数
    encodedParams = encodeConstructorParams(constructorParams)
    deploymentData = contractBytecode + encodedParams
    
    // 第5步：创建合约部署交易
    transaction = createContractCreationTransaction({
        nonce: nonce,
        value: 0, // 通常不发送ETH给构造函数
        gasLimit: gasLimit,
        gasPrice: gasPrice,
        data: deploymentData
    })
    
    // 第6步：签名交易
    chainID = client.getChainID()
    signedTx = signTransaction(transaction, privateKey, chainID)
    
    // 第7步：发送交易
    txHash = client.sendTransaction(signedTx)
    
    // 第8步：等待确认
    receipt = waitForReceipt(client, txHash)
    
    // 第9步：检查部署结果
    IF receipt.status != 1:
        THROW "合约部署失败"
    END IF
    
    // 第10步：验证合约代码
    contractAddress = receipt.contractAddress
    deployedCode = client.getCode(contractAddress)
    IF deployedCode.length == 0:
        THROW "合约代码验证失败"
    END IF
    
    RETURN {
        contractAddress: contractAddress,
        transactionHash: txHash,
        blockNumber: receipt.blockNumber,
        gasUsed: receipt.gasUsed,
        deployerAddress: deployerAddress
    }
END FUNCTION
```

### 关键步骤
1. **字节码准备**：合约编译后的字节码
2. **参数编码**：构造函数参数ABI编码
3. **交易创建**：创建合约部署交易
4. **地址获取**：从回执中获取合约地址
5. **代码验证**：确认合约部署成功

---

## 智能合约调用流程

### 流程概述
调用已部署的智能合约函数，包括只读调用和写入调用。

### 伪代码

```pseudocode
// 只读调用（不消耗Gas）
FUNCTION callContractReadOnly(contractAddress, functionName, params):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：编码函数调用
    functionSignature = generateFunctionSignature(functionName, params)
    methodID = keccak256(functionSignature)[0:4]
    encodedParams = encodeParameters(params)
    callData = methodID + encodedParams
    
    // 第3步：创建调用消息
    callMsg = {
        to: contractAddress,
        data: callData
    }
    
    // 第4步：执行调用
    result = client.call(callMsg, "latest")
    
    // 第5步：解码返回值
    decodedResult = decodeReturnValue(result, functionName)
    
    RETURN decodedResult
END FUNCTION

// 写入调用（消耗Gas，需要发送交易）
FUNCTION callContractWrite(privateKey, contractAddress, functionName, params, value, gasLimit, gasPrice):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：获取调用者地址
    callerAddress = deriveAddressFromPrivateKey(privateKey)
    
    // 第3步：编码函数调用
    functionSignature = generateFunctionSignature(functionName, params)
    methodID = keccak256(functionSignature)[0:4]
    encodedParams = encodeParameters(params)
    callData = methodID + encodedParams
    
    // 第4步：获取nonce
    nonce = client.getPendingNonce(callerAddress)
    
    // 第5步：创建交易
    transaction = createTransaction({
        nonce: nonce,
        to: contractAddress,
        value: value, // 发送给合约的ETH数量
        gasLimit: gasLimit,
        gasPrice: gasPrice,
        data: callData
    })
    
    // 第6步：签名并发送
    chainID = client.getChainID()
    signedTx = signTransaction(transaction, privateKey, chainID)
    txHash = client.sendTransaction(signedTx)
    
    // 第7步：等待确认
    receipt = waitForReceipt(client, txHash)
    
    // 第8步：检查执行结果
    IF receipt.status != 1:
        THROW "合约调用失败"
    END IF
    
    RETURN {
        transactionHash: txHash,
        blockNumber: receipt.blockNumber,
        gasUsed: receipt.gasUsed,
        logs: receipt.logs
    }
END FUNCTION

// 辅助函数：生成函数签名
FUNCTION generateFunctionSignature(functionName, params):
    paramTypes = []
    FOR each param IN params:
        paramTypes.append(getParameterType(param))
    END FOR
    RETURN functionName + "(" + join(paramTypes, ",") + ")"
END FUNCTION
```

### 关键步骤
1. **函数签名生成**：根据函数名和参数类型
2. **参数编码**：ABI编码函数参数
3. **调用方式选择**：只读call或写入transaction
4. **结果解码**：解析返回值或事件日志

---

## 智能合约加载流程

### 流程概述
加载已部署的智能合约，验证合约存在性。

### 伪代码

```pseudocode
FUNCTION loadContract(contractAddress, contractABI):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：验证地址格式
    IF NOT isValidAddress(contractAddress):
        THROW "无效的合约地址"
    END IF
    
    // 第3步：检查合约是否存在
    bytecode = client.getCode(contractAddress, "latest")
    IF bytecode.length == 0:
        THROW "指定地址没有部署合约"
    END IF
    
    // 第4步：解析合约ABI
    parsedABI = parseABI(contractABI)
    IF parsedABI == NULL:
        THROW "ABI解析失败"
    END IF
    
    // 第5步：创建合约实例
    contract = createContractInstance({
        address: contractAddress,
        abi: parsedABI,
        client: client
    })
    
    // 第6步：验证合约接口（可选）
    IF hasFunction(parsedABI, "supportsInterface"):
        // 检查ERC165接口支持
        interfaceSupport = contract.call("supportsInterface", ["0x01ffc9a7"])
        IF interfaceSupport:
            contract.supportsERC165 = true
        END IF
    END IF
    
    // 第7步：获取合约基本信息
    contractInfo = {
        address: contractAddress,
        bytecodeSize: bytecode.length,
        functions: extractFunctions(parsedABI),
        events: extractEvents(parsedABI),
        isVerified: true
    }
    
    RETURN {
        contract: contract,
        info: contractInfo
    }
END FUNCTION

// 辅助函数：提取函数列表
FUNCTION extractFunctions(abi):
    functions = []
    FOR each item IN abi:
        IF item.type == "function":
            functions.append({
                name: item.name,
                inputs: item.inputs,
                outputs: item.outputs,
                stateMutability: item.stateMutability
            })
        END IF
    END FOR
    RETURN functions
END FUNCTION

// 辅助函数：提取事件列表
FUNCTION extractEvents(abi):
    events = []
    FOR each item IN abi:
        IF item.type == "event":
            events.append({
                name: item.name,
                inputs: item.inputs,
                anonymous: item.anonymous
            })
        END IF
    END FOR
    RETURN events
END FUNCTION
```

### 关键步骤
1. **地址验证**：确保合约地址有效
2. **代码检查**：验证合约已部署
3. **ABI解析**：解析合约接口定义
4. **实例创建**：创建可调用的合约对象
5. **接口验证**：检查合约支持的标准

---

## 事件监听流程

### 流程概述
实时监听智能合约事件，需要WebSocket连接。

### 伪代码

```pseudocode
FUNCTION subscribeContractEvents(contractAddress, eventSignatures, callback):
    // 第1步：建立WebSocket连接
    client = connectToEthereumWS(WEBSOCKET_ENDPOINT)
    IF client == NULL:
        THROW "WebSocket连接失败"
    END IF
    
    // 第2步：创建事件过滤器
    topics = []
    FOR each eventSig IN eventSignatures:
        eventHash = keccak256(eventSig)
        topics.append(eventHash)
    END FOR
    
    filter = {
        addresses: [contractAddress],
        topics: [topics] // 第一层topics是事件签名
    }
    
    // 第3步：创建日志通道
    logChannel = createChannel()
    
    // 第4步：订阅日志
    subscription = client.subscribeFilterLogs(filter, logChannel)
    IF subscription == NULL:
        THROW "事件订阅失败"
    END IF
    
    // 第5步：启动监听循环
    WHILE true:
        SELECT:
            CASE error FROM subscription.errorChannel:
                PRINT "订阅错误: " + error
                BREAK
            
            CASE log FROM logChannel:
                // 第6步：解析事件日志
                eventData = parseEventLog(log)
                
                // 第7步：调用回调函数
                callback(eventData)
        END SELECT
    END WHILE
    
    // 第8步：清理资源
    subscription.unsubscribe()
    client.close()
END FUNCTION

// 事件日志解析
FUNCTION parseEventLog(log):
    RETURN {
        address: log.address,
        blockNumber: log.blockNumber,
        transactionHash: log.transactionHash,
        logIndex: log.logIndex,
        topics: log.topics,
        data: log.data,
        eventSignature: log.topics[0], // 第一个topic是事件签名
        timestamp: getCurrentTimestamp()
    }
END FUNCTION

// 示例回调函数
FUNCTION eventCallback(eventData):
    PRINT "收到事件:"
    PRINT "  合约地址: " + eventData.address
    PRINT "  区块号: " + eventData.blockNumber
    PRINT "  交易哈希: " + eventData.transactionHash
    PRINT "  事件数据: " + eventData.data
END FUNCTION
```

### 关键步骤
1. **WebSocket连接**：建立实时通信
2. **过滤器创建**：指定监听条件
3. **事件订阅**：注册事件监听器
4. **日志解析**：提取事件信息
5. **回调处理**：处理接收到的事件

---

## 历史事件查询流程

### 流程概述
查询指定区块范围内的历史事件日志。

### 伪代码

```pseudocode
FUNCTION queryHistoricalEvents(contractAddress, eventSignatures, fromBlock, toBlock):
    // 第1步：连接网络
    client = connectToEthereum(RPC_ENDPOINT)
    
    // 第2步：构建事件主题过滤器
    topics = []
    IF eventSignatures.length > 0:
        eventHashes = []
        FOR each eventSig IN eventSignatures:
            eventHash = keccak256(eventSig)
            eventHashes.append(eventHash)
        END FOR
        topics.append(eventHashes)
    END IF
    
    // 第3步：创建查询过滤器
    filter = {
        fromBlock: fromBlock,
        toBlock: toBlock,
        addresses: [contractAddress],
        topics: topics
    }
    
    // 第4步：执行查询
    logs = client.filterLogs(filter)
    IF logs == NULL:
        THROW "事件查询失败"
    END IF
    
    // 第5步：处理查询结果
    events = []
    FOR each log IN logs:
        eventData = {
            address: log.address,
            blockNumber: log.blockNumber,
            blockHash: log.blockHash,
            transactionHash: log.transactionHash,
            transactionIndex: log.transactionIndex,
            logIndex: log.logIndex,
            topics: log.topics,
            data: log.data,
            eventSignature: log.topics[0]
        }
        
        // 第6步：解码事件数据（如果有ABI）
        IF hasABI(contractAddress):
            decodedData = decodeEventData(log, getContractABI(contractAddress))
            eventData.decodedData = decodedData
        END IF
        
        events.append(eventData)
    END FOR
    
    // 第7步：按区块号排序
    events = sortByBlockNumber(events)
    
    // 第8步：统计信息
    statistics = {
        totalEvents: events.length,
        blockRange: {
            from: fromBlock,
            to: toBlock
        },
        contractAddress: contractAddress,
        uniqueTransactions: countUniqueTransactions(events)
    }
    
    RETURN {
        events: events,
        statistics: statistics
    }
END FUNCTION

// 辅助函数：解码事件数据
FUNCTION decodeEventData(log, contractABI):
    eventABI = findEventABI(contractABI, log.topics[0])
    IF eventABI == NULL:
        RETURN NULL
    END IF
    
    // 解码indexed参数（在topics中）
    indexedParams = []
    FOR i = 1 TO log.topics.length - 1:
        param = decodeIndexedParameter(log.topics[i], eventABI.inputs[i-1])
        indexedParams.append(param)
    END FOR
    
    // 解码非indexed参数（在data中）
    nonIndexedParams = decodeEventParameters(log.data, eventABI)
    
    RETURN {
        eventName: eventABI.name,
        indexedParams: indexedParams,
        nonIndexedParams: nonIndexedParams
    }
END FUNCTION
```

### 关键步骤
1. **过滤器构建**：设置查询条件
2. **区块范围限制**：避免查询范围过大
3. **日志获取**：批量查询历史日志
4. **数据解码**：解析事件参数
5. **结果排序**：按时间顺序整理

---

## 区块订阅流程

### 流程概述
实时订阅新区块，获取最新的区块信息。

### 伪代码

```pseudocode
FUNCTION subscribeNewBlocks(callback, includeTransactions = false):
    // 第1步：建立WebSocket连接
    client = connectToEthereumWS(WEBSOCKET_ENDPOINT)
    IF client == NULL:
        THROW "WebSocket连接失败"
    END IF
    
    // 第2步：创建区块头通道
    headerChannel = createChannel()
    
    // 第3步：订阅新区块头
    subscription = client.subscribeNewHead(headerChannel)
    IF subscription == NULL:
        THROW "区块订阅失败"
    END IF
    
    PRINT "开始监听新区块..."
    
    // 第4步：监听循环
    WHILE true:
        SELECT:
            CASE error FROM subscription.errorChannel:
                PRINT "订阅错误: " + error
                // 尝试重新连接
                subscription = reconnectAndSubscribe(client, headerChannel)
                IF subscription == NULL:
                    BREAK
                END IF
            
            CASE header FROM headerChannel:
                // 第5步：处理新区块头
                blockInfo = {
                    number: header.number,
                    hash: header.hash,
                    timestamp: header.timestamp,
                    difficulty: header.difficulty,
                    gasLimit: header.gasLimit,
                    gasUsed: header.gasUsed,
                    parentHash: header.parentHash
                }
                
                // 第6步：获取完整区块（如果需要）
                IF includeTransactions:
                    fullBlock = client.getBlockByHash(header.hash, true)
                    blockInfo.transactions = fullBlock.transactions
                    blockInfo.transactionCount = fullBlock.transactions.length
                END IF
                
                // 第7步：调用回调函数
                callback(blockInfo)
        END SELECT
    END WHILE
    
    // 第8步：清理资源
    subscription.unsubscribe()
    client.close()
END FUNCTION

// 重连机制
FUNCTION reconnectAndSubscribe(client, headerChannel):
    maxRetries = 3
    retryCount = 0
    
    WHILE retryCount < maxRetries:
        SLEEP(1000 * (retryCount + 1)) // 递增延迟
        
        // 重新连接
        client = connectToEthereumWS(WEBSOCKET_ENDPOINT)
        IF client != NULL:
            subscription = client.subscribeNewHead(headerChannel)
            IF subscription != NULL:
                PRINT "重连成功"
                RETURN subscription
            END IF
        END IF
        
        retryCount += 1
        PRINT "重连失败，重试中... (" + retryCount + "/" + maxRetries + ")"
    END WHILE
    
    PRINT "重连失败，停止监听"
    RETURN NULL
END FUNCTION

// 示例回调函数
FUNCTION blockCallback(blockInfo):
    PRINT "新区块到达:"
    PRINT "  区块号: " + blockInfo.number
    PRINT "  区块哈希: " + blockInfo.hash
    PRINT "  时间戳: " + blockInfo.timestamp
    PRINT "  Gas使用: " + blockInfo.gasUsed + "/" + blockInfo.gasLimit
    
    IF blockInfo.transactions != NULL:
        PRINT "  交易数量: " + blockInfo.transactionCount
    END IF
    
    PRINT "  " + "-".repeat(50)
END FUNCTION
```

### 关键步骤
1. **WebSocket连接**：建立实时通信
2. **区块头订阅**：监听新区块产生
3. **完整区块获取**：根据需要获取交易详情
4. **错误处理**：实现重连机制
5. **回调处理**：处理新区块信息

---

## 通用错误处理模式

### 网络错误处理

```pseudocode
FUNCTION handleNetworkError(operation, maxRetries = 3):
    retryCount = 0
    
    WHILE retryCount < maxRetries:
        TRY:
            result = operation()
            RETURN result
        CATCH NetworkError as e:
            retryCount += 1
            IF retryCount >= maxRetries:
                THROW "网络操作失败，已重试" + maxRetries + "次: " + e.message
            END IF
            
            delay = 1000 * retryCount // 递增延迟
            PRINT "网络错误，" + delay + "ms后重试... (" + retryCount + "/" + maxRetries + ")"
            SLEEP(delay)
        END TRY
    END WHILE
END FUNCTION
```

### 交易确认等待

```pseudocode
FUNCTION waitForConfirmation(client, txHash, confirmations = 1, timeout = 300000):
    startTime = getCurrentTime()
    
    WHILE getCurrentTime() - startTime < timeout:
        receipt = client.getTransactionReceipt(txHash)
        
        IF receipt != NULL:
            currentBlock = client.getBlockNumber()
            confirmationCount = currentBlock - receipt.blockNumber + 1
            
            IF confirmationCount >= confirmations:
                RETURN receipt
            END IF
            
            PRINT "等待确认... (" + confirmationCount + "/" + confirmations + ")"
        END IF
        
        SLEEP(1000) // 等待1秒
    END WHILE
    
    THROW "交易确认超时"
END FUNCTION
```

---

## 性能优化建议

### 1. 批量操作
- 使用批量RPC调用减少网络往返
- 合并多个查询到单个请求
- 使用连接池管理网络连接

### 2. 缓存策略
- 缓存不变的数据（如历史区块）
- 使用本地数据库存储常用信息
- 实现智能缓存失效机制

### 3. 错误恢复
- 实现指数退避重试策略
- 区分可重试和不可重试错误
- 提供优雅的降级方案

### 4. 资源管理
- 及时关闭WebSocket连接
- 限制并发请求数量
- 监控内存和CPU使用情况

---

*本文档提供了以太坊开发中各个模块的标准工作流程，可作为实现参考和代码审查的依据。*