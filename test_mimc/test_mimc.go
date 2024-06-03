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
	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(r1cs)

	mimc := MIMC.NewMiMC()
	data := []byte("11")

	mimc.Write(data)
	test := mimc.Sum(nil)

	assignment := MyCircuit{witness{Hash_input: data}, statement{Hash_output: test}}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
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
	Hash_output frontend.Variable `gnark:",public"`
}

type witness struct {
	//statement
	Hash_input frontend.Variable
}

// Relation
func (circuit *MyCircuit) Define(api frontend.API) error {

	// tool
	mimc, _ := mimc.NewMiMC(api)

	// ... see Circuit API section

	for i := 0; i < 1000; i++ {
		mimc.Write(circuit.W.Hash_input)
		api.AssertIsEqual(circuit.S.Hash_output, mimc.Sum())
		mimc.Reset()
	}

	// ??
	return nil

}

// go mod tidy = 관련된 모듈 설치
