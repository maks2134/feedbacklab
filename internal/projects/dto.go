package projects

type CreateProjectDTO struct {
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Description     *string `json:"description,omitempty"`
	GitlabProjectID *int    `json:"gitlab_project_id,omitempty"`
	MattermostTeam  *string `json:"mattermost_team,omitempty"`
}

type UpdateProjectDTO struct {
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Description     *string `json:"description,omitempty"`
	GitlabProjectID *int    `json:"gitlab_project_id,omitempty"`
	MattermostTeam  *string `json:"mattermost_team,omitempty"`
}
