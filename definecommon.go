package gameprovider

// SES_CONTEXT 网络层保存 数据
var SES_CONTEXT = "SES_CONTEXT"

//桌子状态
const (
	//未开始
	TABLESTATUS_WAIT = iota
	//已开始
	TABLESTATUS_PLAYING
	//结束
	TABLESTATUS_ENDING
)
