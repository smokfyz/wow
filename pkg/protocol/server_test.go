package protocol

import (
	"bytes"
	"testing"
	"time"

	"github.com/smokfyz/wow/pkg/pow"
)

type ChallengeMock struct{}

func (c ChallengeMock) GetDifficulty() byte {
	return 3
}

func (c ChallengeMock) GeneratePuzzle() []byte {
	puzzle := make([]byte, pow.PuzzleSize)
	return puzzle
}

func (c ChallengeMock) Verify(puzzle []byte, nonce []byte) bool {
	return true
}

func TestServerHandleRequestGeneratePuzzle(t *testing.T) {
	server := NewServer(ChallengeMock{}, 10*time.Second)
	body := ChallengeRequestBody{}
	respRaw := server.HandleRequest(body.Encode())
	expected := ChallengeResponseBody{make([]byte, pow.PuzzleSize), 3}.Encode()
	if !bytes.Equal(respRaw, expected) {
		t.Errorf("expected %v, got %v", expected, respRaw)
	}
}

func TestServerHandleRequestVerify(t *testing.T) {
	server := NewServer(ChallengeMock{}, 10*time.Second)
	body := ChallengeRequestBody{}
	respRaw := server.HandleRequest(body.Encode())
	bodyRaw := VerifyRequestBody{respRaw[1 : len(respRaw)-1], []byte{0x03, 0x04}}
	respRaw = server.HandleRequest(bodyRaw.Encode())
	if respRaw[0] != verifiedResponse {
		t.Errorf("expected %v, got %v", verifiedResponse, respRaw)
	}
}

func TestServerHandleRequestVerifyDecodeError(t *testing.T) {
	server := NewServer(ChallengeMock{}, 10*time.Second)
	body := VerifyRequestBody{[]byte{0x00, 0x01, 0x02, 0x03}, []byte{0x04, 0x05}}
	respRaw := server.HandleRequest(body.Encode())
	if respRaw[0] != errorResponse {
		t.Errorf("expected %v, got %v", errorResponse, respRaw)
	}
}

func TestServerHandleRequestUnexpectedType(t *testing.T) {
	server := NewServer(ChallengeMock{}, 10*time.Second)
	respRaw := server.HandleRequest([]byte{0x06})
	if respRaw[0] != errorResponse {
		t.Errorf("expected %v, got %v", errorResponse, respRaw)
	}
}
