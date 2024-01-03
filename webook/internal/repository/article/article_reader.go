// 线上库
package article

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/dao/article"
)

type RepositoryArticleReader interface {
	//有id 就是修该 没有就是创建
	Save(ctx context.Context, art domain.Article) error
}
type repositoryArticleReader struct {
	dao article.DaoArticleReader
}

func NewRepositoryArticleReader(dao article.DaoArticleReader) RepositoryArticleReader {
	return &repositoryArticleReader{
		dao: dao,
	}
}

func (r *repositoryArticleReader) Save(ctx context.Context, art domain.Article) error {
	return r.dao.Save(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
