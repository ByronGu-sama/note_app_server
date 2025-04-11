package kafkaAction

const (
	LikeNote       = iota // 点赞笔记
	DislikeNote           // 取消点赞笔记
	CollectNote           // 收藏笔记
	AbandonNote           // 取消收藏笔记
	DelNoteComment        // 删除笔记评论
	LikeComment           // 点赞评论
	DislikeComment        // 取消点赞评论
	SyncNote              // 同步帖子
	DelNote               // 删除笔记
	SyncMessage           // 同步聊天记录
)
