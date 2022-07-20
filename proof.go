package main

import (
	"fmt"
	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/db/memdb"
	"github.com/cosmos/cosmos-sdk/store/v2alpha1/multi"

	types "github.com/cosmos/cosmos-sdk/store/v2alpha1"
	ics23types "github.com/cosmos/ibc-go/v3/modules/core/23-commitment/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func main() {
	db := memdb.NewDB()
	opts := multi.DefaultStoreConfig()
	err := opts.RegisterSubstore("store1", types.StoreTypePersistent)
	if err != nil {
		fmt.Println("ERROR: Substore cannot be registered: ", err)
		return
	}
	store, err := multi.NewStore(db, opts)
	if err != nil {
		fmt.Println("ERROR: multi store cannot be created: ", err)
		return
	}

	skey_1 := types.NewKVStoreKey("store1")
	substore := store.GetKVStore(skey_1)
	substore.Set([]byte("MYKEY"), []byte("MYVALUE"))
	cid := store.Commit()

	res := store.Query(abci.RequestQuery{
		Path:  "/store1/key", // required path to get key/value+proof
		Data:  []byte("MYKEY"),
		Prove: true,
	})

	proof, err := ics23types.ConvertProofs(res.ProofOps)
	if err != nil {
		fmt.Println("ERROR: Proof conversion failed: ", err)
		return
	}

	root := cid.Hash
	paths := []string{"store1", "MYKEY"}
	value := []byte("MYVALUE")

	substoreSpec := ics23.SmtSpec
	// This spec change is for the forked ics23 supporting SMT
	// Please edit `go.mod` to use the forked ics23 package
	// Set prehash for key comparisons (The value in the proof isn't hashed)
	substoreSpec.PrehashComparedKey = ics23.HashOp_SHA256

	specs := []*ics23.ProofSpec{substoreSpec, ics23.TendermintSpec}
	merkleRoot := ics23types.NewMerkleRoot(root)
	path := ics23types.NewMerklePath(paths...)

	err = proof.VerifyMembership(specs, &merkleRoot, path, value)
	if err == nil {
		fmt.Println("The existence proof was verified successfully")
	} else {
		fmt.Println("ERROR: The verification of the existence proof failed unexpectedly: ", err)
	}

	// test for non-existence proof
	res = store.Query(abci.RequestQuery{
		Path:  "/store1/key", // required path to get key/value+proof
		Data:  []byte("MYABSENTKEY"),
		Prove: true,
	})

	proof, err = ics23types.ConvertProofs(res.ProofOps)
	if err != nil {
		fmt.Println("ERROR: Proof conversion failed: ", err)
		return
	}

	paths = []string{"store1", "MYABSENTKEY"}
	path = ics23types.NewMerklePath(paths...)

	err = proof.VerifyNonMembership(specs, &merkleRoot, path)
	if err == nil {
		fmt.Println("The non-existence proof was verified successfully")
	} else {
		fmt.Println("ERROR: The verification of the non-existence proof failed unexpectedly: ", err)
	}
}
