package main

import (
	"./Blockchain"
	"fmt"
)

func main() {

	// 创建带有创世区块的区块链
	blockchain := Blockchain.CreateBlockchainWithGensisBlock()

	// 添加新区块
	blockchain.AddBlockToBlockchain("first Block",
		blockchain.Blocks[len(blockchain.Blocks) - 1].Height,
		blockchain.Blocks[len(blockchain.Blocks) - 1].Hash)
	
	blockchain.AddBlockToBlockchain("second Block",
		blockchain.Blocks[len(blockchain.Blocks) - 1].Height,
		blockchain.Blocks[len(blockchain.Blocks) - 1].Hash)

	blockchain.AddBlockToBlockchain("third Block",
		blockchain.Blocks[len(blockchain.Blocks) - 1].Height,
		blockchain.Blocks[len(blockchain.Blocks) - 1].Hash)
	

	fmt.Println(blockchain)
}