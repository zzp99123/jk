package web

//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/mock/gomock"
//	"goFoundation/webook/internal/domain"
//	"goFoundation/webook/internal/service"
//	svcmocks "goFoundation/webook/internal/service/mocks"
//	ijwt "goFoundation/webook/internal/web/jwt"
//	"goFoundation/webook/pkg/logger"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestArticleHandler_Publish(t *testing.T) {
//	testCases := []struct {
//		name     string
//		mock     func(ctrl *gomock.Controller) service.ArticleService
//		reqBody  string
//		wantCode int
//		wantRes  Result
//	}{
//		{
//			name: "新建并发表",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Title:   "标题",
//					Content: "内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(1), nil)
//				return svc
//			},
//			reqBody: `{
//
//"title":"标题",
//"content":"内容"
//}`,
//			wantCode: 200,
//			wantRes: Result{
//				Data: float64(1),
//				Msg:  "ok",
//			},
//		},
//		{
//			name: "pulish失败",
//			mock: func(ctrl *gomock.Controller) service.ArticleService {
//				svc := svcmocks.NewMockArticleService(ctrl)
//				svc.EXPECT().Publish(gomock.Any(), domain.Article{
//					Title:   "标题",
//					Content: "内容",
//					Author: domain.Author{
//						Id: 123,
//					},
//				}).Return(int64(1), nil)
//				return svc
//			},
//			reqBody: `{
//
//"title":"标题",
//"content":"内容"
//}`,
//			wantCode: 200,
//			wantRes: Result{
//				Code: 5,
//				Msg:  "系统错误",
//			},
//		},
//	}
//	for _, v := range testCases {
//		t.Run(v.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			server := gin.Default()
//			server.Use(func(ctx *gin.Context) {
//				ctx.Set("claims", ijwt.UserClaims{Uid: 123})
//			})
//			a := NewArticleHandler(v.mock(ctrl), &logger.NopLogger{})
//			a.RegisterRoutes(server)
//			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer([]byte(v.reqBody)))
//			require.NoError(t, err)
//			req.Header.Set("Content-Type", "application/json")
//			resp := httptest.NewRecorder()
//			server.ServeHTTP(resp, req)
//			assert.Equal(t, v.wantCode, resp.Code)
//			if resp.Code != 200 {
//				return
//			}
//			var webRes Result
//			err = json.NewDecoder(resp.Body).Decode(&webRes)
//			require.NoError(t, err)
//			assert.Equal(t, v.wantRes, webRes)
//		})
//	}
//}
