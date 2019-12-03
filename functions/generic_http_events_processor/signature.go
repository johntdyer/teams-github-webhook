package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

// CheckMAC reports whether messageSig is a valid HMAC tag for message.
func checkMAC(unsignedData, receivedHMAC, key string) bool {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(unsignedData))
	expectedMAC := mac.Sum(nil)
	// fmt.Println(hex.EncodeToString([]byte(expectedMAC)))
	// log.Debugf("## checkMAC sig: secret: %s", key)
	log.Debugf("## checkMAC sig: messageSig: %x", string(receivedHMAC))
	log.Debugf("## checkMAC sig: computedSig: %x", hex.EncodeToString([]byte(expectedMAC)))
	return receivedHMAC == hex.EncodeToString([]byte(expectedMAC))
}
