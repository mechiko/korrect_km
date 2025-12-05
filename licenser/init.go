package licenser

import (
	crypto_rand "crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/nacl/box"
)

func init() {
	boxLocal = &BoxKeys{
		ServerPub: serverCertPub,
	}
	if keyExists(root, nameModuleKey) {
		if valueExists(root, nameModuleKey, nameLicenseKey) {
			licenseKeyValue, _ = readStringValueWithDefault(root, nameModuleKey, nameLicenseKey, "")
		}
		if valueExists(root, nameModuleKey, nameClientCertPubKey) {
			boxLocal.LocalPub, _ = readStringValueWithDefault(root, nameModuleKey, nameClientCertPubKey, "")
		}
		if valueExists(root, nameModuleKey, nameClientCertPrivKey) {
			boxLocal.LocalPriv, _ = readStringValueWithDefault(root, nameModuleKey, nameClientCertPrivKey, "")
		}
	} else {
		key, err := createKey(root, nameModuleKey)
		if err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameLicenseKey, ""); err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameClientCertPrivKey, ""); err != nil {
			panic(err)
		}
		if err := key.SetStringValue(nameClientCertPubKey, ""); err != nil {
			panic(err)
		}
	}
	if boxLocal.LocalPriv == "" || boxLocal.LocalPub == "" {
		// создаем ключи и прописываем
		publicKey, privateKey, err := box.GenerateKey(crypto_rand.Reader)
		if err != nil {
			panic(err)
		}
		boxLocal.LocalPub = hex.EncodeToString(publicKey[:])
		boxLocal.LocalPriv = hex.EncodeToString(privateKey[:])
		if err := writeStringValue(root, nameModuleKey, nameClientCertPrivKey, boxLocal.LocalPriv); err != nil {
			panic(err)
		}
		if err := writeStringValue(root, nameModuleKey, nameClientCertPubKey, boxLocal.LocalPub); err != nil {
			panic(err)
		}
	}
}
