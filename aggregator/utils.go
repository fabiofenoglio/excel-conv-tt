package aggregator

import (
	"crypto/sha1"
	"encoding/base64"
)

func Base64Sha(content []byte) string {
	h := sha1.New()
	h.Write(content)
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}
