package consts

const (
	FileStatusActivate = "activate"
	FileStatusDeleted  = "deleted"
	FileStatusReleased = "released"

	DefaultAvatarName                = "default.png"
	AvatarBucketName                 = "avatar"
	FileBucketName                   = "file"
	PresignedGetAvatarURLExpiresSecs = 86400 // 秒
	PresignedPutAvatarURLExpiresSecs = 86400 // 秒

	NeuronIDLength   = 30
	DocumentIDLength = 30

	MainExchangeName = "main_exchange"
	DeadExchangeName = "dead_exchange"
	DeadQueueName    = "dead_queue"
	MessageTTL       = 300000
)
