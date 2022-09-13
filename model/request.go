package model

import "sync"

var RowRequestPool = sync.Pool{
	New: func() interface{} {
		return new(RowRequest)
	},
}

// 请求信息
type RowRequest struct {
	RuleKey   string        //规则
	LogName   string        //日志文件的名称，和LogPos组合可唯一
	LogPos    uint32        //事件结束位置，和LogName组合可唯一
	Action    string        //动作
	Timestamp uint32        //提交时间戳，如果在事务内，则为事务开始时间
	Old       []interface{} //修改前的数据
	Row       []interface{} //修改后的数据
}

// 位置信息
type PosRequest struct {
	Name  string
	Pos   uint32
	Force bool
}

func BuildRowRequest() *RowRequest {
	return RowRequestPool.Get().(*RowRequest)
}

func ReleaseRowRequest(t *RowRequest) {
	RowRequestPool.Put(t)
}
