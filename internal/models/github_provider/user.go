package githubprovider

import "time"

type GithubUser struct {
	AvatarURL               string      `json:"avatar_url"`
	Bio                     *string     `json:"bio"`
	Blog                    *string     `json:"blog"`
	BusinessPlus            *bool       `json:"business_plus,omitempty"`
	Collaborators           *int64      `json:"collaborators,omitempty"`
	Company                 *string     `json:"company"`
	CreatedAt               time.Time   `json:"created_at"`
	DiskUsage               *int64      `json:"disk_usage,omitempty"`
	Email                   *string     `json:"email"`
	EventsURL               string      `json:"events_url"`
	Followers               int64       `json:"followers"`
	FollowersURL            string      `json:"followers_url"`
	Following               int64       `json:"following"`
	FollowingURL            string      `json:"following_url"`
	GistsURL                string      `json:"gists_url"`
	GravatarID              *string     `json:"gravatar_id"`
	Hireable                *bool       `json:"hireable"`
	HTMLURL                 string      `json:"html_url"`
	ID                      int64       `json:"id"`
	LDAPDN                  *string     `json:"ldap_dn,omitempty"`
	Location                *string     `json:"location"`
	Login                   string      `json:"login"`
	Name                    *string     `json:"name"`
	NodeID                  string      `json:"node_id"`
	NotificationEmail       *string     `json:"notification_email"`
	OrganizationsURL        string      `json:"organizations_url"`
	OwnedPrivateRepos       *int64      `json:"owned_private_repos,omitempty"`
	Plan                    *GithubPlan `json:"plan,omitempty"`
	PrivateGists            *int64      `json:"private_gists,omitempty"`
	PublicGists             int64       `json:"public_gists"`
	PublicRepos             int64       `json:"public_repos"`
	ReceivedEventsURL       string      `json:"received_events_url"`
	ReposURL                string      `json:"repos_url"`
	SiteAdmin               bool        `json:"site_admin"`
	StarredURL              string      `json:"starred_url"`
	SubscriptionsURL        string      `json:"subscriptions_url"`
	TotalPrivateRepos       *int64      `json:"total_private_repos,omitempty"`
	TwitterUsername         *string     `json:"twitter_username"`
	TwoFactorAuthentication *bool       `json:"two_factor_authentication,omitempty"`
	Type                    string      `json:"type"`
	UpdatedAt               time.Time   `json:"updated_at"`
	URL                     string      `json:"url"`
	UserViewType            *string     `json:"user_view_type,omitempty"`
}

type GithubPlan struct {
	Collaborators int64  `json:"collaborators"`
	Name          string `json:"name"`
	PrivateRepos  int64  `json:"private_repos"`
	Space         int64  `json:"space"`
}
