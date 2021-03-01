package eth

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func (e *Eth) SendRawTransaction(
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (common.Hash, error) {
	nonce, err := e.GetNonce(e.address, nil)
	var hash common.Hash
	if err != nil {
		return hash, err
	}
	// fmt.Printf("nonce %v\n", nonce)

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	// fmt.Println(tx)
	signedTx, err := eTypes.SignTx(tx, eTypes.NewEIP155Signer(e.chainId), e.privateKey)
	if err != nil {
		return hash, err
	}
	// fmt.Println("signTx", signedTx)
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return hash, err
	}
	// fmt.Printf("serializedTx 0x%x\n", serializedTx)

	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	return hash, err

}

func (e *Eth) SyncSendRawTransaction(
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	data []byte,
) (*eTypes.Receipt, error) {
	nonce, err := e.GetNonce(e.address, nil)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("nonce %v\n", nonce)

	tx := eTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	// fmt.Println(tx)
	signedTx, err := eTypes.SignTx(tx, eTypes.NewEIP155Signer(e.chainId), e.privateKey)
	if err != nil {
		return nil, err
	}
	// fmt.Println("signTx", signedTx)
	serializedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("serializedTx 0x%x\n", serializedTx)
	var hash common.Hash
	err = e.c.Call("eth_sendRawTransaction", &hash, fmt.Sprintf("0x%x", serializedTx))
	if err != nil {
		return nil, err
	}

	fmt.Printf("hash %v\n", hash)

	for {
		receipt, err := e.GetTransactionReceipt(hash)
		if err != nil {
			if err.Error() != "not found" {
				return nil, err
			}
		}
		if receipt != nil {
			return receipt, nil
		}
	}
}
