package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/webinar/chat/x/chat/types"
)

func createNMessage(keeper *Keeper, ctx sdk.Context, n int) []types.Message {
	items := make([]types.Message, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Id = keeper.AppendMessage(ctx, items[i])
	}
	return items
}

func TestMessageGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNMessage(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetMessage(ctx, item.Id))
	}
}

func TestMessageExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNMessage(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasMessage(ctx, item.Id))
	}
}

func TestMessageRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNMessage(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveMessage(ctx, item.Id)
		assert.False(t, keeper.HasMessage(ctx, item.Id))
	}
}

func TestMessageGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNMessage(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllMessage(ctx))
}

func TestMessageCount(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNMessage(keeper, ctx, 10)
	count := uint64(len(items))
	assert.Equal(t, count, keeper.GetMessageCount(ctx))
}
