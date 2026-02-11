package comments

import (
	"net/http"
	"strings"
	"time"

	"restorapp/db"
	"restorapp/db/client"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func formatTimestamp(t pgtype.Timestamp) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

func getCommentsHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, err := uuid.Parse(productIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := ctx.GetString("userId")
	var comments []CommentResponse

	if userID != "" {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		rows, err := db.Queries.GetCommentsByProductIdWithUserVote(ctx, client.GetCommentsByProductIdWithUserVoteParams{
			ProductID: productUUID,
			UserID:    userUUID,
		})
		if err != nil {
			log.Error("Failed to get comments", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
			return
		}

		for _, row := range rows {
			var parentID *string
			if row.ParentID.Valid {
				s := uuid.UUID(row.ParentID.Bytes).String()
				parentID = &s
			}

			comments = append(comments, CommentResponse{
				ID:        row.ID.String(),
				ProductID: row.ProductID.String(),
				UserID:    row.UserID.String(),
				ParentID:  parentID,
				Content:   row.Content,
				CreatedAt: formatTimestamp(row.CreatedAt),
				UpdatedAt: formatTimestamp(row.UpdatedAt),
				Author: CommentAuthor{
					Name:  row.AuthorName,
					Image: row.AuthorImage.String,
				},
				VoteCounts: VoteCounts{
					Likes:    row.Likes,
					Dislikes: row.Dislikes,
				},
				UserVote: row.UserVote,
			})
		}
	} else {
		rows, err := db.Queries.GetCommentsByProductId(ctx, productUUID)
		if err != nil {
			log.Error("Failed to get comments", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
			return
		}

		for _, row := range rows {
			var parentID *string
			if row.ParentID.Valid {
				s := uuid.UUID(row.ParentID.Bytes).String()
				parentID = &s
			}

			comments = append(comments, CommentResponse{
				ID:        row.ID.String(),
				ProductID: row.ProductID.String(),
				UserID:    row.UserID.String(),
				ParentID:  parentID,
				Content:   row.Content,
				CreatedAt: formatTimestamp(row.CreatedAt),
				UpdatedAt: formatTimestamp(row.UpdatedAt),
				Author: CommentAuthor{
					Name:  row.AuthorName,
					Image: row.AuthorImage.String,
				},
				VoteCounts: VoteCounts{
					Likes:    row.Likes,
					Dislikes: row.Dislikes,
				},
				UserVote: "",
			})
		}
	}

	if comments == nil {
		comments = []CommentResponse{}
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
}

func createCommentHandler(ctx *gin.Context) {
	productIdParam := ctx.Param("id")
	productUUID, err := uuid.Parse(productIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	var parentID pgtype.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		parentUUID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent comment ID"})
			return
		}
		parentID = pgtype.UUID{Bytes: parentUUID, Valid: true}
	}

	comment, err := db.Queries.CreateComment(ctx, client.CreateCommentParams{
		ProductID: productUUID,
		UserID:    userUUID,
		ParentID:  parentID,
		Content:   content,
	})
	if err != nil {
		log.Error("Failed to create comment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Fetch the user to get author info
	user, err := db.Queries.GetUserById(ctx, userUUID)
	if err != nil {
		log.Error("Failed to get user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	var responseParentID *string
	if comment.ParentID.Valid {
		s := uuid.UUID(comment.ParentID.Bytes).String()
		responseParentID = &s
	}

	ctx.JSON(http.StatusCreated, CommentResponse{
		ID:        comment.ID.String(),
		ProductID: comment.ProductID.String(),
		UserID:    comment.UserID.String(),
		ParentID:  responseParentID,
		Content:   comment.Content,
		CreatedAt: formatTimestamp(comment.CreatedAt),
		UpdatedAt: formatTimestamp(comment.UpdatedAt),
		Author: CommentAuthor{
			Name:  user.Name,
			Image: user.Image.String,
		},
		VoteCounts: VoteCounts{
			Likes:    0,
			Dislikes: 0,
		},
		UserVote: "",
	})
}

func deleteCommentHandler(ctx *gin.Context) {
	commentIdParam := ctx.Param("id")
	commentUUID, err := uuid.Parse(commentIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = db.Queries.DeleteComment(ctx, client.DeleteCommentParams{
		ID:     commentUUID,
		UserID: userUUID,
	})
	if err != nil {
		log.Error("Failed to delete comment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func voteCommentHandler(ctx *gin.Context) {
	commentIdParam := ctx.Param("id")
	commentUUID, err := uuid.Parse(commentIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req VoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.VoteType != "like" && req.VoteType != "dislike" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Vote type must be 'like' or 'dislike'"})
		return
	}

	vote, err := db.Queries.UpsertCommentVote(ctx, client.UpsertCommentVoteParams{
		CommentID: commentUUID,
		UserID:    userUUID,
		VoteType:  req.VoteType,
	})
	if err != nil {
		log.Error("Failed to vote on comment", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"voteType": vote.VoteType})
}

func removeVoteHandler(ctx *gin.Context) {
	commentIdParam := ctx.Param("id")
	commentUUID, err := uuid.Parse(commentIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := ctx.GetString("userId")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = db.Queries.DeleteCommentVote(ctx, client.DeleteCommentVoteParams{
		CommentID: commentUUID,
		UserID:    userUUID,
	})
	if err != nil {
		log.Error("Failed to remove vote", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}
