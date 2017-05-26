Minimal HTTP server to store and retrieve files from a folder. Compiled for Linux AMD64, but you can re-compile to any Go supported architecture.

#### Use cases:
 - sharing files via HTTP
 - storing files for long term (like artifact repository)

## Usage examples

#### Running application

Run storage and expose all files in `/home/user/files` to `localhost:8223`:

    ./storage --path="/home/user/files" --port=8223

To keep server up you should run it with some process manager, for example, [supervisor]([http://supervisord.org/).

#### Uploading files

Then, you are able to upload file `file.tar.gz` via HTTP post using [curl]([https://curl.haxx.se/):

    curl -X POST http://localhost:8223/newFolder -F anythingHere=@file.tar.gz

Multiple file uploads are supported too:

    curl -X POST http://localhost:8223 -F anythingHere=@file.tar.gz -F somethingElse=@file2.tar.gz

#### Retrieving files

Uploaded files should be available at `http://localhost:8223/{optionalFolder}/{fileName}`

#### Removing files

    curl -X DELETE http://localhost:8223/{fileName}

### Security considerations

There are no limits who can create and delete files, to secure this server, use proxy with access authentication and specific firewall rules to avoid unwanted traffic.

### Limitations
 - Files with same name are overwritten without any warnings
 - HTTPS is not supported without additional proxy
 - It is not possible to mass delete files, only one by one or whole directory

### Command line options

    Usage of ./storage:
        -debug
        	Enable debug output
        -max int
             Max uploaded file size in bytes (default 65536)
        -path string
             Storage path for files (default "/tmp")
        -port int
             Port for connecting to application (default 8222)
