package tests

/*
	Sorry, this is ugly, I know. I needed this one as the import directive
	on the _test.go files is olympically ignored by the go get command
*/

import (
	_ "gopkg.in/stretchr/testify.v1/assert"
)

func main() {
}
