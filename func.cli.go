package main

import (
	"path/filepath"
)

// FileModeApplication will run the File Mode application logic bootstrap
func (config CLIOpts) FileModeApplication() {
	logStdOut("Source file: " + config.Source)
	logStdOut("Target directory: " + config.Dir)

	absoluteTargetDirectory, err := filepath.Abs(config.Dir)
	check(err)

	// Initialize the Target Directory
	if err := ValidateConfigDirectory(absoluteTargetDirectory); err != nil {
		logStdErr("Target directory unwritable!")
	} else {
		directoryCheck, err := DirectoryExists(absoluteTargetDirectory)
		check(err)
		if !directoryCheck {
			//Directory doesn't exist, create
			CreateDirectory(absoluteTargetDirectory)
		}
	}
	CreateDirectory(absoluteTargetDirectory + "/config")
	CreateDirectory(absoluteTargetDirectory + "/zones")

	// Read in Zones file
	zones, err := NewZones(config)
	check(err)

	_, err = LoopThroughZonesForBindConfig(zones, absoluteTargetDirectory)
	check(err)

	_, err = LoopThroughZonesForBindZonesFiles(zones, absoluteTargetDirectory)
	check(err)

	_, err = LoopThroughZonesForBindReverseV4ZonesFiles(zones, absoluteTargetDirectory)
	check(err)
}
