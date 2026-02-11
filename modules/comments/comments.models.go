package comments

type CommentAuthor struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type VoteCounts struct {
	Likes    int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
}

type CommentResponse struct {
	ID         string        `json:"id"`
	ProductID  string        `json:"productId"`
	UserID     string        `json:"userId"`
	ParentID   *string       `json:"parentId"`
	Content    string        `json:"content"`
	CreatedAt  string        `json:"createdAt"`
	UpdatedAt  string        `json:"updatedAt"`
	Author     CommentAuthor `json:"author"`
	VoteCounts VoteCounts    `json:"voteCounts"`
	UserVote   string        `json:"userVote"`
}

type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parentId"`
}

type VoteRequest struct {
	VoteType string `json:"voteType"`
}
