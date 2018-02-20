# gitlab-api-extension
Extends Gitlab API without having to fork Gitlab, by proxying requests to Gitlab

Added methods : 

		"/api/v4/time_logs"
		"/api/v4/projects/{projectID}/issues/{issueIID}/time_logs"
		"/api/v4/projects/{projectID}/merge_requests/{mergeRequestIID}/time_logs"
		"/api/v4/users/{userID}/time_logs"
		"/api/v4/projects/{projectID}/time_logs"
		"/api/v4/projects/{projectID}/users/{userID}/time_logs"
		"/api/v4/projects/{projectID}/issues/{issueIID}/users/{userID}/time_logs"

Authors:

  - Antoine Huret (@antony360)
  - Jean-Baptiste Guerraz (@jbguerraz)
