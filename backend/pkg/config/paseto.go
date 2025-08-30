package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	"aidanwoods.dev/go-paseto"
	"github.com/sirupsen/logrus"
)

var PasetoKey paseto.V4SymmetricKey

func LoadPasetoSecret() {
	secret := os.Getenv("PASETO_SECRET")

	if secret == "" {
		logrus.Warn("Generating Random Paseto Secret...")
		PasetoKey = paseto.NewV4SymmetricKey()
	} else {
		hash := sha256.Sum256([]byte(secret))

		key, err := paseto.V4SymmetricKeyFromBytes(hash[:])
		if err != nil {
			panic("Failed to create Paseto key: " + err.Error())
		}

		PasetoKey = key
	}

	logrus.Debugf("Generated Paseto Key: %s", hex.EncodeToString(PasetoKey.ExportBytes()))
}
