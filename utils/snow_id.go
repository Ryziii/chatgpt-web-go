package utils

import "github.com/bwmarrin/snowflake"

var (
	snowflakeNode *snowflake.Node
)

func init() {
	snowflakeNode, _ = snowflake.NewNode(1)
}
func GetSnowIdInt64() int64 {
	return snowflakeNode.Generate().Int64()
}

func GetSnowIdUint64() uint64 {
	return uint64(snowflakeNode.Generate().Int64())
}
