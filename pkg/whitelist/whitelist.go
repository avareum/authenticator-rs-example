package whitelist

import (
	"context"
	"crypto/ed25519"
	"io"

	"cloud.google.com/go/storage"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/api/iterator"
)

type WhitelistOptions struct {
	ProjectID string
	Bucket    string
}

type Whitelist struct {
	opt           WhitelistOptions
	storageClient *storage.Client
	serviceKeys   map[string][]byte
}

func NewWhitelist(opt WhitelistOptions) (*Whitelist, error) {
	w := &Whitelist{
		opt:         opt,
		serviceKeys: map[string][]byte{},
	}
	err := w.init()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Whitelist) init() error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	w.storageClient = client
	return nil
}

func (w *Whitelist) Reload() error {
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

func (w *Whitelist) GetPublicKey(serviceName string) []byte {
	return w.serviceKeys[serviceName]
}

func (w *Whitelist) Verify(pub ed25519.PublicKey, payload []byte, payloadSignature []byte) bool {
	return ed25519.Verify(pub, payload, payloadSignature)
}

func (w *Whitelist) CanCall(serviceName string, payload []byte, payloadSignature []byte) bool {
	pubBytes := w.GetPublicKey(serviceName)
	pub := solana.PublicKeyFromBytes(pubBytes)
	return w.Verify(pub[:], payload, payloadSignature)
}
