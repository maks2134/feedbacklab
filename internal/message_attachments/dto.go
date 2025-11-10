package message_attachments

type CreateMessageAttachmentDTO struct {
	ChatID     int     `json:"chat_id" validate:"required"`
	FilePath   string  `json:"file_path" validate:"required"`
	UploadedBy string  `json:"uploaded_by" validate:"required,uuid4"`
	FileType   *string `json:"file_type,omitempty"`
}

type UpdateMessageAttachmentDTO struct {
	FilePath string  `json:"file_path" validate:"omitempty"`
	FileType *string `json:"file_type,omitempty"`
}
