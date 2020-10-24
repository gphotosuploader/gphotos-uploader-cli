# Getting started

## Install
You can install the pre-compiled binary (in several different ways) or compile from source.

Here are the steps for each of them:

### Install the pre-compiled binary

**homebrew tap** (only on macOS for now):
```
$ brew install gphotosuploader/tap/gphotos-uploader-cli
```

**manually**

Download the pre-compiled binaries from the [releases page](https://github.com/gphotosuploader/gphotos-uploader-cli/releases/latest) and copy to the desired location.

### Compiling from source

> This project will maintain compatibility with the last two Go major versions published. It could work with other versions but we can't support it. 

You can compile the source code in your system.

```
$ git clone https://github.com/gphotosuploader/gphotos-uploader-cli
$ cd gphotos-uploader-cli
$ make build
```

Or you can use `go get` if you prefer it:

```
$ go get github.com/gphotosuploader/gphotos-uploader-cli
```

## Configure
First initialize the config file using this command:
```
$ gphotos-uploader-cli init
```

> Default configuration folder is `~/.gphotos-uploader-cli` but you can specify your own folder using `--config /my/config/dir`. Configuration is kept in the `config.hjson` file inside this folder.

You must review the [documentation](configuration.md) to specify your **Google Photos API credentials**, `APIAppCredentials`. You should tune your `jobs` configuration also.

## Run
Once it's configured you can start uploading files in this way:
``` 
$ gphotos-uploader-cli push
```

### First time run
The first time you run `gphotos-uploader-cli`, after setting your configuration ([Google Photos API credentials](configuration.md#APIAppCredentials)), few manual steps are needed:

1. You should get an output like this one:

```
Visit the following URL in your browser:
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=...

After completing the authorization flow, enter the authorization code here:
```

2. Open a browser and point to the previous URL. Select the account where you wan to upload your files (the same you configured in the config file). You will see something like this:

![Google asking for Google Photos API credentials](images/ask_Google_Photos_API_credentials.png) 

3. After that, you should confirm that you trust on `gphotos-uploader-cli` to access to your Google Photos account, click on **Go to gphotos-uploader**:

![Google ask you to verify gphotos-upload-cli](images/ask_for_application_verification.png)

4. Finally Google will ask you to confirm permission Google Photos account:

![Google ask permission to your Google Photos account](images/ask_for_permission.png)

5. A page with a code is shown in your browser, copy this code and go back to the terminal.

![Final confirmation, all was good](images/final_confirmation.png)

6. Paste the previous code in your terminal to complete the process.

```
After completing the authorization flow, enter the authorization code here: 4/4QFPtCv11dN3a-hVYhHkMryZe5g
```

All auth configuration is in place.