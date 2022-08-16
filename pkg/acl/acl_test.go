package acl

import (
	"encoding/base64"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

type SignerTestCase struct {
	signer    string
	signature string
	payload   string
	expect    bool
}

func Test_SignatureVerification(t *testing.T) {
	assert := assert.New(t)

	t.Run("should verify signatures", func(t *testing.T) {
		w, err := NewServiceACL()
		assert.Nil(err)
		tests := []SignerTestCase{
			{
				signer:    "2V7t5NaKY7aGkwytCWQgvUYZfEr9XMwNChhJEakTExk6",
				signature: "3APVUYY2Qcq9WwPSciQ1xicshQcsXJWhgD1x5YEMNnNnKtpaAJtRhDVMWcvePosamk3JfSKpBM8qt5ZkAQgRNSzo",
				payload:   "AQACBBYPusE6993YBdMXCj3gxr2XEmoeAsDSWdCobvgh1uXHoaZGX0wuvyRMMdgLyVwnNFo0JOQowt7zPs7Z6Q0/cBsGp9UXGMd0yShWY5hpHV62i164o5tLbVxzVVshAAAAANzl6+HknDufEUy1VExQqZ7A1pLWP1Z5WuAprIPZ6ovia//gAIGoVfMJgSuJp8/VsCUq8qvMdiHNGrMrQAxoS2QBAwMAAQIoAgAAAAcAAAABAAAAAAAAAM/ncrkCAAAACJ9fAQAAAAD7FSgHAAAAAA==",
				expect:    true,
			},
			{
				signer:    "AXUChvpRwUUPMJhA4d23WcoyAL7W8zgAeo7KoH57c75F",
				signature: "2c5u22N6Yyjj7qQdGtEshp4n9r4akLyRLJijnXL9HnQAXqq7S9cfnGdi5mwnfQ5kJHQuRv62T366SKaztHiUUqnK",
				payload:   "AQACBo2HXSJ3fw0pdiiQj9fBoRA5Wicjyla4ARhrK90xeo0mD4K9wHSEzSq9lEFUOSPloeUrUL/2uV3S6+lTeGNtTMtTP2NRTSucECyiWtS2HzfxwLnbnxVbXNH0T0egDFCfVK5ESO7I2Pz56XFyxewO0gbya1rvqPPu0CGg/LC/Wl5BhQ8tbgKkevgk0Jq2ncQtcMsoy/okn7fuV7nSVsEnYu8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIvotbwkk/hzTtAxHIaRCSJjB4WdQUXrn+6aEapl6XgPAgQFAQIDAgIHAAMAAAABAAUCAAAMAgAAAOsNAAAAAAAA",
				expect:    true,
			},
			{
				signer:    "2vKtu3nW1TS6iPvJPK8R88B5QfDrwJDwwB11Uu1CN9o7",
				signature: "3APVUYY2Qcq9WwPSciQ1xicshQcsXJWhgD1x5YEMNnNnKtpaAJtRhDVMWcvePosamk3JfSKpBM8qt5ZkAQgRNSzo",
				payload:   "AQACBBYPusE6993YBdMXCj3gxr2XEmoeAsDSWdCobvgh1uXHoaZGX0wuvyRMMdgLyVwnNFo0JOQowt7zPs7Z6Q0/cBsGp9UXGMd0yShWY5hpHV62i164o5tLbVxzVVshAAAAANzl6+HknDufEUy1VExQqZ7A1pLWP1Z5WuAprIPZ6ovia//gAIGoVfMJgSuJp8/VsCUq8qvMdiHNGrMrQAxoS2QBAwMAAQIoAgAAAAcAAAABAAAAAAAAAM/ncrkCAAAACJ9fAQAAAAD7FSgHAAAAAA==",
				expect:    false,
			},
		}

		for _, test := range tests {
			pub := solana.MustPublicKeyFromBase58(test.signer)
			sig := solana.MustSignatureFromBase58(test.signature)
			msg, err := base64.StdEncoding.DecodeString(test.payload)
			assert.Nil(err)
			assert.Equal(w.Verify(pub[:], msg, sig[:]), test.expect)
		}
	})

}
