package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/xiaolonglong/harborctl/internal/config"
)

// Client represents the Harbor API client
type Client struct {
	Config     *config.Config
	HTTPClient *http.Client
	baseURL   string
	auth     string
}

// NewClient creates a new Harbor API client
func NewClient(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Insecure},
	}

	return &Client{
		Config:     cfg,
		HTTPClient: &http.Client{Timeout: 30 * time.Second, Transport: tr},
		baseURL:   cfg.GetBaseURL(),
		auth:     fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(cfg.Username+":"+cfg.Password))),
	}, nil
}

// Project represents a Harbor project
type Project struct {
	ProjectID    int       `json:"project_id"`
	Name        string    `json:"name"`
	Public      bool      `json:"public"`
	OwnerID     int       `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// ProjectReq represents the request body for creating/updating a project
type ProjectReq struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

// User represents a Harbor user
type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RealName string `json:"realname"`
	Password string `json:"-"`
	Admin    bool   `json:"admin"`
}

// Repository represents a Harbor repository
type Repository struct {
	Name        string    `json:"name"`
	ProjectID   int       `json:"project_id"`
	Description string    `json:"description"`
	PullCount  int       `json:"pull_count"`
	StarCount  int       `json:"star_count"`
	Size       int64     `json:"size"`
	Tags       []string  `json:"tags,omitempty"`
	UpdateTime time.Time `json:"update_time"`
}

// Artifact represents a Harbor artifact
type Artifact struct {
	Digest           string        `json:"digest"`
	Type            string        `json:"type"`
	Tags            interface{}   `json:"tags"`
	PushTime        time.Time     `json:"push_time"`
	PullTime        time.Time     `json:"pull_time"`
	Size            int64        `json:"size"`
	Architecture   string       `json:"architecture"`
	OS             string       `json:"os"`
	ManifestConfig *ManifestConfig `json:"manifest_config"`
	Layers         []Layer       `json:"layers"`
}

// ManifestConfig represents manifest configuration
type ManifestConfig struct {
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}

// Layer represents a container layer
type Layer struct {
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}

// SystemInfo represents Harbor system information
type SystemInfo struct {
	HarborVersion    string `json:"harbor_version"`
	DatabaseType    string `json:"database_type"`
	SelfRegistration bool `json:"self_registration"`
	LDAPEnabled    bool   `json:"ldap_enabled"`
	Scanner        string `json:"scanner"`
}

// OverallHealthStatus represents health status
type OverallHealthStatus struct {
	Status     string           `json:"status"`
	Components []ComponentItem `json:"components"`
	Harbor     *ComponentHealth `json:"harbor,omitempty"`
	Portal     *ComponentHealth `json:"portal,omitempty"`
	Core       *ComponentHealth `json:"core,omitempty"`
	Jobservice *ComponentHealth `json:"jobservice,omitempty"`
	Registry  *ComponentHealth `json:"registry,omitempty"`
	Database  *ComponentHealth `json:"database,omitempty"`
	Redis     *ComponentHealth `json:"redis,omitempty"`
	Proxy     *ComponentHealth `json:"proxy,omitempty"`
}

// ComponentItem represents a component in health check
type ComponentItem struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ComponentHealth represents health of a component
type ComponentHealth struct {
	Status string `json:"status"`
}

// Statistic represents statistics
type Statistic struct {
	ProjectCount int `json:"project_count"`
	RepoCount   int `json:"repo_count"`
}

// ErrorResponse represents error
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ProjectMember represents project member
type ProjectMember struct {
	MemberID   int    `json:"member_id"`
	EntityName string `json:"entity_name"`
	EntityID  int    `json:"entity_id"`
	RoleName  string `json:"role_name"`
	RoleID    int    `json:"role_id"`
}

// SearchResult represents search results
type SearchResult struct {
	Projects    []Project    `json:"project"`
	Repositories []Repository `json:"repository"`
}

// ListProjects lists all projects
func (c *Client) ListProjects(name string, page, pageSize int) ([]Project, error) {
	params := url.Values{}
	if name != "" {
		params.Add("name", name)
	}
	if page > 0 {
		params.Add("page", fmt.Sprintf("%d", page))
	}
	if pageSize > 0 {
		params.Add("page_size", fmt.Sprintf("%d", pageSize))
	}
	params.Add("with_detail", "true")

	uri := c.baseURL + "/projects"
	if len(params) > 0 {
		uri += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject gets a project by name or ID
func (c *Client) GetProject(projectNameOrID string) (*Project, error) {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(projectName string, isPublic bool) (*Project, error) {
	body := ProjectReq{
		Name:   projectName,
		Public: isPublic,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	uri := c.baseURL + "/projects"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectNameOrID string, force bool) error {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID)
	if force {
		uri += "?force=true"
	}

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetProjectDeletable checks if deletable
func (c *Client) GetProjectDeletable(projectNameOrID string) (bool, string, error) {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/_deletable"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	var deletable struct {
		Deletable bool   `json:"deletable"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&deletable); err != nil {
		return false, "", err
	}

	return deletable.Deletable, deletable.Message, nil
}

