package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/service"
	ijwt "goFoundation/webook/internal/web/jwt"
	"goFoundation/webook/pkg/ginx"
	"goFoundation/webook/pkg/logger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc   service.ArticleService
	l     logger.LoggerV1
	inSvc intrv1.InteractiveServiceClient
	biz   string
}

func NewArticleHandlers(svc service.ArticleService, inSvc intrv1.InteractiveServiceClient,
	l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc:   svc,
		l:     l,
		inSvc: inSvc,
		biz:   "article",
	}
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	//作者
	g := s.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	g.POST("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](a.List))
	g.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))
	//读者
	pub := s.Group("/pub")
	pub.GET("/:id", ginx.WrapToken[ijwt.UserClaims](a.PubDetail))
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq, ijwt.UserClaims](a.Like))
	pub.POST("/collect", ginx.WrapBodyAndToken[CollectReq, ijwt.UserClaims](a.Collect))
}

// 收藏
func (a *ArticleHandler) Collect(ctx *gin.Context, c CollectReq, u ijwt.UserClaims) (ginx.Result, error) {
	_, err := a.inSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz: a.biz, Cid: c.Cid, BizId: c.Id, Uid: u.Uid,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{Msg: "OK"}, nil
}

// 点赞
func (a *ArticleHandler) Like(ctx *gin.Context, l LikeReq, c ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if l.Like {
		_, err = a.inSvc.Like(ctx, &intrv1.LikeRequest{})

	} else {
		_, err = a.inSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz: a.biz, BizId: l.Id, Uid: c.Uid,
		})
	}
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统繁忙",
		}, err
	}
	return ginx.Result{Msg: "ok"}, nil
}

// 创作者查id 缓存方案
func (a *ArticleHandler) Detail(ctx *gin.Context, c ijwt.UserClaims) (ginx.Result, error) {
	//从前端获取id
	idStr := ctx.Param("id")
	// 一个string类型的字符串 base转化后所需要的进制 bitsize能且只能输出int64类型数据
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		}, err
	}
	art, err := a.svc.Detail(ctx, id)
	if err != nil {
		//a.l.Error("获得文章信息失败", logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	if art.Author.Id != c.Uid {
		//ctx.JSON(http.StatusOK)
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		//a.l.Error("非法访问文章，创作者 ID 不匹配",
		//	logger.Int64("uid", usr.Id))
		return ginx.Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", c.Uid)
	}
	return ginx.Result{Data: ArticleVO{
		Id:    art.Id,
		Title: art.Title,
		// 不需要这个摘要信息
		//Abstract: art.Abstract(),
		Status:  art.Status.ToUint8(),
		Content: art.Content,
		// 这个是创作者看自己的文章列表，也不需要这个字段
		//Author: art.Author
		Ctime: art.Ctime.Format(time.DateTime),
		Utime: art.Utime.Format(time.DateTime),
	}}, nil
}

// 读者查id 跨领域操作 你需要拿到作者的名字
func (a *ArticleHandler) PubDetail(ctx *gin.Context, c ijwt.UserClaims) (ginx.Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		}, fmt.Errorf("查询文章详情的 ID %s 不正确, %w", idStr, err)

	}
	//为什么使用这个eg.go 因为我必须得到文章和收藏，阅读，点赞的数的数量 然后才能返回给前端
	var (
		eg      errgroup.Group
		getResp *intrv1.GetResponse
		art     domain.Article
	)
	eg.Go(func() error {
		var er error
		art, er = a.svc.PubDetail(ctx, id, c.Uid)
		return er

	})
	eg.Go(func() error {
		var er error
		//为了用户更好的体验 我需要让用户看到 收藏，阅读，点赞的数量 我需要在这回显
		getResp, er = a.inSvc.Get(ctx, &intrv1.GetRequest{
			Biz: a.biz, BizId: id, Uid: c.Uid,
		})
		return er
	})
	//等待上面的2个eg.go执行完事才继续向下执行
	err = eg.Wait()
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, fmt.Errorf("获取文章信息失败 %w", err)
	}

	//计算浏览数量
	go func() {
		//虽然这里异步了 但还是对数据库造成很大的压力 随意用kafka来解耦
		//数据库什么压力？ 在接口里面全是读数据库 或者直接命中 所以用kafka来进行写操作
		_, er := a.inSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
			Biz: a.biz, BizId: art.Id,
		})
		if er != nil {
			a.l.Error("获得文章信息失败", logger.Error(er))

		}
	}()
	return ginx.Result{Data: ArticleVO{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status.ToUint8(),
		Content: art.Content,
		// 要把作者信息带出去
		Author: art.Author.Name,
		Ctime:  art.Ctime.Format(time.DateTime),
		Utime:  art.Utime.Format(time.DateTime),
		//阅读。收藏。点赞
		ReadCnt:    getResp.Intr.ReadCnt,
		CollectCnt: getResp.Intr.CollectCnt,
		LikeCnt:    getResp.Intr.LikeCnt,
		Liked:      getResp.Intr.Liked,
		Collected:  getResp.Intr.Collected,
	}}, nil
}

// 搜索
func (a *ArticleHandler) List(ctx *gin.Context, r ListReq, c ijwt.UserClaims) (ginx.Result, error) {
	res, err := a.svc.List(ctx, c.Uid, r.Offset, r.Limit)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	// 在列表页，不显示全文，只显示一个"摘要"
	// 比如说，简单的摘要就是前几句话
	// 强大的摘要是 AI 帮你生成的
	return ginx.Result{
		//domain.Article转换成ArticleVO
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Id:       src.Id,
				Title:    src.Title,
				Status:   src.Status.ToUint8(),
				Abstract: src.Abstract(), //摘要
				Ctime:    src.Ctime.Format(time.DateTime),
				Utime:    src.Utime.Format(time.DateTime),
				// 这个是创作者看自己的文章列表，也不需要这个字段
				//Author: src.Author
			}
		}),
	}, nil
}

// 仅自己可见的状态
func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: claims.Uid,
		Author: domain.Author{
			Id: req.Id,
		},
	})
	if err != nil {
		a.l.Error("设置为尽自己可见失败", logger.Error(err),
			logger.Field{Key: "id", Value: req.Id})
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

// 发表
func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
		Msg:  "ok",
	})
}

// 创建和修改
func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存数据失败", logger.Field{Key: "error", Value: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
		Msg:  "保存成功",
	})
}
