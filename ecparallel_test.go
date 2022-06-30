package intfact

import (
	"context"
	"fmt"
	"testing"
)

func TestEcParallel(t *testing.T) {
	fac, err := EcParallel(context.Background(), &lcRandom{x: 10}, intval("340282366920938463463374607431768211457"), 20000, 550000)
	fmt.Println(fac, err)
}
