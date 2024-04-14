package handler

import (
	"avito/internal/domain"
	errGroup "avito/internal/error"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h Handler) GetUserBannerHandler(c *gin.Context) {
	var req domain.GetUserBannerRequest

	tagIDStr := c.Query("tag_id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag_id"})
		return
	}

	featureIDStr := c.Query("feature_id")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feature_id"})
		return
	}

	useLastRevisionStr := c.Query("use_last_revision")
	useLastRevision, err := strconv.ParseBool(useLastRevisionStr)
	if err != nil {
		useLastRevision = false
	}

	req.TagIDs = []int{tagID}
	req.FeatureID = featureID
	req.UseLastRevision = useLastRevision

	log.Println(req)

	response, err := h.BannerService.GetUserBannerService(req)

	if err == errGroup.NotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "no banner found"})
		return
	}

	log.Println(err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get user banner"})
		return
	}

	user, _ := c.Get("user")
	if user.(domain.AuthStruct).IsAdmin || response.IsActive {
		c.JSON(http.StatusOK, gin.H{"message": response})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "no banner found"})

}

func (h Handler) GetAllBannersHandler(c *gin.Context) {
	var req domain.GetBannersRequest

	featureIDStr := c.Query("feature_id")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		featureID = -1
	}
	req.FeatureID = featureID

	tagIDStr := c.Query("tag_id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		req.TagIDs = []int{-1}
	} else {
		req.TagIDs = []int{tagID}
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	req.Limit = limit

	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	req.Offset = offset

	response, err := h.BannerService.GetAllBannersService(req)
	if err == errGroup.NotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "no banners found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

func (h Handler) PostBannerHandler(c *gin.Context) {
	var req domain.Banner

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.BannerService.PostBannerService(req)
	if err == errGroup.AlreadyExists {
		c.JSON(http.StatusConflict, gin.H{"message": "banner already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h Handler) PatchBannerHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req domain.Banner

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.BannerID = idInt

	err = h.BannerService.PatchBannerService(req)
	if err == errGroup.NotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "no banner found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "banner updated"})
}

func (h Handler) DeleteBannerHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	req := domain.Banner{
		BannerID: idInt,
	}

	err = h.BannerService.DeleteBannerService(req)
	if err == errGroup.NotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "no banner found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "banner deleted"})
}

func (h Handler) DeleteBannerByDataHandler(c *gin.Context) {
	var req domain.Banner

	tagIDStr := c.Query("tag_id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		tagID = -1
	}

	featureIDStr := c.Query("feature_id")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		featureID = -1
	}
	req.FeatureID = featureID
	req.TagIDs = []int{tagID}

	go func(req domain.Banner) {
		_ = h.BannerService.DeleteBannerService(req)
	}(req)

	c.JSON(http.StatusOK, gin.H{"message": "deleting banner"})

}
