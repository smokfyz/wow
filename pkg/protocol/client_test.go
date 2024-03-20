package protocol

import (
	"testing"

	"github.com/smokfyz/wow/pkg/pow"
)

type VerifierMock struct{}

func (v VerifierMock) SetDifficulty(difficulty byte) error {
	return nil
}

func (v VerifierMock) Solve(puzzle []byte) []byte {
	return []byte{0x00, 0x01, 0x02}
}

func TestClientGetChallengeRequest(t *testing.T) {
	client := NewClient(VerifierMock{})
	challengeRequestBody := client.GetChallengeRequest()
	if len(challengeRequestBody) != 1 && challengeRequestBody[0] != 0x00 {
		t.Errorf("unexpected challenge request body: %x", challengeRequestBody)
	}
}

func TestClientHandleChallengeResponse(t *testing.T) {
	client := NewClient(VerifierMock{})
	puzzle := make([]byte, pow.PuzzleSize)
	challengeResponseBody := ChallengeResponseBody{puzzle, 3}
	rawChallengeResponseBody := challengeResponseBody.Encode()
	verifiedRequestBody, err := client.HandleChallengeResponse(rawChallengeResponseBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if verifiedRequestBody[0] != verifyRequest {
		t.Errorf("unexpected verified request: %x", verifiedRequestBody)
	}
}

func TestClientHandleVerifiedResponse(t *testing.T) {
	client := NewClient(VerifierMock{})
	verifiedResponseBody := VerifiedResponseBody{"result"}
	rawVerifiedResponseBody := verifiedResponseBody.Encode()
	err := client.HandleVerifiedResponse(rawVerifiedResponseBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestClientHandleErrorResponse(t *testing.T) {
	client := NewClient(VerifierMock{})
	errorResponseBody := ErrorResponseBody{"error"}
	rawErrorResponseBody := errorResponseBody.Encode()
	err := client.HandleErrorResponse(rawErrorResponseBody)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestClientIsErrorResponse(t *testing.T) {
	client := NewClient(VerifierMock{})
	errorResponseBody := ErrorResponseBody{"error"}
	rawErrorResponseBody := errorResponseBody.Encode()
	if !client.IsErrorResponse(rawErrorResponseBody) {
		t.Errorf("expected error response")
	}
}
