// Package constants defines shared constants and enumerations for the Labra
// backend. It lives in src/api/constants.
package constants

import (
	"time"

	"github.com/MakeNowJust/heredoc"
)

// Status represents an application or resource status.
type Status string

// Role defines a user role name used for authorization.
type Role string

// Environment indicates the deployment environment such as dev or prod.
type Environment string

// ContextKey identifies values stored in context.Context.
type ContextKey string

// Commonly used boolean and string constants.
var (
	FalseVal                = false
	TrueVal                 = true
	Now                     = time.Now()
	SystemUsername          = "system"
	DefaultSystemKey        = "defaultSystemKey"
	GlobalDepartment        = "Global"
	RootConfig              = "root"
	DefaultRemoteFolderPath = "files"
	PrivateACL              = "private"
	AwsS3Domain             = "https://s3.amazonaws.com/"
	RpcPluginGoModContent   = heredoc.Docf(`
		module graphabc-impl
		go 1.18
		replace graphabc => ../../../../../
	`)
)

// Enumerations and configuration constants used throughout the API.
const (
	Active     Status = "Active"
	Inactive   Status = "Inactive"
	Stopped    Status = "Stopped"
	Running    Status = "Running"
	Suspended  Status = "Suspended"
	Blocked    Status = "Blocked"
	New        Status = "New"
	Retired    Status = "Retired"
	Processing Status = "Processing"
	Ready      Status = "Ready"

	SuperAdmin Role = "SuperAdmin"

	MigrationsDir  = "pkg/infrastructure/entgql/ent/migrate/migrations"
	MigrateDataDir = "pkg/infrastructure/entgql/ent/migrate/migratedata"

	CurrentUserContextValue    = "currentUser"
	DefaultPluginFolderPath    = "goplugin"
	DefaultRpcPluginFolderPath = "rpc"

	RepositoryContextValue       ContextKey = "repository"
	ServiceContextValue          ContextKey = "service"
	AdminRepositoryContextValue  ContextKey = "adminRepository"
	AdminServiceContextValue     ContextKey = "adminService"
	CentrifugeClientContextValue ContextKey = "centrifugeClient"
	UserContextValue             ContextKey = "user"
	RoleContextValue             ContextKey = "role"

	Local Environment = "local"
	Dev   Environment = "dev"
	Prod  Environment = "prod"
)

// String returns the textual representation of the Status value.
func (s Status) String() string {
	switch s {
	case Active:
		return "Active"
	case Inactive:
		return "Inactive"
	case Stopped:
		return "Stopped"
	case Running:
		return "Running"
	case Suspended:
		return "Suspended"
	case Blocked:
		return "Blocked"
	case New:
		return "New"
	case Retired:
		return "Retired"
	case Processing:
		return "Processing"
	case Ready:
		return "Ready"
	default:
		return "Unknown"
	}
}

// String returns the role name.
func (r Role) String() string {
	switch r {
	case SuperAdmin:
		return "SuperAdmin"
	default:
		return "Unknown"
	}
}

// String converts the Environment value to its string form.
func (env Environment) String() string {
	switch env {
	case Local:
		return "local"
	case Dev:
		return "dev"
	case Prod:
		return "prod"
	default:
		return "Unknown"
	}
}
