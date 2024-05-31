package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	// MIMC "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	groth16 "github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	fmt.Println("Hello World")
	var Circuit MyCircuit
	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(r1cs)

	var cat1 fp.Element
	cat1.SetRandom()
	fmt.Println(cat1)
	var cat2 fp.Element
	cat2.SetRandom()
	fmt.Println(cat2)
	var cat3 fp.Element
	cat3.Mul(&cat1, &cat2)
	fmt.Println(cat3)

	assignment := MyCircuit{witness{A: cat1, B: cat2}, statement{C: cat3}}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	fmt.Println(witness)
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
	C fp.Element `gnark:",public"`
}

type witness struct {
	//statement
	A fp.Element
	B fp.Element
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
