package acl

import (
	"context"
	"crypto/ed25519"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/avareum/avareum-hubble-signer/pkg/acl/types"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/api/iterator"
)

type ServiceACLOptions struct {
	ProjectID         string
	Bucket            string
	SkipFetchOnVerify bool
}

type ServiceACL struct {
	types.ACL
	serviceKeys   map[string][]byte
	opt           ServiceACLOptions
	storageClient *storage.Client
}

func NewServiceACL() (*ServiceACL, error) {
	return NewServiceACLWithOpt(ServiceACLOptions{
		ProjectID:         os.Getenv("GCP_PROJECT"),
		Bucket:            "service-keys",
		SkipFetchOnVerify: false,
	})
}

func NewServiceACLWithOpt(opt ServiceACLOptions) (*ServiceACL, error) {
	w := &ServiceACL{
		opt:         opt,
		serviceKeys: map[string][]byte{},
	}
	err := w.init()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *ServiceACL) init() error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	w.storageClient = client
	return nil
}

func (w *ServiceACL) fetchServiceKeys() error {
	bkt := w.storageClient.Bucket(w.opt.Bucket)

	// list all blobs in the bucket
	query := &storage.Query{Prefix: ""}
	iter := bkt.Objects(context.TODO(), query)

	serviceKeys := map[string][]byte{}
	for {
		attrs, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// ready object body
		rc, err := bkt.Object(attrs.Name).NewReader(context.TODO())
		if err != nil {
			return err
		}
		defer rc.Close()
		pub, err := io.ReadAll(rc)
		if err != nil {
			return err
		}

		serviceKeys[attrs.Name] = pub
	}
	w.serviceKeys = serviceKeys
	return nil
}

// setServiceKey sets the service key for the given service name.
func (w *ServiceACL) setServiceKey(serviceName string, pub []byte) {
	w.serviceKeys[serviceName] = pub
}

// getServiceKey returns the public key for the given service name.
func (w *ServiceACL) getServiceKey(serviceName string) []byte {
	return w.serviceKeys[serviceName]
}

/*
 ACL implementaiton for GCP Secret Manager
*/

func (w *ServiceACL) Verify(pub ed25519.PublicKey, payload []byte, payloadSignature []byte) bool {
	return ed25519.Verify(pub, payload, payloadSignature)
}

func (w *ServiceACL) CanCall(serviceName string, payload []byte, payloadSignature []byte) bool {
	// fetch if the service keys are available
	if !w.opt.SkipFetchOnVerify {
		err := w.fetchServiceKeys()
		if err != nil {
			log.Println(err)
			return false
		}
	}

	pubBytes := w.getServiceKey(serviceName)
	if len(pubBytes) == 0 {
		return false
	}
	pub := solana.PublicKeyFromBytes(pubBytes)
	return w.Verify(pub[:], payload, payloadSignature)
}
