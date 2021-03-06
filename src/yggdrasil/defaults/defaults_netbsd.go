// +build netbsd

package defaults

// Sane defaults for the BSD platforms. The "default" options may be
// may be replaced by the running configuration.
func GetDefaults() platformDefaultParameters {
	return platformDefaultParameters{
		// Admin
		DefaultAdminListen: "tcp://localhost:9001",

		// TUN/TAP
		MaximumIfMTU:     9000,
		DefaultIfMTU:     9000,
		DefaultIfName:    "/dev/tap0",
		DefaultIfTAPMode: true,
	}
}
