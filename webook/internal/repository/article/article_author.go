// 制作库
package article

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/dao/article"
)

type RepositoryArticleAuthor interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
type repositoryArticleAuthor struct {
	dao article.DaoArticleAuthor
}

func NewRepositoryArticleAuthor(dao article.DaoArticleAuthor) RepositoryArticleAuthor {
	return &repositoryArticleAuthor{
		dao: dao,
	}
}
func (a *repositoryArticleAuthor) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.dao.Create(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
func (a *repositoryArticleAuthor) Update(ctx context.Context, art domain.Article) error {
	return a.dao.Update(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
