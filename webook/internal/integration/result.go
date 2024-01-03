package integration

// 输入
type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 输出
type Result[T any] struct {
	Data T      `json:"data"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
