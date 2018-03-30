package secp256k1

import (
	"errors"

	"github.com/medibloc/go-medibloc/crypto/signature"
	"github.com/medibloc/go-medibloc/crypto/signature/algorithm"
)

type Signature struct {
	privateKey *PrivateKey
	publicKey  *PublicKey
}

func (s *Signature) Algorithm() algorithm.Algorithm {
	return algorithm.SECP256K1
}

func (s *Signature) InitSign(priv signature.PrivateKey) {
	s.privateKey = priv.(*PrivateKey)
}

func (s *Signature) Sign(data []byte) ([]byte, error) {
	if s.privateKey == nil {
		return nil, errors.New("signature private key is nil")
	}

	sig, err := s.privateKey.Sign(data)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (s *Signature) RecoverPublic(data []byte, sig []byte) (signature.PublicKey, error) {
	pub, err := RecoverPubkey(data, sig)
	if err != nil {
		return nil, err
	}
	pubKey, err := ToECDSAPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return NewPublicKey(*pubKey), nil
}

func (s *Signature) InitVerify(pub signature.PublicKey) {
	s.publicKey = pub.(*PublicKey)
}

func (s *Signature) Verify(data []byte, sig []byte) (bool, error) {
	if s.publicKey == nil {
		return false, errors.New("signature public key is nil")
	}

	return s.publicKey.Verify(data, sig)
}