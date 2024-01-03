package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	domain2 "goFoundation/webook/interactive/domain"
	"goFoundation/webook/interactive/service"
	"goFoundation/webook/internal/domain"
	svcmocks "goFoundation/webook/internal/service/mocks"
	"testing"
	"time"
)

func Test_rangingService(t *testing.T) {
	now := time.Now()
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService)
		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "计算成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, service.InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				artSvc.EXPECT().ListPub(gomock.Any(), now, 0, 3).Return([]domain.Article{
					{Id: 1, Utime: now, Ctime: now},
					{Id: 2, Utime: now, Ctime: now},
					{Id: 3, Utime: now, Ctime: now},
				}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), now, 3, 3).Return([]domain.Article{}, nil)
				interSvc := svcmocks.NewMockInteractiveService(ctrl)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2, 3}).Return(map[int64]domain2.Interactive{
					1: {BizId: 1, LikeCnt: 1},
					2: {BizId: 2, LikeCnt: 2},
					3: {BizId: 3, LikeCnt: 3},
				}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).Return(map[int64]domain2.Interactive{}, nil)
				return artSvc, interSvc
			},
			wantArts: []domain.Article{
				{Id: 3, Utime: now, Ctime: now},
				{Id: 2, Utime: now, Ctime: now},
				{Id: 1, Utime: now, Ctime: now},
			},
		},
	}
	for _, v := range testCase {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			art, inter := v.mock(ctrl)
			svc := NewRangingService(art, inter).(*rangingService)
			// 为了测试
			svc.Limit = 3
			svc.N = 3
			svc.ScoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			res, err := svc.topN(context.Background())
			assert.Equal(t, v.wantErr, err)
			assert.Equal(t, v.wantArts, res)
		})
	}
}
