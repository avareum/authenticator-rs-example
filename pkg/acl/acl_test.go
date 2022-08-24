package acl

import (
	"crypto/ed25519"
	"testing"

	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
	"github.com/avareum/avareum-hubble-signer/tests/fixtures"
	"github.com/gagliardetto/solana-go"
	"github.com/test-go/testify/require"
)

func Test_SignatureVerification(t *testing.T) {

	t.Run("should verify payload & signature", func(t *testing.T) {
		aclSuite := fixtures.NewTestACL()
		acl, err := NewServiceACLWithOpt(ServiceACLOptions{SkipFetchOnVerify: true})
		require.Nil(t, err)

		acc1, err := solana.NewRandomPrivateKey()
		require.Nil(t, err)
		acc2, err := solana.NewRandomPrivateKey()
		require.Nil(t, err)

		type VerifyTestCase struct {
			signer    solana.PublicKey
			signature string
			payload   []byte
			expect    bool
		}

		tests := []VerifyTestCase{
			{
				signer:    acc1.PublicKey(),
				payload:   []byte("hello"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(acc1), []byte("hello")),
				expect:    true,
			},
			{
				signer:    acc2.PublicKey(),
				payload:   []byte("world"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(acc2), []byte("world")),
				expect:    true,
			},
			{
				signer:    acc1.PublicKey(),
				payload:   []byte("payload"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(acc1), []byte("mismatch payload")),
				expect:    false,
			},
			{
				signer:    acc2.PublicKey(),
				payload:   []byte("acc2, sig by acc1"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(acc1), []byte("acc2, sig by acc1")),
				expect:    false,
			},
		}

		for _, test := range tests {
			signature := solana.MustSignatureFromBase58(test.signature)
			require.Equal(t, test.expect, acl.Verify(test.signer[:], test.payload, signature[:]))
		}
	})

	t.Run("should verify can call", func(t *testing.T) {
		aclSuite := fixtures.NewTestACL()
		sm := fixtures.NewTestSecretManager()
		acl, err := NewServiceACLWithOpt(ServiceACLOptions{SkipFetchOnVerify: true, Prefix: "SERVICE_", SecretManager: sm})
		require.Nil(t, err)

		type CanCallTestCase struct {
			service   string
			signature string
			payload   []byte
			expect    bool
		}

		service1, err := solana.NewRandomPrivateKey()
		require.Nil(t, err)
		unauthorizedService1, err := solana.NewRandomPrivateKey()
		require.Nil(t, err)

		// [hack] create service key
		sm.Create(types.NewSecretServiceID("service1"), service1)

		tests := []CanCallTestCase{
			{
				service:   "service1",
				payload:   []byte("hello"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(service1), []byte("hello")),
				expect:    true,
			},
			{
				service:   "service1",
				payload:   []byte("payload"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(service1), []byte("mismatch payload")),
				expect:    false,
			},
			{
				service:   "introduce as service1, signed by unauthorizedService1",
				payload:   []byte("payload"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(unauthorizedService1), []byte("introduce as service1, signed by unauthorizedService1")),
				expect:    false,
			},
			{
				service:   "unauthorizedService1",
				payload:   []byte("unauthorized service"),
				signature: aclSuite.MustSignPayloadWithKey(ed25519.PrivateKey(unauthorizedService1), []byte("unauthorized service")),
				expect:    false,
			},
		}

		for _, test := range tests {
			signature := solana.MustSignatureFromBase58(test.signature)
			require.Equal(t, test.expect, acl.CanCall(test.service, test.payload, signature[:]))
		}
	})

}

func Test_SignatureVerification2(t *testing.T) {

	t.Run("should verify signature", func(t *testing.T) {
		signerPub := solana.MustPublicKeyFromBase58("48oeRt6qvLpDrxwnjcooCSCYXvY1Wj5eYztBeqw938jg")
		payload := []byte(`{"method":"eth_call","params":[{"from":"0x407d73d8a49eeb85d32cf465507dd71d507100c1","to":"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b","value":"0x186a0"}],"id":1,"jsonrpc":"2.0"}`)
		payloadSignature := solana.MustSignatureFromBase58("5SeDH1KWhMcvART54uzCX38DaxVMjkWtwNfaxYL1F1XELKxi4yVxrjKLpskJfBhWsqPsMbQqrTFCtW3AnF9JxP6f")
		require.Equal(t, true, ed25519.Verify(signerPub[:], payload, payloadSignature[:]))
	})

}
