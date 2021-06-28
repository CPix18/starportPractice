package keeper

import (
	"github.com/webinar/chat/x/chat/types"
)

var _ types.QueryServer = Keeper{}
