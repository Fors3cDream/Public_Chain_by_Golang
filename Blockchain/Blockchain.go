package Blockchain

type Blockchain struct {
	// 有序的区块数组
	Blocks []*Block
}

// 创建带有创世区块的区块链
func CreateBlockchainWithGensisBlock() *Blockchain {

	gensisBlock := CreateGensisBlock("Gensis Block!")

	return &Blockchain{[] *Block{gensisBlock}}
}

// 新增一个区块到区块链
func (blchain *Blockchain) AddBlockToBlockchain(data string, height int64, prevHash []byte) {

	// 新增区块
	newBlock := NewBlock(data, height, prevHash)

	// 添加到链上
	blchain.Blocks = append(blchain.Blocks, newBlock)
}