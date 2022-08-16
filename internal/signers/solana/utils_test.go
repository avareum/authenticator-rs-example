package solana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Base64ToTransaction(t *testing.T) {
	assert := assert.New(t)
	/*
		Program id: uvrXzrjzPv7cAsjgnZy3hmMUnff2KiKaJpj8Ab3ffny
		Data:{
		  args: [
		    {
		      name: 'address',
		      type: 'publicKey',
		      data: '8pqnoox4Yi5G97uRb1ZaFh2LnroRG92Uc41tjZweyxUR'
		    },
		    { name: 'action', type: 'Action', data: 'Execute' }
		  ],
		  accounts: [
		    {
		      name: 'Acl',
		      pubkey: [PublicKey],
		      isSigner: false,
		      isWritable: false
		    },
		    {
		      name: 'Acl Certificate',
		      pubkey: [PublicKey],
		      isSigner: false,
		      isWritable: false
		    }
		  ]
		}
	*/
	bytecode := "AbkY5ti8jYg+FEiLVZkNkcuiQg3lVDUhsmqt2gc6WkiTufP1w5XVKLeKIx6tskqB6OG9IyyZMa8kTqqNMlSUNAcBAAMEt8bYqLqc/twuJxYcoviVHbmpnP7SHwY8g7yrCSIvt7A0VuBHGbrhF0xVbvK2PR3B3uyfbOREdYeTuaHgVSDaClv49i3Ee53OplR/+V9cERMCoMHOcaHv1p7yG0/ACnv6DY8YU/p0tPYxNs5WLcYbs7o0YSO7AOib7lNyQz/LNtKbWKt9stl1Udab8KuHBQE1ZR/MmbyJqItpEVF3tAk1ZwEDAgIBKXD+xDS/h+igdENF9i0hZwYrIrvJSZ1WL6zXuMsWPfHnE+JLbr2v7MIA"

	t.Run("should extract program signature", func(t *testing.T) {
		tx, err := Base64ToTransaction(bytecode)
		assert.Nil(err)
		assert.NotNil(tx)
	})

	t.Run("should get program id #0", func(t *testing.T) {
		tx, _ := Base64ToTransaction(bytecode)
		programID, err := ProgramID(*tx, 0)
		assert.Nil(err)
		assert.Equal("uvrXzrjzPv7cAsjgnZy3hmMUnff2KiKaJpj8Ab3ffny", programID.String())
	})

	t.Run("should reject invalid program idx", func(t *testing.T) {
		tx, _ := Base64ToTransaction(bytecode)
		_, err := ProgramID(*tx, 1)
		assert.EqualError(err, "invalid program index")
	})
}
