package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	// 현재는 v4까지 나왔으니 v2는 참고만
	token        paseto.Token
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// TODO
	return nil, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// TODO
	return "", nil	
}
