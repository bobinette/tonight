package tonight_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/labstack/echo"
// 	"github.com/stretchr/testify/require"

// 	"github.com/bobinette/tonight"
// 	"github.com/bobinette/tonight/inmem"
// )

// func TestAPI(t *testing.T) {
// 	e := echo.New()
// 	projectJSON := `{"name": "test project"}`
// 	store := inmem.NewStore()
// 	require.NoError(t, tonight.RegisterHTTP(
// 		e.Group("/api"),
// 		store.EventStore(),
// 		store.TaskStore(),
// 		store.ProjectStore(),
// 		store.UserStore(),
// 	))

// 	ts := httptest.NewServer(e)
// 	defer ts.Close()
// 	client := http.Client{}

// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/projects", ts.URL), strings.NewReader(projectJSON))
// 	require.NoError(t, err)
// 	req.Header.Set("Token-Claim-Sub", "user")
// 	res, err := client.Do(req)
// 	require.NoError(t, err)
// 	require.Equal(t, 200, res.StatusCode)
// 	defer res.Body.Close()

// 	var resBody struct {
// 		Data tonight.Project
// 	}
// 	require.NoError(t, json.NewDecoder(res.Body).Decode(&resBody))
// 	res.Body.Close()
// 	project := resBody.Data

// 	taskJSON := fmt.Sprintf(`{"title": "task", "project": {"uuid": "%s"}}`, project.UUID)
// 	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/tasks", ts.URL), strings.NewReader(taskJSON))
// 	require.NoError(t, err)
// 	req.Header.Set("Token-Claim-Sub", "user")
// 	res, err = client.Do(req)
// 	require.NoError(t, err)
// 	require.Equal(t, 200, res.StatusCode)
// 	defer res.Body.Close()

// 	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/projects", ts.URL), strings.NewReader(projectJSON))
// 	require.NoError(t, err)
// 	req.Header.Set("Token-Claim-Sub", "user")
// 	res, err = client.Do(req)
// 	require.NoError(t, err)
// 	require.Equal(t, 200, res.StatusCode)
// 	defer res.Body.Close()

// 	var response struct {
// 		Projects []tonight.Project `json:"data"`
// 	}
// 	require.NoError(t, json.NewDecoder(res.Body).Decode(&response))
// 	require.Len(t, response.Projects, 1)
// 	require.Equal(t, "test project", response.Projects[0].Name)
// 	require.Len(t, response.Projects[0].Tasks, 1)
// 	require.Equal(t, "task", response.Projects[0].Tasks[0].Title)

// 	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/projects", ts.URL), strings.NewReader(projectJSON))
// 	require.NoError(t, err)
// 	req.Header.Set("Token-Claim-Sub", "other_user")
// 	res, err = client.Do(req)
// 	require.NoError(t, err)
// 	require.Equal(t, 200, res.StatusCode)
// 	defer res.Body.Close()

// 	response.Projects = nil
// 	require.NoError(t, json.NewDecoder(res.Body).Decode(&response))
// 	require.Len(t, response.Projects, 0)
// }
