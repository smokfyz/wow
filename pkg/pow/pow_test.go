package pow_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/smokfyz/wow/pkg/pow"
)

func TestNewChallenge(t *testing.T) {
	_, err := pow.NewChallenge(3)
	if err != nil {
		t.Errorf("NewChallenge(3) failed: %v", err)
	}
}

func TestNewChallengeUnsupportedDifficulty(t *testing.T) {
	difficulties := []byte{0, 10, 20}
	for _, difficulty := range difficulties {
		_, err := pow.NewChallenge(difficulty)
		if err == nil {
			t.Errorf("NewChallenge(%d) should fail", difficulty)
		}
		if !errors.Is(err, pow.ErrInvalidDifficulty) {
			t.Errorf("NewChallenge(%d) failed with unexpected error: %v", difficulty, err)
		}
	}
}

func TestGeneratePuzzle(t *testing.T) {
	c, _ := pow.NewChallenge(3)
	puzzle := c.GeneratePuzzle()
	if len(puzzle) != pow.PuzzleSize {
		t.Errorf("GeneratePuzzle() failed: unexpected puzzle size")
	}
}

func TestSolve(t *testing.T) {
	c, _ := pow.NewChallenge(3)

	testCases := []struct {
		puzzle []byte
		nonce  []byte
	}{
		{
			[]byte{6, 75, 143, 78, 73, 225, 255, 145, 249, 105, 157, 49, 93, 48, 16, 101, 138, 195, 136, 156, 231, 120, 55, 217, 23, 15, 111, 193, 176, 29, 33, 71},
			[]byte{57, 16, 14, 1},
		},
		{
			[]byte{99, 251, 192, 99, 248, 247, 100, 73, 2, 51, 35, 220, 134, 146, 103, 160, 100, 108, 40, 175, 191, 249, 252, 101, 34, 236, 208, 177, 206, 101, 198, 108},
			[]byte{123, 100, 112},
		},
		{
			[]byte{173, 116, 159, 66, 187, 198, 82, 19, 215, 120, 36, 209, 240, 88, 87, 151, 90, 18, 27, 47, 47, 164, 200, 55, 240, 43, 172, 181, 94, 3, 104, 138},
			[]byte{94, 235, 72, 1},
		},
	}

	for _, tc := range testCases {
		nonce := c.Solve(tc.puzzle)
		if !bytes.Equal(nonce, tc.nonce) {
			t.Errorf("Solve(%v) failed: expected %v, got %v", tc.puzzle, tc.nonce, nonce)
		}
	}
}

func TestVerify(t *testing.T) {
	c, _ := pow.NewChallenge(3)

	testCases := []struct {
		puzzle []byte
		nonce  []byte
		result bool
	}{
		{
			[]byte{6, 75, 143, 78, 73, 225, 255, 145, 249, 105, 157, 49, 93, 48, 16, 101, 138, 195, 136, 156, 231, 120, 55, 217, 23, 15, 111, 193, 176, 29, 33, 71},
			[]byte{57, 16, 14, 1},
			true,
		},
		{
			[]byte{99, 251, 192, 99, 248, 247, 100, 73, 2, 51, 35, 220, 134, 146, 103, 160, 100, 108, 40, 175, 191, 249, 252, 101, 34, 236, 208, 177, 206, 101, 198, 108},
			[]byte{123, 100, 112},
			true,
		},
		{
			[]byte{173, 116, 159, 66, 187, 198, 82, 19, 215, 120, 36, 209, 240, 88, 87, 151, 90, 18, 27, 47, 47, 164, 200, 55, 240, 43, 172, 181, 94, 3, 104, 138},
			[]byte{94, 235, 72, 0},
			false,
		},
	}

	for _, tc := range testCases {
		if c.Verify(tc.puzzle, tc.nonce) != tc.result {
			t.Errorf("Verify(%v, %v) failed: expected %v, got %v", tc.puzzle, tc.nonce, tc.result, !tc.result)
		}
	}
}
