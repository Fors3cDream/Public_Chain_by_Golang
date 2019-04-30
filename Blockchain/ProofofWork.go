package Blockchain

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	//"fmt"
)

// 达标条件 -- 前多少位为0
const conditionBits = 20

type ProofofWork struct {

	// 待求值的 Block
	Block *Block

	// 工作量难度 big.Int大数存储
	condition *big.Int
}

// 新工作量证明
func NewProofofWork(block *Block) *ProofofWork {
	/**
	condition计算方式  假设：Hash为10位，conditionBits为2位
	eg:000000 0001(10位的Hash)
	1.10-2 = 8 将上值左移8位
	2.00 0000 0001 << 8 = 01 0000 0000 = condition
	3.只要计算的Hash满足 ：hash < condition，便是符合POW的哈希值
	*/

	// 创建一个初始值为1的condition
	condition := big.NewInt(1)

	// 左移 bits(Hash) - conditionBits 位
	condition = condition.Lsh(condition, 256-conditionBits)

	return &ProofofWork{block, condition}
}

// 拼接区块中的所有属性，返回字节数组
func (pow *ProofofWork) prepareData(nonce int) []byte{

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(conditionBits)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		}, []byte{},
	)

	return data
}

// 判断当前区块是否有效
func (proofofWork * ProofofWork) IsValid() bool {

	// 比较当前区块的哈希值与目标哈希值
	var hashInt big.Int
	hashInt.SetBytes(proofofWork.Block.Hash)

	if proofofWork.condition.Cmp(&hashInt) == 1{
		return true
	}
	return false
}

// 运行工作量证明函数
func (proofofWork *ProofofWork) Start() ([]byte, int64) {

	//1.将Block属性拼接成字节数组

	//2.生成hash
	//3.判断Hash值有效性，如果满足条件跳出循环

	//用于寻找目标hash值的随机数
	nonce := 0
	
	// 存储新生成的Hash值
	var hashInt big.Int
	var hash [32]byte

	for {
		// 准备数据
		dataBytes := proofofWork.prepareData(nonce)

		// Hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("努力计算中，当前Hash为: 0x%x\n", hash)

		//存储Hash到hashInt
		hashInt.SetBytes(hash[:])
		//验证Hash
		if proofofWork.condition.Cmp(&hashInt) == 1 {
			//fmt.Printf("计算完毕! Nonce 为: %d\n", nonce)
			break
		}
		nonce++
	}

	return hash[:], int64(nonce)
}