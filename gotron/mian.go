//package main

// import (
// 	cron "eeth/cron"
// )

// func main() {
// 	// 启动定时任务
// 	cron.Run()

//		select {}
//	}
package main

import (
	"crypto/ecdsa"
	"eeth/tron"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func endsWithFiveSameChars(str string) bool {

	// 取字符串最后一个字符
	lastChar := str[len(str)-1]

	// 检查字符串最后5个字符是否都与最后一个字符相同
	for i := len(str) - 2; i >= len(str)-7; i-- {
		if str[i] != lastChar {
			return false
		}

	}

	// 所有字符都相同
	return true
}

func generateAddress(m *sync.WaitGroup) {
	defer m.Done()
	for {

		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		tronaddress := tron.EthtoTron(address)

		appkys := hexutil.Encode(privateKeyBytes)[2:]
		if endsWithFiveSameChars(tronaddress) {
			apps := tronaddress[33]
			s2 := string(apps)
			str := s2 + "wS2" + ".txt"
			file, errs := os.OpenFile(str, os.O_CREATE, 0)
			if errs != nil {
				log.Fatal(errs)
			}
			count, err := file.WriteString(tronaddress + ":kys:" + appkys)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("写入了长度为 %d 的字符串\n", count)
			fmt.Println(tronaddress)
			fmt.Println(appkys)
			break

		}

	}
}

func main() {

	numCPU := runtime.NumCPU()
	fmt.Println("Number of CPUs:", numCPU)

	// 设置 GOMAXPROCS 为 CPU 核心数
	runtime.GOMAXPROCS(numCPU)
	var m sync.WaitGroup
	//启动多个协程

	for i := 0; i < numCPU; i++ {
		m.Add(1)
		go generateAddress(&m)

	}
	m.Wait()
}
