package repository_test

import (
	"avito/internal/domain"
	errGroup "avito/internal/error"
	"avito/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func Test_GetExistingBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBannerRepository(db, nil)

	expectedBanner := domain.Banner{
		Content: domain.Content{
			Title: "Test Title",
			Text:  "Test Text",
			URL:   "http://test.url",
		},
	}

	rows := sqlmock.NewRows([]string{"title", "text", "url"}).
		AddRow(expectedBanner.Content.Title, expectedBanner.Content.Text, expectedBanner.Content.URL)

	mock.ExpectQuery("^SELECT (.+) FROM banner WHERE feature_id = \\$1 AND tag_ids = \\$2$").
		WithArgs(1, pq.Array([]int{1, 2})).
		WillReturnRows(rows)

	banner := domain.Banner{
		FeatureID: 1,
		TagIDs:    []int{1, 2},
	}

	resultBanner, err := repo.GetUserBannerRepo(domain.GetUserBannerRequest{Banner: banner, UseLastRevision: true})

	assert.NoError(t, err)
	assert.Equal(t, expectedBanner, resultBanner)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_GetNonExistingBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBannerRepository(db, nil)

	mock.ExpectQuery("^SELECT (.+) FROM banner WHERE feature_id = \\$1 AND tag_ids = \\$2$").
		WithArgs(1, pq.Array([]int{1, 2})).
		WillReturnError(errGroup.NotFound)

	banner := domain.Banner{
		FeatureID: 1,
		TagIDs:    []int{1, 2},
	}

	_, err = repo.GetUserBannerRepo(domain.GetUserBannerRequest{Banner: banner, UseLastRevision: true})

	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_GetExistingBannerFromCache(t *testing.T) {
	db, mock := redismock.NewClientMock()

	repo := repository.NewBannerRepository(nil, db)

	expectedBanner := domain.Banner{
		Content: domain.Content{
			Title: "Test Title",
			Text:  "Test Text",
			URL:   "http://test.url",
		},
	}
	bannerBytes, _ := json.Marshal(expectedBanner)
	mock.ExpectGet("tag_ids:[1 2]&feature_id:1").SetVal(string(bannerBytes))

	banner := domain.Banner{
		FeatureID: 1,
		TagIDs:    []int{1, 2},
	}

	resultBanner, err := repo.GetUserBannerRepo(domain.GetUserBannerRequest{Banner: banner, UseLastRevision: false})

	assert.NoError(t, err)
	assert.Equal(t, expectedBanner, resultBanner)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllBannersRepo(t *testing.T) {
	banners := []domain.Banner{
		{
			BannerID:  1,
			TagIDs:    []int{1, 2},
			FeatureID: 1,
			Content:   domain.Content{Title: "Sample Title 1", Text: "Sample Text 1", URL: "http://example1.com"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			BannerID:  2,
			TagIDs:    []int{3, 4},
			FeatureID: 2,
			Content:   domain.Content{Title: "Sample Title 2", Text: "Sample Text 2", URL: "http://example2.com"},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"banner_id", "tag_ids", "feature_id", "title", "text", "url", "is_active", "created_at", "updated_at"})
	for _, banner := range banners {
		rows.AddRow(banner.BannerID, pq.Array(banner.TagIDs), banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive, banner.CreatedAt, banner.UpdatedAt)
	}

	req := domain.GetBannersRequest{
		FeatureID: 1,
		TagIDs:    []int{1, 2},
		Limit:     10,
		Offset:    0,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM banner WHERE feature_id = $1 OR $2 && tag_ids LIMIT $3 OFFSET $4")).
		WithArgs(req.FeatureID, pq.Array(req.TagIDs), req.Limit, req.Offset).
		WillReturnRows(rows)

	b := repository.NewBannerRepository(db, nil)

	result, err := b.GetAllBannersRepo(req)

	assert.NoError(t, err)

	assert.Equal(t, banners, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostBannerRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	redisDB, redisMock := redismock.NewClientMock()

	repo := repository.NewBannerRepository(db, redisDB)

	expectedBanner := domain.Banner{
		BannerID:  1,
		TagIDs:    []int{1, 2},
		FeatureID: 1,
		Content: domain.Content{
			Title: "Test Title",
			Text:  "Test Text",
			URL:   "http://test.url",
		},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("^INSERT INTO banner").
		WithArgs(pq.Array(expectedBanner.TagIDs), expectedBanner.FeatureID, expectedBanner.Content.Title, expectedBanner.Content.Text, expectedBanner.Content.URL, expectedBanner.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"banner_id"}).AddRow(1))

	bannerBytes, _ := json.Marshal(expectedBanner)
	redisMock.ExpectSet(fmt.Sprintf("tag_ids:%v&feature_id:%v", expectedBanner.TagIDs, expectedBanner.FeatureID), bannerBytes, 5*time.Minute).SetVal("OK")

	resultBanner, err := repo.PostBannerRepo(expectedBanner)

	assert.NoError(t, err)
	assert.Equal(t, expectedBanner, resultBanner)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, redisMock.ExpectationsWereMet())
}

func TestPatchBannerRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBannerRepository(db, nil)

	expectedBanner := domain.Banner{
		BannerID:  1,
		TagIDs:    []int{1, 2},
		FeatureID: 1,
		Content: domain.Content{
			Title: "Test Title",
			Text:  "Test Text",
			URL:   "http://test.url",
		},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec("^UPDATE banner").
		WithArgs(pq.Array(expectedBanner.TagIDs), expectedBanner.FeatureID, expectedBanner.Content.Title, expectedBanner.Content.Text, expectedBanner.Content.URL, expectedBanner.IsActive, sqlmock.AnyArg(), expectedBanner.BannerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.PatchBannerRepo(expectedBanner)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBannerRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewBannerRepository(db, nil)

	req := domain.Banner{
		BannerID: 1,
	}

	mock.ExpectExec("^DELETE FROM banner WHERE banner_id = ?").
		WithArgs(req.BannerID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteBannerRepo(req)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
