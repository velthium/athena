package profiles

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/djuno/database"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ modules.Module = &Module{}
var _ modules.GenesisModule = &Module{}
var _ modules.MessageModule = &Module{}

// Module represents the x/profiles module handler
type Module struct {
	encodingConfig *params.EncodingConfig
	db             *database.DesmosDb
}

// NewModule allows to build a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, db *database.DesmosDb) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(appState, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.db)
}