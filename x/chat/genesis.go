package chat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/webinar/chat/x/chat/keeper"
	"github.com/webinar/chat/x/chat/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the message
	for _, elem := range genState.MessageList {
		k.SetMessage(ctx, *elem)
	}

	// Set message count
	k.SetMessageCount(ctx, genState.MessageCount)

	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all message
	messageList := k.GetAllMessage(ctx)
	for _, elem := range messageList {
		elem := elem
		genesis.MessageList = append(genesis.MessageList, &elem)
	}

	// Set the current count
	genesis.MessageCount = k.GetMessageCount(ctx)

	genesis.PortId = k.GetPort(ctx)

	return genesis
}
