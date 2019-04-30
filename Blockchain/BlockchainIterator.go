package Blockchain

import (
	"github.com/boltdb/bolt"
	"log"
)

// 区块链迭代器结构
type BlockchainIterator struct {
	
	// 当前Hash
	CurrHash []byte
	// 区块链数据库
	DB *bolt.DB
}

func (blcIterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := blcIterator.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 获取当前迭代器对应的区块
			currBlockBytes := b.Get(blcIterator.CurrHash)
			block = DeSerializeBlock(currBlockBytes)

			// 更新迭代器
			blcIterator.CurrHash = block.PrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return block
}