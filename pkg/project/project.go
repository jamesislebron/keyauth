package project

// Project tenant resource container
type Project struct {
	// the project's domain id
	DomainID string
	// Project id, UUID as a unique logo
	ID string
	// Project name, allow repeat
	Name string
	// Project description
	Description string
	// Whether to enable
	Enabled bool
	// CrateAt create time
	CreateAt int64
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Manager is project service
type Manager interface {
	// Create a Project, super admin & domain admin are
	// allowed to operate, Named in Domain, does not allow renaming
	CreateProject(domainID, name, description string, enabled bool) (*Project, error)
	// Get a project with project id
	GetProject(id string) (*Project, error)
	// List all Project in domain_id, else all project
	ListProject(domainID string) ([]*Project, error)
	// Update a Project, super admin & domain admin are allowed to operate
	UpdateProject(id, name, description string) (*Project, error)
	// Soft Delete a Project,project still in persistence storage, super admin & domain admin are allowed to operate
	DeleteProject(id string) error
	// IsExist use to check the project is exist by project id
	IsExist(id string) (bool, error)
}
