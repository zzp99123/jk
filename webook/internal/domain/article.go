// 帖子
package domain

import "time"

// 帖子标题和内容
type Article struct {
	Id      int64         //id
	Title   string        //标题
	Content string        //内容
	Author  Author        //作者id
	Status  ArticleStatus //状态
	Ctime   time.Time     //创建时间
	Utime   time.Time     //更新时间
}

// 帖子用户
type Author struct {
	Id   int64
	Name string
}
type ArticleStatus uint8

// 定义三种状态
// 新建-保存，那么直接就是未发表。
// 修改-保存，那么依旧保持未发表。
// 发表-修改-保存，那么应该是从发表状态变回未发表状态
const (
	ArticleStatusUnknown     ArticleStatus = iota //文章状态未知
	ArticleStatusUnpublished                      //文章状态未发布
	ArticleStatusPublished                        //文章状态已发布
	ArticleStatusPrivate                          //文章状态不可见
)

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusUnpublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	case ArticleStatusPrivate:
		return "private"
	default:
		return "unknown"
	}
}

//go:inline
func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

// 中文
func (a Article) Abstract() string {
	// 摘要我们取前几句。
	// 要考虑一个中文问题
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	// 英文怎么截取一个完整的单词，我的看法是……不需要纠结，就截断拉到
	// 词组、介词，往后找标点符号
	return string(cs[:100])
}

// ArticleStatusV1 如果你的状态很复杂，有很多行为（就是你要搞很多方法），状态里面需要一些额外字段
// 就用这个版本
//
//	type ArticleStatusV1 struct {
//		Val  uint8
//		Name string
//	}
//
// var (
//
//	ArticleStatusV1Unknown = ArticleStatusV1{Val: 0, Name: "unknown"}
//
// )
