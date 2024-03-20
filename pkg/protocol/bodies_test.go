package protocol

import (
	"bytes"
	"testing"

	"github.com/smokfyz/wow/pkg/pow"
)

func TestVerifyRequestBody(t *testing.T) {
	t.Run("DecodeVerifyRequestBody", func(t *testing.T) {
		puzzle := make([]byte, pow.PuzzleSize)
		nonce := []byte{0x00}
		verifyRequestBody := VerifyRequestBody{puzzle, nonce}
		body, err := DecodeVerifyRequestBody(verifyRequestBody.Encode())
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !bytes.Equal(body.puzzle, puzzle) || !bytes.Equal(body.nonce, nonce) {
			t.Errorf("unexpected verify request body: %x", body)
		}
	})
}

// TODO: add tests for other bodies
