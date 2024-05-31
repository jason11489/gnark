package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"

	MIMC "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

func main() {
	fmt.Println("Hello World")
	var Circuit MyCircuit
	r1cs, _ := frontend.Compile(*ecc.BN254, r1cs.NewBuilder, &Circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(r1cs)

	mimc := MIMC.NewMiMC()
	data := []byte("11")

	// var arr [3]byte
	// for i := 0; i < arr.Len(); i++ {
	// 	arr.Push(data)
	// }

	mimc.Write(data)
	test := mimc.Sum(nil)

	assignment := MyCircuit{witness{A: 2, B: 3, Hash_input: data}, statement{C: 6, Hash_output: test}}
	witness, _ := frontend.NewWitness(&assignment, *ecc.BN254)
	publicWitness, _ := witness.Public()

	proof, _ := groth16.Prove(r1cs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)

}

// statement=Y as public, witness=C
type MyCircuit struct {
	W witness
	S statement
}

type statement struct {
	C           frontend.Variable `gnark:",public"`
	Hash_output frontend.Variable `gnark:",public"`
}

type witness struct {
	//statement
	A          frontend.Variable
	B          frontend.Variable
	Hash_input frontend.Variable
}

// Relation
func (circuit *MyCircuit) Define(api frontend.API) error {

	// tool
	mimc, _ := mimc.NewMiMC(api)

	// ... see Circuit API section
	Check_C := api.Mul(circuit.W.A, circuit.W.B)
	api.AssertIsEqual(circuit.S.C, Check_C)

	for i := 0; i < 1000; i++ {
		mimc.Write(circuit.W.Hash_input)
		api.AssertIsEqual(circuit.S.Hash_output, mimc.Sum())
		mimc.Reset()
	}

	// ??
	return nil

}

// go mod tidy = 관련된 모듈 설치
