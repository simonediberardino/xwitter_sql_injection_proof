package dto

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"-"`
	DisplayName  string `json:"displayName"`
	Bio          string `json:"bio"`
	PasswordHash string `json:"-"`
	CreatedAt    string `json:"createdAt"`
	PostCount    string `json:"postCount"`
}

type Author struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"-"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	Author    Author `json:"author"`
}

type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Bio         string `json:"bio"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreatePostRequest struct {
	Content string `json:"content"`
}

type SearchQuery struct {
	Query string `json:"q"`
}

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Bio         string `json:"bio"`
}

type PostFilter struct {
	UserID string `json:"userId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
