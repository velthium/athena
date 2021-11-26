package profiles

import (
	"context"
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v2/node/remote"

	"github.com/desmos-labs/djuno/types"
)

// UpdateProfiles updates the profiles associated with the given addresses, if any.
func (m *Module) UpdateProfiles(height int64, addresses []string) error {
	for _, address := range addresses {
		res, err := m.profilesClient.Profile(
			context.Background(),
			profilestypes.NewQueryProfileRequest(address),
			remote.GetHeightRequestHeader(height),
		)
		if err != nil {
			return fmt.Errorf("error while getting profile from gRPC: %s", err)
		}

		if res.Profile != nil {
			var account authtypes.AccountI
			err = m.cdc.UnpackAny(res.Profile, &account)
			if err != nil {
				return fmt.Errorf("error while unpacking profile: %s", err)
			}

			err = m.db.SaveProfile(types.NewProfile(account.(*profilestypes.Profile), height))
			if err != nil {
				return fmt.Errorf("error while saving profile: %s", err)
			}
		}
	}

	return nil
}

// updateParams allows to update the profiles params by fetching them from the chain
func (m *Module) updateParams() error {
	height, err := m.db.LastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	res, err := m.profilesClient.Params(
		context.Background(),
		&profilestypes.QueryParamsRequest{},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	err = m.db.SaveProfilesParams(types.NewProfilesParams(res.Params, height))
	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}