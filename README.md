# Cronjob App for Backup
Little app to backup multiple databases and apps/folder/dir with zip archive written in Go.

# Features
* Backup multiple folder app or dir.
* Follow and walkthrough `symlink`.
* Backup multiple databases. supported __PostgreSQL__ & __MariaDB__.
* Custom max days to retain old backup before deleted.
* Backup multiple app and database concurrently at once.
* Pack the backup using ZIP archive (*deflate*).
* No root privileges is needed. (*as long as the user running this app has sufficient privileges*).
* Throttle CPU usage by setting up max worker.


# Installation
1. Clone the repo.
```sh
$ git clone https://github.com/mdanialr/go-cron-backup.git
```
> assuming that you are in the root path of the repo.
2. Get dependencies.
```sh
$ go mod tidy
```

3. Create new config file.
```sh
$ cp config.yaml.example config.yaml
```

4. Fill in the config.yaml file as needed.

5. Build the project.
```sh
$ go build -o build/go-cron-backup main.go
```

6. Run a test to check if the app is working properly.
```sh
$ ./build/go-cron-backup -test -d
```
> If there is no error message in terminal then go to next step.

7. Run the app.
```sh
$ ./build/go-cron-backup
```

8. (optional) Create a cronjob to run this app.
> Example
```sh
@daily cd /path/to/repo/go-cron-backup && ./build/go-cron-backup
```

# Arguments
* `-test` : test the app. (*will only delete the zip files that created by this test*).
* `-d` : **delete** backup and log folder recursively including every files in that directoris. so be careful with this argument. (**_never use this argument when you have already run the app in production, otherwise this will delete all of your backup files in that directory_**)
* `-no-app` : exclude app from this testing.
* `-no-db` : exclude db from this testing.
* `-sample` : the number of sample to be tested for both app and db.
* `-sam-app` : specifically set the number of sample for app.
* `-sam-db` : specifically set the number of sample for db.


# Notes
* `-sample`, `-sam-app` & `-sam-db` default value is 1 if not specified or overridden.
* Tested in linux. since this app only uses stdlib then this should also work with any other golang supported platform.
* See log file to check if there are some errors or successfull backup. (in **go-cron-backup--log**)
* Run this app when you are in the root path of the repo, otherwise you will see error regarding the config file is not found.

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.