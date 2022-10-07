package component

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ipfsOptions "github.com/ipfs/interface-go-ipfs-core/options"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/metadata"
)

// Handler for the "pin-ls" operation, which removes a pin
func (b *IPFSBinding) pinLsOperation(ctx context.Context, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	reqMetadata := &pinLsRequestMetadata{}
	err := reqMetadata.FromMap(req.Metadata)
	if err != nil {
		return nil, err
	}

	opts, err := reqMetadata.PinLsOptions()
	if err != nil {
		return nil, err
	}
	ls, err := b.ipfsAPI.Pin().Ls(ctx, opts...)
	if err != nil {
		return nil, err
	}

	res := pinLsOperationResponse{}
	for e := range ls {
		err = e.Err()
		if err != nil {
			return nil, err
		}
		res = append(res, pinLsOperationResponseItem{
			Type: e.Type(),
			Cid:  e.Path().Cid().String(),
		})
	}

	j, _ := json.Marshal(res)
	return &bindings.InvokeResponse{
		Data:     j,
		Metadata: nil,
	}, nil
}

type pinLsOperationResponseItem struct {
	Cid  string `json:"cid,omitempty"`
	Type string `json:"type,omitempty"`
}

type pinLsOperationResponse []pinLsOperationResponseItem

type pinLsRequestMetadata struct {
	Type *string `mapstructure:"type"`
}

func (m *pinLsRequestMetadata) FromMap(mp map[string]string) (err error) {
	if len(mp) > 0 {
		err = metadata.DecodeMetadata(mp, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *pinLsRequestMetadata) PinLsOptions() ([]ipfsOptions.PinLsOption, error) {
	opts := []ipfsOptions.PinLsOption{}
	if m.Type != nil {
		switch strings.ToLower(*m.Type) {
		case "direct":
			opts = append(opts, ipfsOptions.Pin.Ls.Direct())
		case "recursive":
			opts = append(opts, ipfsOptions.Pin.Ls.Recursive())
		case "indirect":
			opts = append(opts, ipfsOptions.Pin.Ls.Indirect())
		case "all":
			opts = append(opts, ipfsOptions.Pin.Ls.All())
		default:
			return nil, fmt.Errorf("invalid value for metadata property 'type'")
		}
	} else {
		opts = append(opts, ipfsOptions.Pin.Ls.All())
	}
	return opts, nil
}
