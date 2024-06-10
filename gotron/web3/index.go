package web3

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	*ethclient.Client
	ctx context.Context
}

var Eth Client

// NewClient 创建并返回一个新的以太坊客户端
func init() {
	// rpc := "https://rpc-bsc.48.club"
	rpc := "https://api.trongrid.io/jsonrpc"
	// // rpc := "https://trx.mytokenpocket.vip/jsonrpc"
	// // rpc := "https://powerful-hardworking-panorama.tron-mainnet.quiknode.pro/93c8383167554e6531cf780c1d42679822142ca4/jsonrpc/"
	ctx := context.Background()
	client, err := ethclient.Dial(rpc)
	if err != nil {
		panic(err)
	}
	Eth = Client{client, ctx}
}

// GetCurrentBlockNumber 获取当前区块号
func (e *Client) GetBlockNumber() (uint64, error) {
	blockNumber, err := e.BlockNumber(e.ctx)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

// 获取合约日志
func (e *Client) Getlog(query ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := Eth.FilterLogs(e.ctx, query)
	if err != nil {

		return []types.Log{}, err
	}
	return logs, nil
}
