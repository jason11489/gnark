package main

import (
	"fmt"
	"math/big"
	"reflect"

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

	// var cat1 fp.Element
	// cat1.SetRandom()
	// // fmt.Println(cat1)
	// var cat2 fp.Element
	// cat2.SetRandom()
	// // fmt.Println(cat2)
	// var cat3 fp.Element
	// cat3.Mul(&cat1, &cat2)
	// // fmt.Println(cat3)

	aElement := fp.NewElement(2)
	bElement := fp.NewElement(3)
	cElement := fp.NewElement(6)

	aValue := new(big.Int)
	bValue := new(big.Int)
	cValue := new(big.Int)
	aElement.BigInt(aValue)
	bElement.BigInt(bValue)
	cElement.BigInt(cValue)

	assignment := MyCircuit{witness{A: aValue, B: bValue}, statement{C: cValue}}
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

	// relation
	fmt.Println("Relation start")
	// tt := Bn254.NewElement(circuit.W.A)

	// teat := fp.NewElement(circuit.W.A)
	var A fp.Element
	A.SetInterface(circuit.W.A)
	var B fp.Element
	B.SetInterface(circuit.W.B)
	var C fp.Element
	C.SetInterface(circuit.S.C)

	fmt.Printf("t1: %s\n", reflect.TypeOf(circuit.W.A))

	Check_C := api.Mul(A, B)
	api.AssertIsEqual(C, Check_C)

	api.Println(Check_C)
	// api.Println(circuit.W.B)

	fmt.Println("Relation finish")
	// ??
	return nil

}

// go mod tidy = 관련된 모듈 설치
