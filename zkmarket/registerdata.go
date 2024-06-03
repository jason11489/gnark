package main

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	// MIMC "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	groth16 "github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/emulated"
)

func main() {
	var Circuit MyCircuit[emulated.BN254Fp]
	r1cs,
		_ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit)

	// r1cs := groth16.NewCS(ecc.BN254)

	fmt.Println(ecc.BN254.ScalarField())
	fmt.Println(ecc.BN254.BaseField())

	// groth16 zkSNARK: Setup
	pk,
		vk,
		_ := groth16.Setup(r1cs)

	aElement := fp.NewElement(2)
	bElement := fp.NewElement(3)
	cElement := fp.NewElement(5)

	aValue := new(big.Int)
	bValue := new(big.Int)
	cValue := new(big.Int)
	aElement.BigInt(aValue)
	bElement.BigInt(bValue)
	cElement.BigInt(cValue)

	assignment := MyCircuit[emulated.BN254Fp]{
		witness[emulated.BN254Fp]{
			A: emulated.ValueOf[emulated.BN254Fp](aValue),
			B: emulated.ValueOf[emulated.BN254Fp](bValue),
		},
		statement[emulated.BN254Fp]{
			C: emulated.ValueOf[emulated.BN254Fp](cValue),
		},
	}
	witness,
		_ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness,
		_ := witness.Public()

	proof,
		_ := groth16.Prove(r1cs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)

}

// statement=Y as public, witness=C
type MyCircuit[T emulated.FieldParams] struct {
	W witness[T]
	S statement[T]
}

type statement[T emulated.FieldParams] struct {
	h_ct emulated.Element[T] `gnark:",public"`
	h_k  emulated.Element[T] `gnark:",public"`
	ENA  emulated.Element[T] `gnark:",public"`
}

type witness[T emulated.FieldParams] struct {
	//statement
	data    []emulated.Element[T]
	CT_data []emulated.Element[T]
	k_data  emulated.Element[T]
}

// Relation
func (circuit *MyCircuit[T]) Define(api frontend.API) error {

	// tool relation

	Fp_api,
		_ := emulated.NewField[T](api)
	A := Fp_api.NewElement(circuit.W.A)
	B := Fp_api.NewElement(circuit.W.B)
	C := Fp_api.NewElement(circuit.S.C)
	api.Println(A)

	for i := 0; i < 10; i++ {
		Check_C := Fp_api.Add(A, B)
		Fp_api.AssertIsEqual(C, Check_C)
	}

	return nil

}

// go mod tidy = 관련된 모듈 설치
