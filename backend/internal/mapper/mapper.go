package mapper

import (
	"strconv"
	"time"

	"xwitter/internal/db"
	"xwitter/internal/dto"
)

func UserToDTO(user db.User) dto.User {
	return dto.User{
		ID:           strconv.Itoa(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Bio:          user.Bio,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.UTC().Format(time.RFC3339),
		PostCount:    "0",
	}
}

func UsersToDTO(users []db.User) []dto.User {
	out := make([]dto.User, 0, len(users))
	for _, user := range users {
		out = append(out, UserToDTO(user))
	}
	return out
}

func UserToDB(user dto.User) (db.User, error) {
	out := db.User{
		Username:     user.Username,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Bio:          user.Bio,
		PasswordHash: user.PasswordHash,
	}

	if user.ID != "" {
		id, err := strconv.Atoi(user.ID)
		if err != nil {
			return db.User{}, err
		}
		out.ID = id
	}

	if user.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
		if err != nil {
			return db.User{}, err
		}
		out.CreatedAt = createdAt
	}

	return out, nil
}

func AuthorToDTO(user db.User) dto.Author {
	return dto.Author{
		ID:          strconv.Itoa(user.ID),
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}
}

func PostToDTO(post db.Post, author db.User) dto.Post {
	return dto.Post{
		ID:        strconv.Itoa(post.ID),
		UserID:    strconv.Itoa(post.UserID),
		Content:   post.Content,
		CreatedAt: post.CreatedAt.UTC().Format(time.RFC3339),
		Author:    AuthorToDTO(author),
	}
}

func PostsToDTO(posts []db.Post, authors []db.User) []dto.Post {
	out := make([]dto.Post, 0, len(posts))
	for i, post := range posts {
		out = append(out, PostToDTO(post, authors[i]))
	}
	return out
}

func PostToDB(post dto.Post) (db.Post, error) {
	out := db.Post{Content: post.Content}

	if post.ID != "" {
		id, err := strconv.Atoi(post.ID)
		if err != nil {
			return db.Post{}, err
		}
		out.ID = id
	}

	if post.UserID != "" {
		userID, err := strconv.Atoi(post.UserID)
		if err != nil {
			return db.Post{}, err
		}
		out.UserID = userID
	}

	if post.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			return db.Post{}, err
		}
		out.CreatedAt = createdAt
	}

	return out, nil
}
