package silence

import (
	"errors"

	"github.com/evcc-io/evcc/util/log"
	"github.com/evcc-io/evcc/util/request"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
)

// API is an api.Vehicle implementation for Silence scooters
type API struct {
	*request.Helper
}

func NewAPI(log log.Logger, identity *Identity) *API {
	v := &API{
		Helper: request.NewHelper(log),
	}

	v.Client.Transport = &oauth2.Transport{
		Source: identity,
		Base:   v.Client.Transport,
	}

	return v
}

// vehicles provides list of vehicles response
func (v *API) vehicles() ([]Vehicle, error) {
	var resp []Vehicle
	err := v.GetJSON(ApiUri, &resp)
	return resp, err
}

// Vehicles provides list of Vehicle IDs
func (v *API) Vehicles() ([]string, error) {
	var resp []Vehicle

	err := v.GetJSON(ApiUri, &resp)
	if err == nil {
		return lo.Map(resp, func(v Vehicle, _ int) string {
			return v.FrameNo
		}), nil
	}

	return nil, err
}

// Status provides the vehicle status response
func (v *API) Status(vin string) (Vehicle, error) {
	resp, err := v.vehicles()

	if err == nil {
		for _, vv := range resp {
			if vv.FrameNo == vin {
				return vv, nil
			}
		}

		err = errors.New("vehicle not found")
	}

	return Vehicle{}, err
}
