package transport

type CreateModuleDTO struct {
	ProjectID         int     `json:"project_id" validate:"required"`
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	Description       *string `json:"description,omitempty"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdateModuleDTO struct {
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	Description       *string `json:"description,omitempty"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty" validate:"omitempty,uuid4"`
}
