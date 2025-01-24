package permissions

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
)

var rolePermissions = map[string][]Permission{
	"developer": {
		CreateCourt, UpdateCourt, DeleteCourt,
		CreateSeason, UpdateSeason, DeleteSeason,
		CreateLeague, UpdateLeague, DeleteLeague,
		CreateMatch, UpdateMatch, DeleteMatch, SubmitScore,
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

func HasPermission(role string, perm Permission) bool {
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
