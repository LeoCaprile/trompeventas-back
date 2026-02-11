package comments

import (
	"restorapp/modules/auth"

	"github.com/gin-gonic/gin"
)

func CommentsController(router *gin.Engine) {
	// Public with optional auth (to include user's votes)
	router.GET("/products/:id/comments", auth.OptionalAuthMiddleware(), getCommentsHandler)

	// Authenticated routes for comment CRUD
	productComments := router.Group("/products/:id/comments")
	productComments.Use(auth.AuthMiddleware())
	productComments.POST("/", createCommentHandler)

	// Comment-level routes (delete, vote)
	commentRoutes := router.Group("/comments")
	commentRoutes.Use(auth.AuthMiddleware())
	commentRoutes.DELETE("/:id", deleteCommentHandler)
	commentRoutes.PUT("/:id/vote", voteCommentHandler)
	commentRoutes.DELETE("/:id/vote", removeVoteHandler)
}
