package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/store/project"
)

func (s *store) CreateProject(domainID, name, description string, enabled bool) (*project.Project, error) {

	ok, err := s.projectNameExist(domainID, name)
	if err != nil {
		return nil, exception.NewInternalServerError("check project name exist error, %s", err)
	}
	if ok {
		return nil, exception.NewBadRequest("project name %s in this domain is exists", name)
	}

	proj := project.Project{ID: uuid.NewV4().String(), Name: name, Description: description, CreateAt: time.Now().Unix(), Enabled: enabled, DomainID: domainID}
	_, err = s.stmts[CreateProject].Exec(proj.ID, proj.Name, proj.Description, proj.Enabled, proj.DomainID, proj.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert project exec sql err, %s", err)
	}

	return &proj, nil
}

// Notice: if project not exits return nil
func (s *store) GetProject(id string) (*project.Project, error) {
	proj := project.Project{}
	err := s.stmts[FindProjectByID].QueryRow(id).Scan(
		&proj.ID, &proj.Name, &proj.Description, &proj.Enabled, &proj.CreateAt, &proj.DomainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("project %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single project error, %s", err)
	}

	return &proj, nil
}

func (s *store) ListDomainProjects(domainID string) ([]*project.Project, error) {
	rows, err := s.stmts[FindDomainPorjects].Query(domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project list error, %s", err)
	}
	defer rows.Close()

	projects := []*project.Project{}
	for rows.Next() {
		proj := project.Project{}
		if err := rows.Scan(&proj.ID, &proj.Name, &proj.Description, &proj.Enabled, &proj.DomainID, &proj.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project record error, %s", err)
		}
		projects = append(projects, &proj)
	}

	return projects, nil
}

func (s *store) UpdateProject(id, name, description string) (*project.Project, error) {
	return nil, nil
}

func (s *store) DeleteProject(id string) error {
	ret, err := s.stmts[DeleteProject].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete project exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("project %s not exist", id)
	}

	return nil
}

func (s *store) CheckProjectIsExistByID(id string) (bool, error) {
	var pid string
	err := s.stmts[CheckProjectExistByID].QueryRow(id).Scan(&pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("check project exist error, %s", err)
	}

	return true, nil
}

func (s *store) projectNameExist(domainID, projectName string) (bool, error) {
	rows, err := s.stmts[CheckProjectExistByName].Query(projectName, domainID)
	if err != nil {
		return false, fmt.Errorf("query project name exist error, %s", err)
	}
	defer rows.Close()

	projects := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, fmt.Errorf("scan project name exist record error, %s", err)
		}
		projects = append(projects, name)
	}
	if len(projects) != 0 {
		return true, nil
	}

	return false, nil
}

func (s *store) ListProjectUsers(projectID string) ([]string, error) {
	// check the project is exist
	ok, err := s.CheckProjectIsExistByID(projectID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("project %s not exist", projectID)
	}

	rows, err := s.stmts[FindProjectUsers].Query(projectID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project user id error, %s", err)
	}
	defer rows.Close()

	users := []string{}
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}
		users = append(users, userID)
	}
	return users, nil
}

func (s *store) AddUsersToProject(projectID string, userIDs ...string) error {
	// check the project is exist
	ok, err := s.CheckProjectIsExistByID(projectID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("project %s not exist", projectID)
	}

	// check user is in this project
	uids, err := s.ListProjectUsers(projectID)
	if err != nil {
		return err
	}
	existUIDs := []string{}
	for _, uid := range uids {
		for _, inuid := range userIDs {
			if inuid == uid {
				existUIDs = append(existUIDs, inuid)
			}
		}
	}
	if len(existUIDs) != 0 {
		return exception.NewBadRequest("users %s is in this project", existUIDs)
	}

	for _, userID := range userIDs {
		_, err = s.stmts[AddUsersToProject].Exec(userID, projectID)
		if err != nil {
			return fmt.Errorf("insert add users to project mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) RemoveUsersFromProject(projectID string, userIDs ...string) error {
	// check the project is exist
	ok, err := s.CheckProjectIsExistByID(projectID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("project %s not exist", projectID)
	}

	// check user is in this project
	uids, err := s.ListProjectUsers(projectID)
	if err != nil {
		return err
	}
	nExistUIDs := []string{}
	for _, inuid := range userIDs {
		var ok bool
		for _, uid := range uids {
			if uid == inuid {
				ok = true
			}
		}
		if !ok {
			nExistUIDs = append(nExistUIDs, inuid)
		}
	}
	if len(nExistUIDs) != 0 {
		return exception.NewBadRequest("users %s isn't in this project", nExistUIDs)
	}

	for _, userID := range userIDs {
		_, err = s.stmts[RemoveUsersFromProject].Exec(userID, projectID)
		if err != nil {
			return fmt.Errorf("insert remove users from project mapping exec sql err, %s", err)
		}
	}

	return nil
}
