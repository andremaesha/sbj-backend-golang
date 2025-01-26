package tests

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

type GenerateSha1AndMd5StringImpl struct {
}

func NewGenerateSha1AndMd5StringImpl() *GenerateSha1AndMd5StringImpl {
	return &GenerateSha1AndMd5StringImpl{}
}

func (g *GenerateSha1AndMd5StringImpl) GenerateMD5String(data string) string {
	sum := md5.Sum([]byte(data))
	encodeToString := hex.EncodeToString(sum[:])

	return encodeToString
}

func (g *GenerateSha1AndMd5StringImpl) GenerateSHA1String(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	encodeToString := hex.EncodeToString(hash.Sum(nil))
	return encodeToString
}
func (g *GenerateSha1AndMd5StringImpl) GenerateSHA256String(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	encodeToString := hex.EncodeToString(hash.Sum(nil))
	return encodeToString
}
func (g *GenerateSha1AndMd5StringImpl) GenerateSHA256StringByte(data string) []byte {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

func (g *GenerateSha1AndMd5StringImpl) GenerateSign(userUid string, passwd string, txId string, more ...string) string {
	moreData := strings.Join(more, "")

	userUidPasswd := fmt.Sprintf("%s%s%s", userUid, passwd, moreData)
	if txId != "" {
		userUidPasswd = fmt.Sprintf("%s%s%s%s", userUid, passwd, txId, moreData)
	}
	result := g.GenerateSHA1String(g.GenerateMD5String(userUidPasswd))

	return result
}

func EncryptAES256(plaintext, key string) string {

	var enc GenerateSha1AndMd5StringImpl
	password := enc.GenerateSHA256String(key)
	password = password[0:32]

	password1 := enc.GenerateSHA256StringByte(key)

	ini := string(password1)

	iv := enc.GenerateMD5String(key + "faspay2018xAuth@#")
	iv = iv[16:]
	nilai, err := GetAESEncrypted(plaintext, ini[0:32], iv)
	if err != nil {
		panic(err)
	}
	return nilai
}

func GetAESEncrypted(plaintext, key, iv string) (string, error) {

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	content := []byte(plaintext)
	content = PKCS5Padding(content, block.BlockSize())
	ciphertext := make([]byte, len(content))
	mode.CryptBlocks(ciphertext, content)
	str := base64.StdEncoding.EncodeToString(ciphertext)
	return str, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func setSignatureOtherService(token, faspaySecret, clientKey, clientSecret, path, timestamp string) string {
	encode := clientKey + ":" + clientSecret + ":" + token
	encode = base64.StdEncoding.EncodeToString([]byte(encode))

	stringToSign := fmt.Sprintf("POST:%s:%s:%s:",
		path,
		timestamp,
		encode,
	)

	return EncryptAES256(stringToSign, faspaySecret)
}

func TestAnjinglah(t *testing.T) {
	method := "GET"
	path := "/account/api/bill-inquiry"
	clientKey := "c74cfe5c-a54c-4080-9eae-3d5d0ea1aa44"
	clientSecret := "837c4fa2-182e-4867-bd51-26bdde7f6ed7"
	faspaySecret := "19c887a3-49af-4087-8e82-58a1875b1452"
	appKey := "8e734017-d7c1-49f2-8202-377faaff6a1c"
	appSecret := "f26a4e84-cbba-4023-b3d9-40892c5638d4"
	timestamp := "2024-12-27 14:50:00"

	token := fmt.Sprintf("%s:%s", clientKey, clientSecret)
	token = base64.StdEncoding.EncodeToString([]byte(token))

	plainText := fmt.Sprintf("%s:%s:%s:%s:", method, path, timestamp, token)

	signature := EncryptAES256(plainText, faspaySecret)

	plaintext := fmt.Sprintf("%s:%s", appKey, appSecret)
	authorization := EncryptAES256(plaintext, faspaySecret)

	println("signature: ", signature)
	println("authorization: ", authorization)

	//tokenOtherService := "o5cerUCPSESzOGD83EmotQI6kbo6rdal+TTPzm/vIs9r8GBwCrMZarsqh24ADU+BbgHND1dQd5c4AolthaXHWAJbEZITBwsNrvFR4xWywcS/r9mjr8dfFtjFJOdoXR45pBFsbiRo16mvfomIpLmNaeu2wCPX9AzBSsbTNC50PIh38aCaYFOdj8x9xXBdCs2CyPdaB25uX13yM64nuG44Xlo0vZajPghsqXpYKnEIx+O77iEF3x+ZVwI31iUGaT+FVkuLDCLaUWZ+5MOYi1UCOF9pdKXRxVTC8baeuXWo9CWp0F/ofHo59jV7LTNfcWW7gC7DKh5Y5fslZVwq735wH7L1Ed/mszVmtadYX4qG4UelY3zkTU09eOsMih0DPlYkMPw+Xn4pa5q29CT9MhiKSLdtPET7+8LgFDRxvfVOV1HD38jkEMz4+vbs6pavC1MgrnGzfCVwQ1KUwSp1AbXqNa7TeQSY9ydqaCjV1MdLKQjBhj54EddJyxugRatCUu5HuztK5Is8x1TG3bYFV4OukDTPs4ECxRUEo6HMLDmfnsd59DWrXSP01M67fV9uz4D6Ju+mVr4zW+Nn+dV7ba5BcNrMm5opfLrTSeVp019nGwxjS6WRFy2BXy6yS87spw0q//RYCl3kale3vxYS94gGBd0ggrY4TkZsXGH48PvwhwQ="
	//
	//serviceOtherService := setSignatureOtherService(tokenOtherService, faspaySecret, clientKey, clientSecret, path, timestamp)
	//println("signature other service: ", serviceOtherService)
}
