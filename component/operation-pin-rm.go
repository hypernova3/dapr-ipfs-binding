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

// Handler for the "pin-rm" operation, which removes a pin
func (b *IPFSBinding) pinRmOperation(ctx context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	reqMetadata := &pinRmRequestMetadata{}
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

	opts, err := reqMetadata.PinRmOptions()
	if err != nil {
		return nil, err
	}
	err = b.ipfsAPI.Pin().Rm(ctx, p, opts...)
	if err != nil {
		return nil, err
	}

	return &bindings.InvokeResponse{
		Data:     nil,
		Metadata: nil,
	}, nil
}

type pinRmRequestMetadata struct {
	Path      string `mapstructure:"path"`
	Recursive *bool  `mapstructure:"recursive"`
}

func (m *pinRmRequestMetadata) FromMap(mp map[string]string) (err error) {
	if len(mp) > 0 {
		err = metadata.DecodeMetadata(mp, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *pinRmRequestMetadata) PinRmOptions() ([]ipfsOptions.PinRmOption, error) {
	opts := []ipfsOptions.PinRmOption{}
	if m.Recursive != nil {
		opts = append(opts, ipfsOptions.Pin.RmRecursive(*m.Recursive))
	} else {
		opts = append(opts, ipfsOptions.Pin.RmRecursive(true))
	}
	return opts, nil
}
