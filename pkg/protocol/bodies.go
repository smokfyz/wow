package protocol

import (
	"errors"

	"github.com/smokfyz/wow/pkg/pow"
)

const (
	challengeRequest = iota
	verifyRequest
	challengeResponse
	verifiedResponse
	errorResponse
)

const (
	MaxBodySize = 256
)

var (
	ErrFailedToDecodeBody = errors.New("failed to decode body")
)

type ChallengeRequestBody struct{}

func (b ChallengeRequestBody) Encode() []byte {
	return []byte{challengeRequest}
}

type VerifyRequestBody struct {
	puzzle []byte
	nonce  []byte
}

func DecodeVerifyRequestBody(data []byte) (VerifyRequestBody, error) {
	if len(data) < pow.PuzzleSize+2 {
		return VerifyRequestBody{}, ErrFailedToDecodeBody
	}
	return VerifyRequestBody{data[1 : pow.PuzzleSize+1], data[pow.PuzzleSize+1:]}, nil
}

func (b VerifyRequestBody) Encode() []byte {
	payload := append(b.puzzle[:], b.nonce...)
	return append([]byte{verifyRequest}, payload...)
}

type ChallengeResponseBody struct {
	puzzle     []byte
	difficulty byte
}

func DecodeChallengeResponseBody(data []byte) (ChallengeResponseBody, error) {
	if len(data) < pow.PuzzleSize+2 {
		return ChallengeResponseBody{}, ErrFailedToDecodeBody
	}
	return ChallengeResponseBody{data[1 : pow.PuzzleSize+1], data[pow.PuzzleSize+1]}, nil
}

func (b ChallengeResponseBody) Encode() []byte {
	payload := append(b.puzzle[:], b.difficulty)
	return append([]byte{challengeResponse}, payload...)
}

type VerifiedResponseBody struct {
	result string
}

func DecodeVerifiedResponseBody(data []byte) (VerifiedResponseBody, error) {
	if len(data) < 2 {
		return VerifiedResponseBody{}, ErrFailedToDecodeBody
	}
	return VerifiedResponseBody{string(data[1:])}, nil
}

func (b VerifiedResponseBody) Encode() []byte {
	return append([]byte{verifiedResponse}, []byte(b.result)...)
}

type ErrorResponseBody struct {
	err string
}

func DecodeErrorResponseBody(data []byte) (ErrorResponseBody, error) {
	if len(data) < 2 {
		return ErrorResponseBody{}, ErrFailedToDecodeBody
	}
	return ErrorResponseBody{string(data[1:])}, nil
}

func (b ErrorResponseBody) Encode() []byte {
	return append([]byte{errorResponse}, []byte(b.err)...)
}
