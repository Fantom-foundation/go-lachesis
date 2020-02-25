package main

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type staker struct {
	Status           *big.Int
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	StakeAmount      *big.Int
	PaidUntilEpoch   *big.Int
	DelegatedMe      *big.Int
	DagAddress       common.Address
	SfcAddress       common.Address
}

type createdStaker struct {
	Address common.Address
	StakerId *big.Int
}

type createdStakers []createdStaker
//`{"address":"239fa7623354ec26520de878b52f13fe84b06971","crypto":{"cipher":"aes-128-ctr","ciphertext":"d25a3ce3381aef33d8c8e6345e3dc5547514cb4fe3daa78f1bb6a12b1b3a6400","cipherparams":{"iv":"0cdb328cb3b7d71b09efa90c90b23157"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"fb9296bf20661d0c4774514d914b1b0e33da86a616d179b604d06ca7196c6652"},"mac":"6396b71f2dabe1ae2d165692f7ff3f71a71a18ad5a71fbcd17a1f6d92cae15a6"},"id":"d2d05c35-66a9-4972-8216-c5a434cc72ff","version":3}`

// probably replace with a built one
type JsonKey struct {
	Address string `json:"address"`
	Crypto struct{
		Cipher string `json:"cipher"`
		CipherText string `json:"ciphertext"`
		CipherParams struct{
			Iv string `json:"iv"`
			Kdf string `json:"kdf"`
			KdfParams struct{
				DkLen int `json:"dklen"`
				N int `json:"n"`
				P int `json:"p"`
				R int `json:"r"`
				Salt string `json:"salt"`
			} `json:"kdfparams"`
			Mac string `json:"mac"`
		} `json:"cipherparams"`
		Id string `json:"id"`
		Version int `json:"version"`
	} `json:"crypto"`
}
