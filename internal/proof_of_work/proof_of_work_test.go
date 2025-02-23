package proof_of_work

import (
	"bytes"
	"testing"
)

func TestProofOfWork(t *testing.T) {
	pow := NewProofOfWork()

	length := 32
	difficulty := 4

	generatedChallenge, err := pow.GenerateChallenge(length)

	if err != nil {
		t.Fatalf("expected err to be nil, but got %v", err)
	}

	if len(generatedChallenge) != length {
		t.Fatalf("expected length of challenge to be %d, but got %d", length, len(generatedChallenge))
	}

	challenge := []byte("Hello World")

	nonce := pow.Solve(challenge, difficulty)

	if len(nonce) != 6 {
		t.Fatalf("expected length of nonce to be %d, but got %d", 5, len(nonce))
	}

	expectedNone := []byte{49, 48, 55, 49, 48, 53}

	if !bytes.Equal(nonce, expectedNone) {
		t.Fatalf("expected nonce to be %v, but got %v", expectedNone, nonce)
	}

	if !pow.Verify(challenge, nonce, difficulty) {
		t.Fatalf("expected challenge to be verified")
	}
}
