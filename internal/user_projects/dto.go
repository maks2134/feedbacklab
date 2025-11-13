package user_projects

type CreateUserProjectDTO struct {
	UserID      string   `json:"user_id" validate:"required,uuid4"`
	ProjectID   int      `json:"project_id" validate:"required"`
	Role        string   `json:"role" validate:"required,oneof=viewer editor admin owner"`
	Permissions []string `json:"permissions" validate:"required,dive,alpha"`
}

type UpdateUserProjectDTO struct {
	Role        string   `json:"role" validate:"required,oneof=viewer editor admin owner"`
	Permissions []string `json:"permissions" validate:"required,dive,alpha"`
}
