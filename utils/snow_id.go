package utils

import "github.com/bwmarrin/snowflake"

func GetSnowIdInt64() int64 {
	node, _ := snowflake.NewNode(1)
	return node.Generate().Int64()
}
