package media

import "github.com/gin-gonic/gin"

func RouterMedia(r *gin.RouterGroup) {

	sv := NewMediaService()
	ctrl := NewMediaController(sv)
	media := r.Group("/media")
	{
		media.GET("/a", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Media route"})
		})
		media.GET("/", ctrl.GetAllMedia)
		media.POST("/single", ctrl.UploadSingle)
		media.POST("/multiple", ctrl.UploadMultiple)
	}
}
