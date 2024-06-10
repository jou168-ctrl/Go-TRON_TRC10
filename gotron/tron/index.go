package tron

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// 将波场地址转换成以太坊地址
func TronToEeh(TronAddress string) common.Address {
	addrByte := hex.EncodeToString(base58Decode([]byte(TronAddress)))
	address := strings.Replace(addrByte, "41", "0x", 1) //转换后 41替换成0x
	return common.HexToAddress(address)
}

// 将以太坊地址转换成波场地址
func EthtoTron(EthAddress string) string {
	EthAddress = strings.Replace(EthAddress, "0x", "41", 1) // 0x替换成41
	addrByte, err := hex.DecodeString(EthAddress)
	if err != nil {
		return ""
	}

	sha := sha256.New()
	sha.Write(addrByte)
	shaStr := sha.Sum(nil)

	sha2 := sha256.New()
	sha2.Write(shaStr)
	shaStr2 := sha2.Sum(nil)

	addrByte = append(addrByte, shaStr2[:4]...)
	return string(base58Encode(addrByte))
}

var base58Alphabets = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func base58Encode(input []byte) []byte {
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := &big.Int{}
	var result []byte
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabets[mod.Int64()])
	}
	reverseBytes(result)
	return result
}

func base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	for _, b := range input {
		charIndex := bytes.IndexByte(base58Alphabets, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	if input[0] == base58Alphabets[0] {
		decoded = append([]byte{0x00}, decoded...)
	}
	return decoded[:len(decoded)-4]
}

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
