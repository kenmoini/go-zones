package main

// FileModeApplication will run the File Mode application logic bootstrap
func (config CLIOpts) FileModeApplication() {
	logStdOut("Source file: " + config.Source)
	logStdOut("Target directory: " + config.Dir)
}
