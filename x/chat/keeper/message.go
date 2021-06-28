package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/webinar/chat/x/chat/types"
	"strconv"
)

// GetMessageCount get the total number of message
func (k Keeper) GetMessageCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageCountKey))
	byteKey := types.KeyPrefix(types.MessageCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetMessageCount set the total number of message
func (k Keeper) SetMessageCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageCountKey))
	byteKey := types.KeyPrefix(types.MessageCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendMessage appends a message in the store with a new id and update the count
func (k Keeper) AppendMessage(
	ctx sdk.Context,
	message types.Message,
) uint64 {
	// Create the message
	count := k.GetMessageCount(ctx)

	// Set the ID of the appended value
	message.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&message)
	store.Set(GetMessageIDBytes(message.Id), appendedValue)

	// Update message count
	k.SetMessageCount(ctx, count+1)

	return count
}

// SetMessage set a specific message in the store
func (k Keeper) SetMessage(ctx sdk.Context, message types.Message) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	b := k.cdc.MustMarshalBinaryBare(&message)
	store.Set(GetMessageIDBytes(message.Id), b)
}

// GetMessage returns a message from its id
func (k Keeper) GetMessage(ctx sdk.Context, id uint64) types.Message {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	var message types.Message
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetMessageIDBytes(id)), &message)
	return message
}

// HasMessage checks if the message exists in the store
func (k Keeper) HasMessage(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	return store.Has(GetMessageIDBytes(id))
}

// GetMessageOwner returns the creator of the message
func (k Keeper) GetMessageOwner(ctx sdk.Context, id uint64) string {
	return k.GetMessage(ctx, id).Creator
}

// RemoveMessage removes a message from the store
func (k Keeper) RemoveMessage(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	store.Delete(GetMessageIDBytes(id))
}

// GetAllMessage returns all message
func (k Keeper) GetAllMessage(ctx sdk.Context) (list []types.Message) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MessageKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Message
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetMessageIDBytes returns the byte representation of the ID
func GetMessageIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetMessageIDFromBytes returns ID in uint64 format from a byte array
func GetMessageIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
