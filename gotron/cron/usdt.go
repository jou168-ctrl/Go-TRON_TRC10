package cron

import (
	"eeth/redis"
	tron "eeth/tron"
	web3 "eeth/web3"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// 初始化区块
var FromBlock uint64

func init() {
	Block, err := redis.Rd.Get("fromBlock")
	if err != nil {
		FromBlock, _ = web3.Eth.GetBlockNumber()
	}
	FromBlock, _ = strconv.ParseUint(Block, 10, 64)
	fmt.Println("初始化区块", FromBlock)
}

// 要监控得地址列表 放在Topics第三个参数
var addlist = []common.Hash{
	common.HexToHash(tron.TronToEeh("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t").Hex()),
}

// 锁
var mtx sync.Mutex

// 监控地址收款
func GetUsdtlog() {
	mtx.Lock()
	defer mtx.Unlock()
	block, err := web3.Eth.GetBlockNumber()
	if err != nil {
		fmt.Print("获取区块错误", err)
		return
	}
	if block <= FromBlock+10 {
		return
	}
	sum := block - FromBlock - 10
	if sum > 500 {
		sum = 500
	}
	fmt.Println("当前块", block, "FromBlock", FromBlock, "ToBlock", FromBlock+sum, "块差", sum)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(FromBlock)),
		ToBlock:   big.NewInt(int64(FromBlock + sum)),
		Addresses: []common.Address{tron.TronToEeh("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")}, //要监控的合约
		Topics: [][]common.Hash{
			{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}, //主题
			{},      //from  筛选转出地址 addlist 维护要监控得地址列表
			addlist, //to  筛选接收地址
		}}

	logs, err := web3.Eth.Getlog(query)
	if err != nil {
		fmt.Print("获取日志错误", err)
	}
	for _, vLog := range logs {
		from := tron.EthtoTron(common.HexToAddress(vLog.Topics[1].Hex()).String())
		value := big.NewInt(0).Div(big.NewInt(0).SetBytes(vLog.Data[:32]), big.NewInt(1000000)).Int64()
		hash := vLog.TxHash.Hex() //使用redis HSETNX //防止重复  只有能够被插入得哈希才处理 如果已经插入得就不处理
		fmt.Println("from", from)
		fmt.Println("value", value)
		fmt.Println("hash", hash)
	}
	FromBlock += sum + 1
	redis.Rd.Set("fromBlock", FromBlock, 0) //写入redis
}
