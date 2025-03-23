package abstractions

import (
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessageContent interface {
	sqsTypes.Message
}

type MessageId interface {
	string
}

type Number interface {
	int8 | int16 | int32 | int64
}

type Queue[C MessageContent, I MessageId, Num Number] interface {
	GetMessages(Num) (*[]C, error)
	DeleteMessage(*I) error
}
