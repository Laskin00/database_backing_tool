# Zealroom backup tool
    Tool for database backup of zealroom project

## How to run
    golang  -  you need installed golang on the machine that you want to run the tool/backup the database
    build the executable with go build
    
## Available commands
     All commands need Database name, Database address, Database user and password

### Backup
    You can backup the data currently on the database by using:
     ./database_backup -a <database_adress> -u <database_user> -p <password> -d <database_name> -o <folder_containing_data>
    All the data (users,organizations,userOrganizationConnections) will be writen in separate files in the specified folder

### Recover
    You can recover the database from the previously backed-up data using:
    ./database_backup -r true -a <database_adress> -u <database_user> -p <password> -d <database_name> -o <folder_containing_data>

### Seed
    You can seed the database with the available data in ./seed using:
    ./database_backup -s true -a <database_adress> -u <database_user> -p <password> -d <database_name>
