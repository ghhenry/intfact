package intfact

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestEcParallel(t *testing.T) {
	// has a factor 59649589127497217
	n := intval("340282366920938463463374607431768211457")
	fac, err := EcParallel(context.Background(), &lcRandom{x: 10}, n, 10000, 215000, 200)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if fac == nil {
		t.Fatal("factor is nil")
	}
	if new(big.Int).Mod(n, fac).Sign() != 0 {
		t.Error("factor does not divide n:", fac)
	}
}

func TestEcParallel1(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	// has a factor 24844087002162188818957
	n := intval("51975989010489170378207397897900519326437583043996359")
	fac, err := EcParallel(context.Background(), &lcRandom{x: 10}, n, 20000, 550000, 500)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if fac == nil {
		t.Fatal("factor is nil")
	}
	if new(big.Int).Mod(n, fac).Sign() != 0 {
		t.Error("factor does not divide n:", fac)
	}
}

func TestEcParallel2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
	r := &lcRandom{x: 10}
	p1, _ := rand.Prime(r, 90)
	p2, _ := rand.Prime(r, 110)
	n := new(big.Int).Mul(p1, p2)
	t.Log(n, p1, p2)
	fac, err := EcParallel(context.Background(), r, n, 40000, 1400000, 3000)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if fac == nil {
		t.Fatal("factor is nil")
	}
	if new(big.Int).Mod(n, fac).Sign() != 0 {
		t.Error("factor does not divide n:", fac)
	}
}
