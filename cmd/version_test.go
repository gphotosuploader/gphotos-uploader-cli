package cmd

func ExampleVersionCmd_Run() {
	cmd := &VersionCmd{}
	_ = cmd.Run(nil, nil)
	// Output: gphotos-uploader-cli v0.0.0
}
