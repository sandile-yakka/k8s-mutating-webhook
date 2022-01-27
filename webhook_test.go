package main

import (
	"fmt"
	server "mutating-webhook/server"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatch(t *testing.T) {

	patch := server.PatchResourceRequests()
	// ResourceLI
	fmt.Printf("====== %+v", patch)

	require.NotNil(t, patch)
}
