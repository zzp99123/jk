// tdd写法 先写测试 再写实现  gorm写法
package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/integration/startup"
	"goFoundation/webook/internal/repository/dao/article"
	ijwt "goFoundation/webook/internal/web/jwt"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// 初始化一些内容
func (a *ArticleTestSuite) SetupSuite() {
	a.server = gin.Default()
	//模拟已经登录成功
	a.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", ijwt.UserClaims{Uid: 123})
	})
	a.db = startup.InitTestDB()
	ah := startup.InitArticleHandler(article.NewDaoArticle(a.db))
	ah.RegisterRoutes(a.server)
}

// 每次数据都清空从1 开始
func (a *ArticleTestSuite) TearDownTest() {
	// 清空所有数据，并且自增主键恢复到 1
	a.db.Exec("TRUNCATE TABLE articles")
	a.db.Exec("TRUNCATE TABLE published_articles")

}

// 创建和修改
func (a *ArticleTestSuite) TestEdit() {
	t := a.T()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "保存成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				//验证数据库
				var art article.Article
				err := a.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusUnpublished),
				}, art)
			},
			art:      Article{Title: "标题", Content: "内容"},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
				Msg:  "保存成功",
			},
		},
		{
			name: "修改已有帖子，并保存",
			before: func(t *testing.T) {
				// 提前准备数据
				err := a.db.Create(article.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					// 跟时间有关的测试，不是逼不得已，不要用 time.Now()
					// 因为 time.Now() 每次运行都不同，你很难断言
					Ctime: 123,
					Utime: 234,
					// 假设这是一个已经发表了的，然后你去修改，改成了没发表
					Status: uint8(domain.ArticleStatusPublished),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art article.Article
				err := a.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				// 是为了确保我更新了 Utime
				assert.True(t, art.Utime > 234)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					Ctime:    123,
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusUnpublished),
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
				Msg:  "保存成功",
			},
		},
		{
			name: "修改别人的贴子",
			before: func(t *testing.T) {
				//别人的帖子
				err := a.db.Create(article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    123,
					Utime:    234,
					AuthorId: 456,
					Status:   uint8(domain.ArticleStatusPublished),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := a.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					Ctime:    123,
					Utime:    234,
					AuthorId: 456,
					Status:   uint8(domain.ArticleStatusPublished),
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			v.before(t)
			//构造请求
			reqBody, err := json.Marshal(v.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			//json格式数据
			req.Header.Set("Content-Type", "application/json")
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			resp := httptest.NewRecorder()
			a.server.ServeHTTP(resp, req)
			//检验
			assert.Equal(t, v.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			//创建一个webRes
			var webRes Result[int64]
			//把这个resp.Body解码然后赋值给webRes
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			assert.NoError(t, err)
			//验证结果
			assert.Equal(t, v.wantRes, webRes)
			v.after(t)
		})
	}
}

// 发表
func (a *ArticleTestSuite) TestPublish() {
	t := a.T()
	testCase := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建并发表成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				//制作库
				var art article.Article
				err := a.db.Where("author_id = ?", 123).First(&art).Error
				assert.NoError(t, err)
				//确保已经生成了
				assert.True(t, art.Id > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Id = 0
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
				}, art)
				//线上库
				var publishArt article.PublishedArticle
				err = a.db.Where("author_id = ?", 123).First(&publishArt).Error
				assert.NoError(t, err)
				//确保已经生成了
				assert.True(t, publishArt.Id > 0)
				assert.True(t, publishArt.Ctime > 0)
				assert.True(t, publishArt.Utime > 0)
				publishArt.Id = 0
				publishArt.Ctime = 0
				publishArt.Utime = 0
				assert.Equal(t, article.PublishedArticle{
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
				}, publishArt)
			},
			art: Article{
				Title:   "发表标题",
				Content: "发表内容",
			},
			wantCode: 200,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "更新并新发表", //更新完成第一次发表
			before: func(t *testing.T) {
				err := a.db.Create(article.Article{
					Id:       2,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//因为是先更改的
				var art article.Article
				err := a.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 234)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "更改发表标题",
					Content:  "更改发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
					Ctime:    456,
				}, art)
				//线上库
				var publishArt article.PublishedArticle
				err = a.db.Where("id = ?", 2).First(&publishArt).Error
				assert.NoError(t, err)
				//确保已经生成了

				assert.True(t, publishArt.Ctime > 0)
				assert.True(t, publishArt.Utime > 0)

				publishArt.Ctime = 0
				publishArt.Utime = 0
				assert.Equal(t, article.PublishedArticle{
					Id:       2,
					Title:    "更改发表标题",
					Content:  "更改发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
				}, publishArt)
			},
			art: Article{
				Id:      2,
				Title:   "更改发表标题",
				Content: "更改发表内容",
			},
			wantCode: 200,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 2,
			},
		},
		{
			name: "更新帖子，并且重新发表", //更新完成 覆盖上次的发表
			before: func(t *testing.T) {
				err := a.db.Create(article.Article{
					Id:       3,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
				err = a.db.Create(article.PublishedArticle{
					Id:       3,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//因为是先更改的
				var art article.Article
				err := a.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 234)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "更改再次发表标题",
					Content:  "更改再次发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
					Ctime:    456,
				}, art)
				//线上库
				var publishArt article.PublishedArticle
				err = a.db.Where("id = ?", 3).First(&publishArt).Error
				assert.NoError(t, err)
				assert.True(t, publishArt.Ctime > 0)
				assert.True(t, publishArt.Utime > 0)
				publishArt.Ctime = 0
				publishArt.Utime = 0
				assert.Equal(t, article.PublishedArticle{
					Id:       3,
					Title:    "更改再次发表标题",
					Content:  "更改再次发表内容",
					AuthorId: 123,
					Status:   uint8(domain.ArticleStatusPublished),
				}, publishArt)
			},
			art: Article{
				Id:      3,
				Title:   "更改再次发表标题",
				Content: "更改再次发表内容",
			},
			wantCode: 200,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 3,
			},
		},
		{
			name: "更新别人的帖子，并且发表失败",
			before: func(t *testing.T) {
				err := a.db.Create(article.Article{
					Id:       4,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 567,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
				err = a.db.Create(article.PublishedArticle{
					Id:       4,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 567,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//因为是先更改的
				var art article.Article
				err := a.db.Where("id = ?", 4).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       4,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 567,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}, art)
				//线上库
				var publishArt article.PublishedArticle
				err = a.db.Where("id = ?", 4).First(&publishArt).Error
				assert.NoError(t, err)
				assert.Equal(t, article.PublishedArticle{
					Id:       4,
					Title:    "发表标题",
					Content:  "发表内容",
					AuthorId: 567,
					Status:   uint8(domain.ArticleStatusUnpublished),
					Ctime:    456,
					Utime:    234,
				}, publishArt)
			},
			art: Article{
				Id:      4,
				Title:   "更改再次发表标题",
				Content: "更改再次发表内容",
			},
			wantCode: 200,
			wantRes: Result[int64]{
				Msg:  "系统错误",
				Code: 5,
			},
		},
	}
	for _, v := range testCase {
		t.Run(v.name, func(t *testing.T) {
			v.before(t)
			reqBody, err := json.Marshal(v.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			a.server.ServeHTTP(resp, req)
			assert.Equal(t, v.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var r Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&r)
			assert.NoError(t, err)
			//验证结果
			assert.Equal(t, v.wantRes, r)
			v.after(t)
		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}
