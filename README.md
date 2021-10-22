# Cronjob App for Backup
Little app to backup multiple databases and apps/folder/dir with zip archive written in Go.

# Features
* Backup multiple app or folder or dir.
* Backup multiple Postgresql database.
* Custom number of days to retain old backup otherwise will be deleted.
* Backup multiple app and database concurrently at once.


# Installation
1. Clone the repo.
```sh
$ git clone https://github.com/mdanialr/go-cron-backup.git
```

2. Create new config file. (assuming that you are in the root path of the repo)
```sh
$ cp config.yaml.example config.yaml
```

3. Fill in the config.yaml file as needed.

4. Build the project.
```sh
$ go build -o build/go-cron-backup main.go
```

5. Run a test to check if the app is working properly.
```sh
$ ./build/go-cron-backup -test -d
```
> If there is no error message in terminal then go to next step.

6. Run the app.
```sh
$ ./build/go-cron-backup
```

7. (optional) Create cronjob to run this app.
> Example
```sh
@daily cd /path/to/repo/go-cron-backup && ./build/go-cron-backup
```

# Notes and Suggestions
* Run the app with `sudo` privileges, since many app are reside in dir like `/var/www/*` need `sudo` privileges to do something with that dir, and this is also a mandatory to backup Postgresql database since in the app **'sudo -u postgres ...'** is used to do the trick.
* `-test` argument is used to test the app. (*only delete the zip files that created by this test*)
* `-d` argument is used to **delete all directories** in backup and log dir recursively including every files in that directories, so please be careful with this argument. (**_never use this argument when you have already run the app in production use, since this will delete all of your backuped files_**)
* See log file to check if there are some errors or successfull backup. (in **go-cron-backup--log**)