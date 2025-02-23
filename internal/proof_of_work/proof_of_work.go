package proof_of_work

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
)

type ProofOfWork struct {
}

func NewProofOfWork() *ProofOfWork {
	return &ProofOfWork{}
}

func (p *ProofOfWork) GenerateChallenge(length int) ([]byte, error) {
	challenge := make([]byte, length)

	_, err := rand.Read(challenge)

	if err != nil {
		return nil, err
	}

	return challenge, nil
}

func (p *ProofOfWork) Verify(challenge []byte, nonce []byte, difficulty int) bool {
	data := append([]byte(nil), challenge...)
	data = append(data, nonce...)

	hash := sha256.Sum256(data)

	return strings.HasPrefix(fmt.Sprintf("%x", hash), strings.Repeat("0", difficulty))
}

func (p *ProofOfWork) Solve(challenge []byte, difficulty int) []byte {
	nonce := 0

	for {
		nonceBytes := []byte(fmt.Sprintf("%d", nonce))

		data := append([]byte(nil), challenge...)
		data = append(data, nonceBytes...)

		hash := sha256.Sum256(data)

		if strings.HasPrefix(fmt.Sprintf("%x", hash), strings.Repeat("0", difficulty)) {
			return nonceBytes
		}

		nonce++
	}
}
