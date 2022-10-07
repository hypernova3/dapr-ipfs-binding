package component

import (
	"context"
	"errors"
	"fmt"

	ipfsOptions "github.com/ipfs/interface-go-ipfs-core/options"
	ipfsPath "github.com/ipfs/interface-go-ipfs-core/path"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/metadata"
)

// Handler for the "pin-add" operation, which adds a new pin
func (b *IPFSBinding) pinAddOperation(ctx context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	reqMetadata := &pinAddRequestMetadata{}
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

	opts, err := reqMetadata.PinAddOptions()
	if err != nil {
		return nil, err
	}
	err = b.ipfsAPI.Pin().Add(ctx, p, opts...)
	if err != nil {
		return nil, err
	}

	return &bindings.InvokeResponse{
		Data:     nil,
		Metadata: nil,
	}, nil
}

type pinAddRequestMetadata struct {
	Path      string `mapstructure:"path"`
	Recursive *bool  `mapstructure:"recursive"`
}

func (m *pinAddRequestMetadata) FromMap(mp map[string]string) (err error) {
	if len(mp) > 0 {
		err = metadata.DecodeMetadata(mp, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *pinAddRequestMetadata) PinAddOptions() ([]ipfsOptions.PinAddOption, error) {
	opts := []ipfsOptions.PinAddOption{}
	if m.Recursive != nil {
		opts = append(opts, ipfsOptions.Pin.Recursive(*m.Recursive))
	} else {
		opts = append(opts, ipfsOptions.Pin.Recursive(true))
	}
	return opts, nil
}
