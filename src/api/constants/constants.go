package constants

import (
	"time"

	"github.com/MakeNowJust/heredoc"
)

type Status string
type Role string
type Environment string
type ContextKey string

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
	CentrifugeClientContextValue ContextKey = "centrifugeClient"
	UserContextValue             ContextKey = "user"
	RoleContextValue             ContextKey = "role"

	Local Environment = "local"
	Dev   Environment = "dev"
	Prod  Environment = "prod"
)

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

func (r Role) String() string {
	switch r {
	case SuperAdmin:
		return "SuperAdmin"
	default:
		return "Unknown"
	}
}

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
