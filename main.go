package main

import (
	"./Blockchain"
)

func main() {

	// 创建带有创世区块的区块链
	blockchain := Blockchain.CreateBlockchainWithGensisBlock()
	defer blockchain.DB.Close()

	//添加一个新区快
	blockchain.AddBlockToBlockchain("First Block")
	blockchain.AddBlockToBlockchain("Second Block")
	blockchain.AddBlockToBlockchain("Third Block")

	blockchain.Printchain()
}