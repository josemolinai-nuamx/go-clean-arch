package article

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/josemolinai-nuamx/go-clean-arch/domain"
)

// ArticleRepository defines the contract for data access operations on articles.
// This is an INTERFACE (not concrete implementation) to achieve DEPENDENCY INVERSION:
//   - Service depends on this interface, not on MySQL or PostgreSQL directly
//   - For tests: inject a mock repository (no real database needed)
//   - For production: inject the real MySQL repository
//
// The actual SQL queries are in: internal/repository/mysql/article.go
//
//go:generate mockery --name ArticleRepository
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetByTitle(ctx context.Context, title string) (domain.Article, error)
	Update(ctx context.Context, ar *domain.Article) error
	Store(ctx context.Context, a *domain.Article) error
	Delete(ctx context.Context, id int64) error
}

// AuthorRepository defines the contract for fetching author data.
// Same pattern as ArticleRepository: interface-based for testability.
//
//go:generate mockery --name AuthorRepository
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Author, error)
}

// Service orchestrates the business logic (use cases) for articles.
// RESPONSIBILITY:
//   - Execute business rules (validation, conflict checking, etc.)
//   - Coordinate between repositories and domain entities
//   - Transform errors into domain-level errors
//
// DOESN'T do:
//   - Direct SQL queries (that's Repository's job)
//   - HTTP handling (that's REST Handler's job)
type Service struct {
	articleRepo ArticleRepository
	authorRepo  AuthorRepository
}

// NewService creates a new Service with injected dependencies.
// Dependency Injection pattern: The service receives its dependencies,
// doesn't create them itself. This makes it:
//   - Testable: inject mock repositories
//   - Flexible: swap repositories without changing Service code
func NewService(a ArticleRepository, ar AuthorRepository) *Service {
	return &Service{
		articleRepo: a,
		authorRepo:  ar,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *Service) fillAuthorDetails(ctx context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}
	// Get the author's id
	mapAuthors := map[int64]domain.Author{}
	authorIDs := make([]int64, 0, len(data))

	for _, article := range data { //nolint
		authorID := article.Author.ID
		if _, ok := mapAuthors[authorID]; !ok {
			mapAuthors[authorID] = domain.Author{}
			authorIDs = append(authorIDs, authorID)
		}
	}
	// Using goroutine to fetch the author's detail
	for _, authorID := range authorIDs {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			mu.Lock()
			mapAuthors[res.ID] = res
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data { //nolint
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *Service) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *Service) Update(ctx context.Context, ar *domain.Article) (err error) {
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *Service) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *Service) Store(ctx context.Context, m *domain.Article) (err error) {
	existedArticle, _ := a.GetByTitle(ctx, m.Title) // ignore if any error
	if existedArticle != (domain.Article{}) {
		return domain.ErrConflict
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *Service) Delete(ctx context.Context, id int64) (err error) {
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
