package kafka_toolkit

import "context"

const (
	// ContextOffsetKey key de offset en context
	ContextOffsetKey = "OFFSET"
	// ContextPartitionKey key de offset en context
	ContextPartitionKey = "PARTITION"
)

// ContextWithOffset agrega el offset al context
func ContextWithOffset(ctx context.Context, offset int64) context.Context {
	return context.WithValue(ctx, ContextOffsetKey, offset)
}

// ContextWithPartition agrega partion al context
func ContextWithPartition(ctx context.Context, partition int32) context.Context {
	return context.WithValue(ctx, ContextPartitionKey, partition)
}

// OffsetFromContext obtiene offset desde el context
func OffsetFromContext(ctx context.Context) int64 {
	return ctx.Value(ContextOffsetKey).(int64)
}

// PartitionFromContext obtiene partition desde el context
func PartitionFromContext(ctx context.Context) int32 {
	return ctx.Value(ContextPartitionKey).(int32)
}
