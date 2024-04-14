package repository

import (
	"avito/internal/domain"
	errGroup "avito/internal/error"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type BannerRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewBannerRepository(db *sql.DB, redis *redis.Client) *BannerRepository {
	return &BannerRepository{
		db:    db,
		redis: redis,
	}
}

func (b BannerRepository) GetUserBannerRepo(req domain.GetUserBannerRequest) (domain.Banner, error) {
	var banner domain.Banner
	log.Println(req)
	redisQuery := fmt.Sprintf("tag_ids:%v&feature_id:%v", req.TagIDs, req.FeatureID)
	sqlQuery := "SELECT title, text, url FROM banner WHERE feature_id = $1 AND tag_ids = $2"

	getFromDb := func() error {
		row := b.db.QueryRow(sqlQuery, req.FeatureID, pq.Array(req.TagIDs))
		err := row.Scan(&banner.Content.Title, &banner.Content.Text, &banner.Content.URL)
		if err != nil {
			return err
		}

		return nil
	}

	if req.UseLastRevision {
		err := getFromDb()
		if err == sql.ErrNoRows {
			return domain.Banner{}, err
		}
		if err != nil {
			return domain.Banner{}, err
		}
	} else {
		ctx := context.Background()
		val, err := b.redis.Get(ctx, redisQuery).Result()

		if err == redis.Nil {
			err = getFromDb()
			if err == sql.ErrNoRows {
				return domain.Banner{}, err
			}
			if err != nil {
				return domain.Banner{}, err
			}

			setBanner, err := json.Marshal(banner)
			if err != nil {
				return domain.Banner{}, err
			}

			err = b.redis.Set(ctx, redisQuery, setBanner, 5*time.Minute).Err()
			if err != nil {
				return domain.Banner{}, err
			}

		}

		if err != nil {
			return domain.Banner{}, err
		}

		err = json.Unmarshal([]byte(val), &banner)
		if err != nil {
			return domain.Banner{}, err
		}

	}

	return banner, nil
}

func (b BannerRepository) GetAllBannersRepo(req domain.GetBannersRequest) ([]domain.Banner, error) {
	var banners []domain.Banner
	query := "SELECT * FROM banner WHERE feature_id = $1 OR $2 && tag_ids LIMIT $3 OFFSET $4"

	rows, err := b.db.Query(query, req.FeatureID, pq.Array(req.TagIDs), req.Limit, req.Offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var banner domain.Banner
		var tagIDs pq.StringArray

		err = rows.Scan(
			&banner.BannerID,
			&tagIDs,
			&banner.FeatureID,
			&banner.Content.Title,
			&banner.Content.Text,
			&banner.Content.URL,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		banner.TagIDs = make([]int, len(tagIDs))
		for i, v := range tagIDs {
			banner.TagIDs[i], _ = strconv.Atoi(v)
		}

		banners = append(banners, banner)
	}
	return banners, nil
}

func (b BannerRepository) PostBannerRepo(req domain.Banner) (domain.Banner, error) {
	redisQuery := fmt.Sprintf("tag_ids:%v&feature_id:%v", req.TagIDs, req.FeatureID)

	query := `
INSERT INTO banner(
    tag_ids, 
    feature_id, 
    title, 
    text, 
    url, 
    is_active, 
    created_at, 
    updated_at
) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8) 
RETURNING banner_id
`
	err := b.db.QueryRow(
		query,
		pq.Array(req.TagIDs),
		req.FeatureID,
		req.Content.Title,
		req.Content.Text,
		req.Content.URL,
		req.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&req.BannerID)

	if err != nil {
		return domain.Banner{}, err
	}

	model, err := json.Marshal(req)
	if err != nil {
		return domain.Banner{}, err
	}

	err = b.redis.Set(context.Background(), redisQuery, model, 5*time.Minute).Err()
	if err != nil {
		return domain.Banner{}, err
	}

	return req, nil
}

func (b BannerRepository) PatchBannerRepo(req domain.Banner) error {
	query := `
        UPDATE banner
        SET
            tag_ids = $1,
            feature_id = $2,
            title = $3,
            text = $4,
            url = $5,
            is_active = $6,
            updated_at = $7
        WHERE banner_id = $8
        `
	result, err := b.db.Exec(
		query,
		pq.Array(req.TagIDs),
		req.FeatureID,
		req.Content.Title,
		req.Content.Text,
		req.Content.URL,
		req.IsActive,
		time.Now(),
		req.BannerID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errGroup.NotFound
	}

	return nil
}

func (b BannerRepository) DeleteBannerRepo(id int) error {
	query := "DELETE FROM banner WHERE banner_id = $1"
	result, err := b.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errGroup.NotFound
	}

	return nil
}
