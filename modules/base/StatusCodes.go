package base

var StatusCodes = map[int]string{
	1:"Success",
	-100:"Invalid resource type",
	-101:"Error deserializing JSON to BaseMessage",
	-102:"Error deserializing JSON to FileMessage",
}