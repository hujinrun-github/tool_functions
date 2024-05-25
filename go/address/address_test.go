package address

import (
	"fmt"
	"testing"
)

func TestIp(t *testing.T) {
	res, err := GetLocalAddressByInterface("en0")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("res: %s\n", res)
}
