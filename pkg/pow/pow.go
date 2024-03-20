// Package implements a proof of work algorithm based on the SHA256 hash function.
// The algorithm generates a random puzzle and then tries to find a nonce that when combined with the puzzle
// produces a hash with a certain number of leading zeros. The number of leading zeros is determined by the difficulty
// level. The difficulty level must be between MinDifficulty and MaxDifficulty.
package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
)

const (
	MinDifficulty = 1
	MaxDifficulty = 8
	PuzzleSize    = 32
)

var (
	ErrInvalidDifficulty = fmt.Errorf("difficulty must be between %d and %d", MinDifficulty, MaxDifficulty)
)

type Challenge struct {
	difficulty byte
	hasher     hash.Hash
}

func NewChallenge(difficulty byte) (*Challenge, error) {
	if difficulty < MinDifficulty || difficulty > MaxDifficulty {
		return nil, ErrInvalidDifficulty
	}
	return &Challenge{difficulty, sha256.New()}, nil
}

func (c Challenge) GetDifficulty() byte {
	return c.difficulty
}

func (c *Challenge) SetDifficulty(difficulty byte) error {
	if difficulty < MinDifficulty || difficulty > MaxDifficulty {
		return ErrInvalidDifficulty
	}
	c.difficulty = difficulty
	return nil
}

// GeneratePuzzle generates a random puzzle. The puzzle is a byte slice of size PuzzleSize.
func (c Challenge) GeneratePuzzle() []byte {
	randBytes := make([]byte, PuzzleSize)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err) // really unexpected error
	}
	return randBytes
}

// Solve tries to find a nonce that when combined with the puzzle produces a hash with a certain number of leading zeros.
func (c Challenge) Solve(puzzle []byte) []byte {
	nonce := make([]byte, 1)
	for !c.Verify(puzzle, nonce) {
		nonce = incrementBytes(nonce)
	}
	return nonce
}

// Verify checks if the given nonce produces a hash with a certain number of leading zeros.
func (c Challenge) Verify(puzzle []byte, nonce []byte) bool {
	c.hasher.Reset()
	c.hasher.Write(append(puzzle, nonce...))
	return checkLeadingZeros(c.hasher.Sum(nil), c.difficulty)
}

func checkLeadingZeros(b []byte, n byte) bool {
	for i := range n {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func incrementBytes(b []byte) []byte {
	for i := 0; i <= len(b); i++ {
		if i == len(b) {
			return append(b, 1)
		}

		if b[i] == 255 {
			b[i] = 0
		} else {
			b[i]++
			break
		}
	}
	return b
}
