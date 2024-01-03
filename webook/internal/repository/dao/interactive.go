package dao

//
//import (
//	"context"
//	"gorm.io/gorm"
//	"gorm.io/gorm/clause"
//	"time"
//)
//
//type InteractiveDao interface {
//	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
//	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error
//	Like(ctx context.Context, biz string, bizId int64, id int64) error
//	CancelLike(ctx context.Context, biz string, bizId int64, id int64) error
//	GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (Interactive, error)
//	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
//	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
//	InsertCollectionBiz(ctx context.Context, u UserCollectionBiz) error
//	GetByIds(ctx context.Context, biz string, bizId []int64) ([]Interactive, error)
//}
//type interactiveDao struct {
//	db *gorm.DB
//}
//
//func NewInteractiveDao(db *gorm.DB) InteractiveDao {
//	return &interactiveDao{
//		db: db,
//	}
//}
//
//// 热搜中获取id
//func (i *interactiveDao) GetByIds(ctx context.Context, biz string, bizId []int64) ([]Interactive, error) {
//	var res []Interactive
//	err := i.db.WithContext(ctx).Where("biz= ?AND bizId=?", biz, bizId).Find(&i).Error
//	return res, err
//}
//
//// 搜索获取收藏，点赞，阅读的信息
//func (i *interactiveDao) GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (Interactive, error) {
//	var in Interactive
//	err := i.db.WithContext(ctx).Where("biz = ? AND biz_id =?", biz, bizId).First(&in).Error
//	return in, err
//}
//
//// 批量记录浏览数量
//func (i *interactiveDao) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
//	// 可以用 map 合并吗？
//	// 看情况。如果一批次里面，biz 和 bizid 都相等的占很多，那么就map 合并，性能会更好
//	// 不然你合并了没有效果
//
//	// 为什么快？
//	// A：十条消息调用十次 IncrReadCnt，
//	// B 就是批量
//	// 事务本身的开销，A 是 B 的十倍
//	// 刷新 redolog, undolog, binlog 到磁盘，A 是十次，B 是一次
//	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		txDao := NewInteractiveDao(tx)
//		for a := range biz {
//			err := txDao.IncrReadCnt(ctx, biz[a], bizId[a])
//			if err != nil {
//				// 记个日志就拉到
//				// 也可以 return err
//				return err
//			}
//		}
//		return nil
//	})
//}
//func (i *interactiveDao) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
//	return i.incrReadCnt(i.db.WithContext(ctx), biz, bizId)
//}
//
//// 记录浏览数量
//func (i *interactiveDao) incrReadCnt(tx *gorm.DB, biz string, bizId int64) error {
//	//如果有则更新数据 没有的话就创建数据 upsert语句
//	now := time.Now().UnixMilli()
//	return tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{
//		"read_cnt": gorm.Expr("read_cnt + 1"),
//		"utime":    now,
//	}),
//	}).Create(&Interactive{
//		Biz:     biz,
//		BizId:   bizId,
//		ReadCnt: 1,
//		Utime:   now,
//		Ctime:   now,
//	}).Error
//}
//
//// 搜素获取点赞信息
//func (i *interactiveDao) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
//	var u UserLikeBiz
//	err := i.db.WithContext(ctx).Where("biz_id = ? AND biz =? AND uid = ?", bizId, biz, uid).First(&u).Error
//	return u, err
//}
//
//// 点赞
//func (i *interactiveDao) Like(ctx context.Context, biz string, bizId, uid int64) error {
//	now := time.Now().UnixMilli()
//	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		err := tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{
//			"Status": 1,
//			"utime":  now,
//		}),
//		}).Create(&UserLikeBiz{
//			BizId:  bizId,
//			Biz:    biz,
//			Uid:    uid,
//			Status: 1,
//			Utime:  now,
//			Ctime:  now,
//		}).Error
//		if err != nil {
//			return err
//		}
//		return tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{
//			"like_cnt": gorm.Expr("like_cnt + 1"),
//			"utime":    now,
//		}),
//		}).Create(&Interactive{
//			Biz:     biz,
//			BizId:   bizId,
//			LikeCnt: 1,
//			Utime:   now,
//			Ctime:   now,
//		}).Error
//	})
//	return err
//}
//
//// 取消点赞
//func (i *interactiveDao) CancelLike(ctx context.Context, biz string, bizId, uid int64) error {
//	now := time.Now().UnixMilli()
//	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		err := tx.Model(&UserLikeBiz{}).
//			Where("biz =? AND biz_id = ? AND uid = ?", biz, bizId, uid).
//			Updates(map[string]any{
//				"status": 0,
//				"utime":  now,
//			}).Error
//		if err != nil {
//			return err
//		}
//		return i.db.WithContext(ctx).Clauses(clause.OnConflict{
//			DoUpdates: clause.Assignments(map[string]any{
//				"like_cnt": gorm.Expr("`like_cnt`-1"),
//				"utime":    now,
//			}),
//		}).Create(&Interactive{
//			LikeCnt: 1,
//			Ctime:   now,
//			Utime:   now,
//			Biz:     biz,
//			BizId:   bizId,
//		}).Error
//	})
//	return err
//}
//
//// 搜素获取收藏信息
//func (i *interactiveDao) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
//	var u UserCollectionBiz
//	err := i.db.WithContext(ctx).Where("biz=? AND biz_id = ? AND uid = ?", biz, bizId, uid).First(&u).Error
//	return u, err
//}
//
//// 插入收藏记录，并且更新计数
//func (i *interactiveDao) InsertCollectionBiz(ctx context.Context, u UserCollectionBiz) error {
//	now := time.Now().UnixMilli()
//	u.Ctime = now
//	u.Utime = now
//	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		err := i.db.WithContext(ctx).Create(&u).Error
//		if err != nil {
//			return err
//		}
//		return tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{
//			"collect_cnt": gorm.Expr("collect_cnt + 1"),
//			"utime":       now,
//		}),
//		}).Create(&Interactive{
//			Biz:        u.Biz,
//			BizId:      u.BizId,
//			Ctime:      now,
//			Utime:      now,
//			CollectCnt: 1,
//		}).Error
//	})
//}
//
//// SELECT * FROM (SELECT biz, biz_id, COUNT(*) as cnt FROM `interactives` GROUP BY biz, biz_id) ORDER BY cnt LIMIT 100
//// 实时查找，性能贼差，上面这个语句，就是全表扫描，
//// 高性能，我不要求准确性 面试标准答案：用 zset 但是，面试标准答案不够有特色，烂大街了
//// 你可以考虑别的方案
//// 1. 定时计算
//// 1.1 定时计算 + 本地缓存
//// 2. 优化版的 zset，定时筛选 + 实时 zset 计算
//// 还要别的方案你们也可以考虑
//
//// 一个主表专门计算收藏，点赞，阅读总数量的 其他的副表就是收藏，点赞，阅读单独用的表
//type Interactive struct {
//	Id         int64  `gorm:"primaryKey,autoIncrement"`
//	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
//	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
//	Ctime      int64
//	Utime      int64
//	ReadCnt    int64 //阅读量的数量
//	CollectCnt int64 //收藏的数量
//	LikeCnt    int64 //点赞的数量
//}
//
//// 点赞副表
//type UserLikeBiz struct {
//	Id int64 `gorm:"primaryKey,autoIncrement"`
//	//联合索引
//	BizId  int64  `gorm:"uniqueIndex:biz_type_id_uid"`
//	Biz    string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
//	Uid    int64  `gorm:"uniqueIndex:biz_type_id_uid"`
//	Ctime  int64
//	Utime  int64
//	Status uint8 // 依旧是只在 DB 层面生效的状态 1- 有效，0-无效。软删除的用法
//}
//
//// 收藏夹
//type Collection struct {
//	Id    int64
//	Name  string
//	Uid   int64
//	Ctime int64
//	Utime int64
//}
//
//// UserCollectionBiz 收藏的东西
//type UserCollectionBiz struct {
//	Id int64 `gorm:"primaryKey,autoIncrement"`
//	//收藏夹的id
//	Cid int64 `gorm:"index"`
//	//联合索引
//	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
//	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
//	Uid   int64  `gorm:"uniqueIndex:biz_type_id_uid"`
//	Ctime int64
//	Utime int64
//}
