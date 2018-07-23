package config

import "github.com/simonedegiacomi/gphotosuploader/auth"

func Authenticate(authFilePath string) *auth.CookieCredentials {
	// TODO: fix and uncomment

	// authOpts := auth.AuthenticationOptions{
	// 	AuthFilePath: authFilePath,
	// 	Silent:       true,
	// }

	// credentials := auth.Authenticate(authOpts)
	// return &credentials
	return &auth.CookieCredentials{}
}

// func loadCredentialsOrAuthenticate() *auth.CookieCredentials {
// 	// Load cookie for credentials from a json file
// 	credentials, err := auth.NewCookieCredentialsFromFile("auth.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Get a new API token using the TokenScraper from the api package
// 	token, err := api.NewAtTokenScraper(credentials).ScrapeNewAtToken()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Add the token to the credentials
// 	credentials.GetRuntimeParameters().AtToken = token
// 	return credentials
// }

// func initAuthentication(authFile string) auth.Credentials {
// 	// Load authentication parameters
// 	credentials, err := auth.NewCookieCredentialsFromFile(authFile)
// 	if err != nil {
// 		log.Printf("Can't use '%v' as auth file\n", authFile)
// 		credentials = nil
// 	} else {
// 		log.Println("Auth file loaded, checking validity ...")
// 		validity, err := credentials.TestCredentials()
// 		if err != nil {
// 			log.Fatalf("Can't check validity of credentials (%v)\n", err)
// 			credentials = nil
// 		} else if !validity.Valid {
// 			log.Printf("Credentials are not valid! %v\n", validity.Reason)
// 			credentials = nil
// 		} else {
// 			log.Println("Auth file seems to be valid")
// 		}
// 	}

// 	if credentials == nil {
// 		credentials, err = utils.StartWebDriverCookieCredentialsWizard()
// 		if err != nil {
// 			log.Fatalf("Can't complete the login wizard, got: %v\n", err)
// 		} else {
// 			// TODO: Handle error
// 			credentials.SerializeToFile(authFile)
// 		}

// 	}

// 	// Get a new At token
// 	token, err := api.NewAtTokenScraper(credentials).ScrapeNewAtToken()
// 	if err != nil {
// 		log.Fatalf("Can't scrape a new At token (%v)\n", err)
// 	}
// 	credentials.GetRuntimeParameters().AtToken = token

// 	return credentials
// }
