package transport

type CreateDocumentationDTO struct {
	ProjectID  int     `json:"project_id" validate:"required"`
	FilePath   string  `json:"file_path" validate:"required"`
	Version    *string `json:"version,omitempty" validate:"omitempty"`
	UploadedBy *string `json:"uploaded_by,omitempty" validate:"omitempty,uuid4"`
}

type UpdateDocumentationDTO struct {
	FilePath   string  `json:"file_path" validate:"required"`
	Version    *string `json:"version,omitempty" validate:"omitempty"`
	UploadedBy *string `json:"uploaded_by,omitempty" validate:"omitempty,uuid4"`
}
