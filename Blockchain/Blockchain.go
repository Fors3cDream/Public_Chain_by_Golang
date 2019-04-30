package Blockchain

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
)


// 数据库相关
const dbName = "ZDBlockchain.db"
const blockTableName = "zdBlocks"
const newestBlockKey = "zdNewestBlockKey"

type Blockchain struct {
	// 有序的区块数组
	// Blocks []*Block

	LatestHash [] byte // 最新区块的 Hash

	DB *bolt.DB // 存储区块的数据库
}

// 创建带有创世区块的区块链
func CreateBlockchainWithGensisBlock() *Blockchain {

	var blockchain *Blockchain

	//判断数据库是否存在
	if IsDBExists(dbName) {

		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			fmt.Println("打开数据库失败!")
			log.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(blockTableName))
			if b != nil {

				hash := b.Get([]byte(newestBlockKey))
				blockchain = &Blockchain{hash, db}
				fmt.Printf("%x", hash)
			}

			return nil
		})
		if err != nil {
			fmt.Println("打开表失败!")
			log.Panic(err)
		}

		return blockchain
	}

	//创建并打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		fmt.Println("创建数据库失败")
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		//blockTableName不存在再去创建表
		if b == nil {

			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				fmt.Println("创建表失败")
				log.Panic(err)
			}
		}

		if b != nil {

			//创世区块
			gensisBlock := CreateGensisBlock("Gensis Block...")
			//存入数据库
			err := b.Put(gensisBlock.Hash, gensisBlock.Serialize())
			if err != nil {
				fmt.Println("存入数据库失败")
				log.Panic(err)
			}

			//存储最新区块hash
			err = b.Put([]byte(newestBlockKey), gensisBlock.Hash)
			if err != nil {
				fmt.Println("存入最新区块失败")
				log.Panic(err)
			}

			blockchain = &Blockchain{gensisBlock.Hash, db}
		}

		return nil
	})
	//更新数据库失败
	if err != nil {
		fmt.Println("更新数据库失败")
		log.Fatal(err)
	}

	return blockchain
}

// 新增一个区块到区块链
func (blchain *Blockchain) AddBlockToBlockchain(data string) {

	err := blchain.DB.Update(func(tx *bolt.Tx) error {

		// 取表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// height, prevHash 都可以从数据库中获得
			blockBytes := b.Get(blchain.LatestHash)
			block := DeSerializeBlock(blockBytes)

			// 创建新区块
			newBlock := NewBlock(data, block.Height + 1, block.Hash)

			// 区块序列化入库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Fatal(err)
			}

			// 更新数据库里最新区块
			err = b.Put([]byte(newestBlockKey), newBlock.Hash)
			if err != nil {
				log.Fatal(err)
			}

			// 更新区块链最新区块
			blchain.LatestHash = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (blchain *Blockchain) Printchain1() {

	var block *Block

	// 当前遍历的区块 Hash
	var currHash []byte = blchain.LatestHash
	for {
		err := blchain.DB.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte(blockTableName))
			if b!=nil {
				blockBytes := b.Get(currHash)
				block = DeSerializeBlock(blockBytes)
				
				fmt.Printf("\n#####\nHeight:%d\nPrevHash:%x\nHash:%x\nData:%s\nTime:%s\nNonce:%d\n#####\n",
				block.Height, block.PrevBlockHash, block.Hash, block.Data, time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"), block.Nonce)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		//遍历到创世区块，跳出循环  创世区块哈希为0
		if big.NewInt(0).Cmp(&hashInt) == 0 {

			break
		}
		currHash = block.PrevBlockHash
	}
}

func (blchain *Blockchain) Printchain() {
	// 迭代器
	blcIterator := blchain.Interator()
	for {
		block := blcIterator.Next()
		fmt.Printf("\n####\nHeight:%d\nPrevHash:%x\nHash:%x\nData:%s\nTime:%s\nNonce:%d\n####\n",
					block.Height, block.PrevBlockHash, block.Hash, block.Data, time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"), block.Nonce)
					
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0{
			break
		}
	}
}

// 迭代器生成方法
func (blchain *Blockchain) Interator() *BlockchainIterator {
	return &BlockchainIterator{blchain.LatestHash, blchain.DB}
}

// 判断数据库是否存在
func IsDBExists(dbName string) bool {
	_, err := os.Stat(dbName)

	if err == nil {
		fmt.Println("数据库存在")
		return true
	}

	if os.IsNotExist(err) {
		fmt.Println("数据库不存在")
		return false
	}

	return true
}