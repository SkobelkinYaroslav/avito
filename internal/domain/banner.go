package domain

import "time"

// GetUserBannerRequest для GET /user_banner
// Параметры запроса для получения баннера пользователя.
type GetUserBannerRequest struct {
	Banner
	UseLastRevision bool `json:"use_last_revision"`
}

// GetBannersRequest для GET /banner
// Параметры запроса для получения всех баннеров с фильтрацией.
type GetBannersRequest struct {
	TagIDs    []int `json:"tag_ids,omitempty"`
	FeatureID int   `json:"feature_id,omitempty"`
	Limit     int
	Offset    int
}
type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Banner struct {
	BannerID  int       `json:"banner_id,omitempty"`
	TagIDs    []int     `json:"tag_ids,omitempty"`
	FeatureID int       `json:"feature_id,omitempty"`
	Content   Content   `json:"content,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
