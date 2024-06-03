package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	fmt.Println("Hello World")
	var Circuit MyCircuit
	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(r1cs)

	assignment := MyCircuit{witness{A: 2, B: 3}, statement{C: 6}}
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
	C frontend.Variable `gnark:",public"`
}

type witness struct {
	//statement
	A frontend.Variable
	B frontend.Variable
}

// Relation
func (circuit *MyCircuit) Define(api frontend.API) error {

	// tool

	// ... see Circuit API section
	Check_C := api.Mul(circuit.W.A, circuit.W.B)
	api.AssertIsEqual(circuit.S.C, Check_C)

	// ??
	return nil

}

// go mod tidy = 관련된 모듈 설치
