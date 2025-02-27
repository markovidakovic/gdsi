package permission

type Permission string

const (
	// court permissions
	CreateCourt Permission = "create:court"
	UpdateCourt Permission = "update:court"
	DeleteCourt Permission = "delete:court"

	// season permissions
	CreateSeason Permission = "create:season"
	UpdateSeason Permission = "update:season"
	DeleteSeason Permission = "delete:season"

	// league permissions
	CreateLeague Permission = "create:league"
	UpdateLeague Permission = "update:league"
	DeleteLeague Permission = "delete:league"

	// match permissions
	CreateMatch Permission = "create:match"
	UpdateMatch Permission = "update:match"
	DeleteMatch Permission = "delete:match"
	SubmitScore Permission = "submit:score"

	// player permission
	UpdatePlayer Permission = "update:player"
	DeletePlayer Permission = "delete:player"
)

var rolePermissions = map[string][]Permission{
	"developer": {
		CreateCourt, UpdateCourt, DeleteCourt,
		CreateSeason, UpdateSeason, DeleteSeason,
		CreateLeague, UpdateLeague, DeleteLeague,
		CreateMatch, UpdateMatch, DeleteMatch, SubmitScore,
		UpdatePlayer, DeletePlayer,
	},
	"admin": {
		CreateCourt, UpdateCourt, DeleteCourt,
		CreateSeason, UpdateSeason, DeleteSeason,
		CreateLeague, UpdateLeague, DeleteLeague,
		CreateMatch, UpdateMatch, DeleteMatch, SubmitScore,
	},
	"user": {
		CreateMatch,
		SubmitScore,
	},
}

// Has accepts the account role and the needed permission to access the resource.
// It gets the permissions for a specific role and looks if the account role
// satisfies the required permission
func Has(role string, perm Permission) bool {
	perms, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range perms {
		if p == perm {
			return true
		}
	}

	return false
}
