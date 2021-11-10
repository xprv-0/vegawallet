package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestDescribeKey(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	walletName := vgrand.RandomStr(5)

	cmd := []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	}

	// when
	generateKeyResp, err := KeyGenerate(t, append(cmd,
		"--meta", "name:key-1,role:validation",
	))

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithWalletCreation().
		WithName(walletName).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	descResp, err := KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithPubKey(generateKeyResp.Key.KeyPair.PublicKey).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		WithAlgorithm("vega/ed25519", 1).
		WithTainted(false)

	// when non-existent public key
	descResp, err = KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey[1:],
	))
	require.Error(t, err)
	require.Nil(t, descResp)
}