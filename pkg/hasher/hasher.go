package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	Cost int
}

func New(cost int) Hasher {
	return Hasher{Cost: cost}
}

func (h Hasher) Hash(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(str), h.Cost)
	return string(b), err
}

func (h Hasher) Compare(str, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
