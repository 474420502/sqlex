package sqlex

import (
	"testing"
)

func Test(t *testing.T) {
	s := "insert into ?t<name>[1](a,b,name,field) value(?v)"
	check(&s)
}
