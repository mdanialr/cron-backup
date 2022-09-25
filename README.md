# Cronjob App for Backup
CLI for backup multiple databases and apps/folders/directories using zip archive format written in Go.

# Features
* Backup multiple folder app or dir.
* Follow and walk-through `symlink`.
* Backup multiple databases. support __PostgreSQL__, __MariaDB__ & __MySQL__ database.
* Custom max days to retain old backup before deleted.
* Backup multiple app and database concurrently at once.
* Pack the backup using ZIP archive (*deflate*).
* No root privileges is needed. (*as long as the user running this app has sufficient privileges which is __read-only__*).
* Throttle CPU & Memory usage by setting up max worker.
* ZIP archive without any external dependency (*thanks to Go's __stdlib__*).

# How to Use
1. Download the binary from [GitHub Releases](https://github.com/mdanialr/webhook/releases)
2. Create new config file with a filename `app.yml`. __The file name `app.yml`__ is mandatory otherwise
      [Viper](https://github.com/spf13/viper) will not find it
3. Extract then run to check if there is any error in config file
    ```bash
    tar -xzf cron-backup....tar.gz
    ./cron-backup
    ```
4. Create a cronjob to run this app (__optional__).
    > Example
    ```bash
    @daily cd /path/to/binary/file && ./cron-backup
    ```

## Example
Create config file with the filename `app.yml`
```yaml
root: /full/path/to/dir # required. root dir for backup, must be full path
max_days: # default to 6 days
db:
  max_days: # default to follow root max_days
  max_worker: # default to 1
  databases:
    - type: { pg/md/my } # pg for postgresql, md for mariadb, my for mysql
      host: ip_to_server # default to localhost
      port: db_port in integer # default to 5432, 3306 for pg and mdb or my respectively
      name: db_name # required. the database name
      user: db_usr # default to postgres for pg, root for md and my
      pass: db_pwd # default to postgres for pg, root for md and my
      backup: db-backup-name # required. unique. directory name for backup this database
      params: --opt --skip-lock-tables --single-transaction # can only be used for mariadb & mysql (md/my). will be ignored if the type is pg
    - type: { pg/md/my }
      docker: container_name # backup a database that's inside a docker container. the following host, port, name, user & pass config should be for the database inside the container
      host: ip_to_server # default to localhost. should point to database host within the container not where this app is run
      port: db_port in integer # default to 5432, 3306 for pg and mdb respectively. should point to database port within the container not where this app is run
      name: db_name # required. the database name inside the container
      user: db_usr # default to postgres for pg, root for md and my. the username that is used within the container, usually is root
      pass: db_pwd # default to postgres for pg, root for md and my. the password that is used within the container, usually is also root
      backup: db-another-name # required. unique. directory name for backup this database inside the container
app:
  max_days: # default to follow root max_days
  max_worker: 2 # default to 1
  apps:
    - dir: /full/path/to/app/dir # required. the app directory that will be archived
      name: some-app-name # required. unique. directory name for backup this app/directory
    - dir: /full/path/to/another/app/dir
      name: some-another-app-name
```

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.