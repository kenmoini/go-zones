package main

// Func main should be as small as possible and do as little as possible by convention
func main() {
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := ParseFlags()
	checkAndFail(err)

	// Run preflight
	PreflightSetup()

	if cfgPath.Mode == "server" {

		// Setup server config
		cfg, err := NewConfig(cfgPath)
		checkAndFail(err)

		// Run server preflight
		ServerPreflightSetup()

		// Run the server
		cfg.RunHTTPServer()

	} else {
		// Run file mode preflight setup
		FilePreflightSetup()

		// Run file mode application via CLI
		cfgPath.FileModeApplication()
	}

}
