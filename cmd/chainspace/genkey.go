package main

import (
	"encoding/base32"
	"fmt"
	"os"

	"chainspace.io/prototype/config"
	"chainspace.io/prototype/crypto/signature"
	"chainspace.io/prototype/crypto/transport"
	"github.com/tav/golly/log"
	"gopkg.in/yaml.v2"
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

func genKeys(path string) (signature.KeyPair, *transport.Cert, error) {
	signingKey, err := signature.GenKeyPair(signature.Ed25519)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate signing key: %s", err)
	}
	cert, err := transport.GenCert(transport.ECDSA)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate transport cert: %s", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	cfg := config.Keys{
		SigningKey: &config.SigningKey{
			Private: b32.EncodeToString(signingKey.PrivateKey().Value()),
			Public:  b32.EncodeToString(signingKey.PublicKey().Value()),
			Type:    signingKey.Algorithm().String(),
		},
		TransportCert: &config.TransportCert{
			Private: cert.Private,
			Public:  cert.Public,
			Type:    cert.Type.String(),
		},
	}
	enc := yaml.NewEncoder(f)
	err = enc.Encode(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("could not write data to %s: %s", path, err)
	}
	return signingKey, cert, nil
}

func cmdGenKey(args []string, usage string) {
	opts := newOpts("genkey OPTIONS", usage)
	path := opts.Flags("-o", "--output").Label("PATH").String("path to write the generated keys [keys.yaml]")
	opts.Parse(args)
	if _, _, err := genKeys(*path); err != nil {
		log.Fatalf("%s", err)
	}
	log.Infof("Generated keys successfully written to %s", *path)
}