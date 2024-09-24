package users

import "testing"

func TestHashPwd(t *testing.T) {
	p, _ := HashPwd("cvWH+Iz7xMXcKd8d3Uvy5VJnnM4ex6ml")
	t.Log(p)
}
