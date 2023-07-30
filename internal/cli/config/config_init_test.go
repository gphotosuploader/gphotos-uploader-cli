package config_test

//func TestNewInitCmd(t *testing.T) {
//	testCases := []struct {
//		name          string
//		input         string
//		args          []string
//		isErrExpected bool
//	}{
//		{"Should success", "", []string{"config", "init", "--overwrite"}, false},
//		{"Should fail if input exists", "/foo", []string{"config", "init"}, true},
//		{"Should success if input exists and force is set", "/foo", []string{"config", "init", "--overwrite"}, false},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			cli.Os = afero.NewMemMapFs()
//			//createTestConfigurationFile(t, cli.Os, tc.input)
//
//			actual := new(bytes.Buffer)
//			configuration.Settings = configuration.Init("")
//			configuration.Settings.SetFs(afero.NewMemMapFs())
//			rootCommand := cli.New()
//			//rootCommand.SetOut(actual)
//			//rootCommand.SetErr(actual)
//			rootCommand.SetArgs(tc.args)
//
//			err := rootCommand.Execute()
//			if tc.isErrExpected {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//
//			assert.Equal(t, "hola", actual.String())
//		})
//	}
//}

//func createTestConfigurationFile(t *testing.T, fs afero.Fs, path string) {
//	if path == "" {
//		return
//	}
//	if err := fs.MkdirAll(path, 0700); err != nil {
//		t.Fatalf("creating test dir, err: %s", err)
//	}
//	filename := filepath.Join(path, app.DefaultConfigFilename)
//	if err := afero.WriteFile(fs, filename, []byte("my"), 0600); err != nil {
//		t.Fatalf("creating test configuration file, err: %s", err)
//	}
//}