// ListRepositories lists repositories
func (c *Client) ListRepositories(projectNameOrID string) ([]Repository, error) {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/repositories"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// ListArtifacts lists artifacts
func (c *Client) ListArtifacts(projectNameOrID, repositoryName string) ([]Artifact, error) {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/repositories/" + url.PathEscape(repositoryName) + "/artifacts"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artifacts []Artifact
	if err := json.NewDecoder(resp.Body).Decode(&artifacts); err != nil {
		return nil, err
	}

	return artifacts, nil
}

// DeleteArtifact deletes an artifact
func (c *Client) DeleteArtifact(projectNameOrID, repositoryName, digest string) error {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/repositories/" + url.PathEscape(repositoryName) + "/artifacts/" + url.PathEscape(digest)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GCJob represents a garbage collection job
type GCJob struct {
	JobID       int       `json:"id"`
	JobKind     string    `json:"job_kind"`
	JobStatus  string    `json:"job_status"`
	JobDetail  string    `json:"job_detail,omitempty"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime time.Time `json:"update_time"`
	Schedule   *GCSchedule `json:"schedule,omitempty"`
}

// GCSchedule represents GC schedule
type GCSchedule struct {
	ScheduleType string `json:"schedule_type"`
	Cron        string `json:"cron"`
	Updatetime  int64    `json:"updatetime,omitempty"`
	Parameters *GCParameters `json:"parameters,omitempty"`
}

// GCParameters represents GC parameters
type GCParameters struct {
	DeleteUntagged bool `json:"delete_untagged"`
	AgeDays       int    `json:"age_days,omitempty"`
	DryRun       bool   `json:"dry_run,omitempty"`
}

// GetGCJobs gets GC jobs
func (c *Client) GetGCJobs() ([]GCJob, error) {
	uri := c.baseURL + "/gc/executions"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs []GCJob
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

// TriggerGC triggers a garbage collection
func (c *Client) TriggerGC(dryRun bool) (*GCJob, error) {
	body := map[string]bool{
		"dry_run": dryRun,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	uri := c.baseURL + "/gc/executions"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var job GCJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

// GetGCSchedule gets GC schedule
func (c *Client) GetGCSchedule() (*GCSchedule, error) {
	uri := c.baseURL + "/gc/schedule"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()


	var schedule GCSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, err
	}

	return &schedule, nil
}

// UpdateGCSchedule updates GC schedule
func (c *Client) UpdateGCSchedule(scheduleType, cron string, deleteUntagged bool) error {
	body := map[string]interface{}{
		"schedule": map[string]interface{}{
			"schedule_type": scheduleType,
			"cron":       cron,
			"parameters": map[string]bool{
				"delete_untagged": deleteUntagged,
			},
		},
	}


	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	uri := c.baseURL + "/gc/schedule"
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetSystemInfo gets system info
func (c *Client) GetSystemInfo() (*SystemInfo, error) {
	uri := c.baseURL + "/systeminfo"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info SystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

// GetHealth gets health status
func (c *Client) GetHealth() (*OverallHealthStatus, error) {
	uri := c.baseURL + "/health"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var health OverallHealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}

// GetStatistics gets statistics
func (c *Client) GetStatistics() (*Statistic, error) {
	uri := c.baseURL + "/statistics"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stat Statistic
	if err := json.NewDecoder(resp.Body).Decode(&stat); err != nil {
		return nil, err
	}

	return &stat, nil
}

// ListUsers lists users
func (c *Client) ListUsers() ([]User, error) {
	uri := c.baseURL + "/users"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser gets user by ID
func (c *Client) GetUser(userID int) (*User, error) {
	uri := c.baseURL + "/users/" + fmt.Sprintf("%d", userID)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a user
func (c *Client) CreateUser(user *User) (*User, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	uri := c.baseURL + "/users"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newUser User
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(userID int) error {
	uri := c.baseURL + "/users/" + fmt.Sprintf("%d", userID)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateUser updates a user
func (c *Client) UpdateUser(user *User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	uri := c.baseURL + "/users/" + fmt.Sprintf("%d", user.UserID)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListProjectMembers lists project members
func (c *Client) ListProjectMembers(projectNameOrID string) ([]ProjectMember, error) {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/members"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var members []ProjectMember
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}

	return members, nil
}

// AddProjectMember adds a member
func (c *Client) AddProjectMember(projectNameOrID string, member *ProjectMember) error {
	body := map[string]interface{}{
		"member_id": member.EntityID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/members"
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RemoveProjectMember removes a member
func (c *Client) RemoveProjectMember(projectNameOrID string, memberID int) error {
	uri := c.baseURL + "/projects/" + url.PathEscape(projectNameOrID) + "/members/" + fmt.Sprintf("%d", memberID)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Search searches
func (c *Client) Search(query string) (*SearchResult, error) {
	uri := c.baseURL + "/search?q=" + url.QueryEscape(query)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.auth)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// doRequest handles HTTP requests
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)

		switch resp.StatusCode {
		case 401:
			return nil, fmt.Errorf("authentication failed: check credentials")
		case 403:
			return nil, fmt.Errorf("forbidden: permission denied")
		case 404:
			return nil, fmt.Errorf("not found: %s", errResp.Message)
		case 409:
			return nil, fmt.Errorf("conflict: %s", errResp.Message)
		case 412:
			return nil, fmt.Errorf("precondition failed: %s", errResp.Message)
		default:
			if errResp.Message != "" {
				return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Message)
			}
			return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
		}
	}

	return resp, nil
}
