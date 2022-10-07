package component

import (
	"context"
	"errors"
	"fmt"
	"io"

	ipfsFiles "github.com/ipfs/go-ipfs-files"
	ipfsPath "github.com/ipfs/interface-go-ipfs-core/path"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/metadata"
)

// Handler for the "get" operation, which retrieves a document
func (b *IPFSBinding) getOperation(ctx context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	reqMetadata := &getRequestMetadata{}
	err := reqMetadata.FromMap(req.Metadata)
	if err != nil {
		return nil, err
	}

	if reqMetadata.Path == "" {
		return nil, errors.New("metadata property 'path' is empty")
	}
	p := ipfsPath.New(reqMetadata.Path)
	err = p.IsValid()
	if err != nil {
		return nil, fmt.Errorf("invalid value for metadata property 'path': %v", err)
	}

	res, err := b.ipfsAPI.Unixfs().Get(ctx, p)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	f, ok := res.(ipfsFiles.File)
	if !ok {
		return nil, errors.New("path does not represent a file")
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &bindings.InvokeResponse{
		Data:     data,
		Metadata: nil,
	}, nil
}

type getRequestMetadata struct {
	Path string `mapstructure:"path"`
}

func (m *getRequestMetadata) FromMap(mp map[string]string) (err error) {
	if len(mp) > 0 {
		err = metadata.DecodeMetadata(mp, m)
		if err != nil {
			return err
		}
	}
	return nil
}
