package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/webinar/chat/x/chat/types"
)

func (k msgServer) CreateMessage(goCtx context.Context, msg *types.MsgCreateMessage) (*types.MsgCreateMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var message = types.Message{
		Creator: msg.Creator,
		Body:    msg.Body,
	}

	id := k.AppendMessage(
		ctx,
		message,
	)

	return &types.MsgCreateMessageResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateMessage(goCtx context.Context, msg *types.MsgUpdateMessage) (*types.MsgUpdateMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var message = types.Message{
		Creator: msg.Creator,
		Id:      msg.Id,
		Body:    msg.Body,
	}

	// Checks that the element exists
	if !k.HasMessage(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetMessageOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetMessage(ctx, message)

	return &types.MsgUpdateMessageResponse{}, nil
}

func (k msgServer) DeleteMessage(goCtx context.Context, msg *types.MsgDeleteMessage) (*types.MsgDeleteMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasMessage(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetMessageOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveMessage(ctx, msg.Id)

	return &types.MsgDeleteMessageResponse{}, nil
}
