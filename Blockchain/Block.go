package Blockchain

import (
	"time"
	//"strconv"
	"fmt"
	"bytes"
	//"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	// 区块高度
	Height int64
	// 上一个区块的 hash
	PrevBlockHash []byte
	// 交易数据 - 交易池
	Data []byte
	// 时间戳
	Timestamp int64
	// 区块hash
	Hash []byte
	// Nonce - 工作量证明
	Nonce int64
}

//区块序列化
func (block *Block) Serialize() []byte  {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil{

		log.Panic(err)
	}

	return result.Bytes()
}

//区块反序列化
func DeSerializeBlock(blockBytes []byte) *Block  {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)
	if err != nil {

		log.Panic(err)
	}

	return &block
}

// // 计算当前区块hash
// func (block *Block) Hash_func() {
// 	// 将高度、时间戳转为字节数组
// 	heightBytes := IntToHex(block.Height)

// 	timestampStr := strconv.FormatInt(block.Timestamp, 2) // base 2 - 2进制形式
// 	timestamp := []byte(timestampStr)

// 	// 拼接所有属性，便于计算hash
// 	blockBytes := bytes.Join([][]byte{
// 		heightBytes,
// 		block.PrevBlockHash,
// 		block.Data,
// 		timestamp,
// 		block.Hash}, []byte{})

// 	// 计算hash
// 	hash := sha256.Sum256(blockBytes)

// 	fmt.Println(hash)

// 	block.Hash = hash[:]
// }

// 创建新的区块
func NewBlock(data string, height int64, PrevBlockHash []byte) *Block {

	// 创建区块
	block := &Block{
		Height:		height,
		PrevBlockHash: PrevBlockHash,
		Data:	[]byte(data),
		Timestamp: time.Now().Unix(),
		Hash:	nil,
		Nonce:  0}

	// 调用工作量证明返回有效的hash
	pow := NewProofofWork(block)
	hash, nonce := pow.Start()
	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Printf("\r############Nonce: %d - Hash: %x\n", nonce, hash)

	// block.Hash_func()

	return block
}

// 创世区块的生成
func CreateGensisBlock(data string) *Block {
	fmt.Println("Create Gensisblock")
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}