package db

import (
	"fmt"
	"testing"
)

func TestReTry(t *testing.T) {
	if e := DB.DB().Ping(); e != nil {
		panic(e)
	}
	fmt.Println("ping success")
}
