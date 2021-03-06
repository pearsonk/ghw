//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package option

import "os"

const (
	defaultChroot           = "/"
	envKeyChroot            = "GHW_CHROOT"
	envKeySnapshotPath      = "GHW_SNAPSHOT_PATH"
	envKeySnapshotRoot      = "GHW_SNAPSHOT_ROOT"
	envKeySnapshotExclusive = "GHW_SNAPSHOT_EXCLUSIVE"
	envKeySnapshotPreserve  = "GHW_SNAPSHOT_PRESERVE"
)

// EnvOrDefaultChroot returns the value of the GHW_CHROOT environs variable or
// the default value of "/" if not set
func EnvOrDefaultChroot() string {
	// Grab options from the environs by default
	if val, exists := os.LookupEnv(envKeyChroot); exists {
		return val
	}
	return defaultChroot
}

// EnvOrDefaultSnapshotPath returns the value of the GHW_SNAPSHOT_PATH environs variable
// or the default value of "" (disable snapshot consumption) if not set
func EnvOrDefaultSnapshotPath() string {
	if val, exists := os.LookupEnv(envKeySnapshotPath); exists {
		return val
	}
	return "" // default is no snapshot
}

// EnvOrDefaultSnapshotRoot returns the value of the the GHW_SNAPSHOT_ROOT environs variable
// or the default value of "" (self-manage the snapshot unpack directory, if relevant) if not set
func EnvOrDefaultSnapshotRoot() string {
	if val, exists := os.LookupEnv(envKeySnapshotRoot); exists {
		return val
	}
	return "" // default is to self-manage the snapshot directory
}

// EnvOrDefaultSnapshotExclusive returns the value of the GHW_SNAPSHOT_EXCLUSIVE environs variable
// or the default value of false if not set
func EnvOrDefaultSnapshotExclusive() bool {
	if _, exists := os.LookupEnv(envKeySnapshotExclusive); exists {
		return true
	}
	return false
}

// EnvOrDefaultSnapshotPreserve returns the value of the GHW_SNAPSHOT_PRESERVE environs variable
// or the default value of false if not set
func EnvOrDefaultSnapshotPreserve() bool {
	if _, exists := os.LookupEnv(envKeySnapshotPreserve); exists {
		return true
	}
	return false
}

// Option is used to represent optionally-configured settings. Each field is a
// pointer to some concrete value so that we can tell when something has been
// set or left unset.
type Option struct {
	// To facilitate querying of sysfs filesystems that are bind-mounted to a
	// non-default root mountpoint, we allow users to set the GHW_CHROOT environ
	// vairable to an alternate mountpoint. For instance, assume that the user of
	// ghw is a Golang binary being executed from an application container that has
	// certain host filesystems bind-mounted into the container at /host. The user
	// would ensure the GHW_CHROOT environ variable is set to "/host" and ghw will
	// build its paths from that location instead of /
	Chroot *string

	// Snapshot contains options for handling ghw snapshots
	Snapshot *SnapshotOptions
}

// SnapshotOptions contains options for handling of ghw snapshots
type SnapshotOptions struct {
	// Path allows users to specify a snapshot (captured using ghw-snapshot) to be
	// automatically consumed. Users need to supply the path of the snapshot, and
	// ghw will take care of unpacking it on a temporary directory.
	// Set the environment variable "GHW_SNAPSHOT_PRESERVE" to make ghw skip the cleanup
	// stage and keep the unpacked snapshot in the temporary directory.
	Path string
	// Root is the directory on which the snapshot must be unpacked. This allows
	// the users to manage their snapshot directory instead of ghw doing that on
	// their behalf. Relevant only if SnapshotPath is given.
	Root *string
	// Exclusive tells ghw if the given directory should be considered of exclusive
	// usage of ghw or not If the user provides a Root. If the flag is set, ghw will
	// unpack the snapshot in the given SnapshotRoot iff the directory is empty; otherwise
	// any existing content will be left untouched and the unpack stage will exit silently.
	// As additional side effect, give both this option and SnapshotRoot to make each
	// context try to unpack the snapshot only once.
	Exclusive bool
}

func WithChroot(dir string) *Option {
	return &Option{Chroot: &dir}
}

// WithSnapshot sets snapshot-processing options for a ghw run
func WithSnapshot(opts SnapshotOptions) *Option {
	return &Option{
		Snapshot: &opts,
	}
}

// There is intentionally no Option related to GHW_SNAPSHOT_PRESERVE because we see that as
// a debug/troubleshoot aid more something users wants to do regularly.
// Hence we allow that only via the environment variable for the time being.

func Merge(opts ...*Option) *Option {
	merged := &Option{}
	for _, opt := range opts {
		if opt.Chroot != nil {
			merged.Chroot = opt.Chroot
		}
		if opt.Snapshot != nil {
			merged.Snapshot = opt.Snapshot
		}
	}
	// Set the default value if missing from mergeOpts
	if merged.Chroot == nil {
		chroot := EnvOrDefaultChroot()
		merged.Chroot = &chroot
	}
	if merged.Snapshot == nil {
		snapRoot := EnvOrDefaultSnapshotRoot()
		merged.Snapshot = &SnapshotOptions{
			Path:      EnvOrDefaultSnapshotPath(),
			Root:      &snapRoot,
			Exclusive: EnvOrDefaultSnapshotExclusive(),
		}
	}
	return merged
}
