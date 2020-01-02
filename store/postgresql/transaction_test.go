package postgresql

import "testing"

func TestBeginTransactionCommit(t *testing.T) {
	testStr(t).TestBeginTransactionCommit()
}

func TestBeginTransactionAbort(t *testing.T) {
	testStr(t).TestBeginTransactionAbort()
}
