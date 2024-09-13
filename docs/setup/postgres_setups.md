#### Install PostGIS on PostgreSQL 16 for GeoDjango Setup

1. **Update Package Lists**
- Update your package lists to get the latest version of packages.
```bash
sudo apt update
sudo apt upgrade
```

2. **Install PostgreSQL**
- Add the PostgreSQL repository and install PostgreSQL 16 if it is not already installed.
```bash
sudo apt install postgresql-16
```

- Install common and important packages:
```bash
sudo apt install postgresql-16 postgresql-client-16 postgresql-doc-16 libpq-dev postgresql-server-dev-16
```

3. **Install PostGIS**
- Install the relevant PostGIS package for your PostgreSQL version.
```bash
sudo apt install postgis postgresql-16-postgis-3
```

4. **Enable PostGIS Extension**
- Switch to the postgres user to configure your PostgreSQL database.
```bash
sudo -i -u postgres
```

- Create a new database or use an existing one. In this example, a new database called `pi_taxi_gisdb` is created.
```bash
createdb ecommerce
```

- Connect to your database.
```bash
psql -d ecommerce
```
-- Create uuid extension
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

-- Create case-insensitive character string type
```sql
CREATE EXTENSION IF NOT EXISTS citext;
```

- IF NEEDED: (OPTIONAL) Enable the PostGIS extension.
```sql
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
```

-- Set time zone of database to Ashgabat
```sql
ALTER DATABASE ecommerce SET TIME ZONE 'Asia/Ashgabat';
```

5. **Create a New User**
- In the PostgreSQL shell, create a new user. Replace `new_username` and `new_password` with your desired username and password.
```sql
CREATE USER new_username WITH PASSWORD 'new_password';
```

6. **Grant Privileges to the New User**
- Grant all privileges on the `ecommerce` database to the new user.
```sql
GRANT ALL PRIVILEGES ON DATABASE ecommerce TO new_username;
```

7. **Ensure the User Has the Correct Privileges**
- You may also want to grant the new user the ability to create tables, sequences, etc. in the database:
```sql
\c ecommerce
GRANT ALL ON SCHEMA public TO dev;
GRANT ALL PRIVILEGES ON SCHEMA public TO dev;
GRANT CREATE ON SCHEMA public TO dev;
GRANT USAGE ON SCHEMA public TO dev;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO dev;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO dev;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO dev;
```

8. **Test the New User**
- Exit the PostgreSQL shell:
```bash
\q
```

9. Connecting as another user
- Open a new terminal and run the command below
```bash
psql --host=localhost --dbname=ecommerce --username=dev
```
OR

10. OPTIONAL! **Connect to the Database with the New User**
- Try connecting to the `ecommerce` database using the new user to ensure everything is set up correctly. 
- Replace `new_username` with your new username and `ecommerce` with your database name.
```bash
psql -U new_username -d ecommerce
```

- If "Peer authentication failed" error occured:
    - The "Peer authentication failed" error typically means that PostgreSQL is configured to use peer authentication for the user, which relies on the operating system's user account for authentication. This is common for local connections. To resolve this, you need to configure PostgreSQL to use password authentication for your new user.

    - Here’s how you can fix this:

    ##### 1. Edit the PostgreSQL Configuration File
        - Locate the `pg_hba.conf` file, which controls the client authentication configuration. The location of this file can vary, but it's typically found in the PostgreSQL data directory, such as `/etc/postgresql/16/main/pg_hba.conf` on Debian-based systems.

        - Open the file for editing:
        ```bash
        sudo nano /etc/postgresql/16/main/pg_hba.conf
        ```

    ##### 2. Modify the Authentication Method
        - Find the lines that look like this:
        ```
        local   all             all                                     peer
        ```

        - Change the authentication method from `peer` to `md5` or `password` for the `local` connections. For example:
        ```
        local   all             all                                     md5
        ```

    ##### 3. Restart PostgreSQL Service
        - After modifying the `pg_hba.conf` file, restart the PostgreSQL service to apply the changes:
        ```bash
        sudo systemctl restart postgresql
        ```

    ##### 4. Test the Connection Again
        - Now, try to connect to the `pi_taxi_gisdb` database using the new user:
        ```bash
        psql -U new_username -d pi_taxi_gisdb
        ```

    ##### 5. Update Your PostgreSQL User Password
        - If you still encounter issues, ensure that the user password is set correctly and the user exists. Log in as the `postgres` user and set the password for the new user:
        ```bash
        sudo -i -u postgres
        psql
        ```

        - Then, in the PostgreSQL shell:
        ```sql
        ALTER USER new_username WITH PASSWORD 'new_password';
        ```

    ##### 6. Ensure Network Access (if applicable)
        - If you need to connect from a remote machine, make sure that `pg_hba.conf` and `postgresql.conf` are configured to allow remote connections.

        - In `pg_hba.conf`, add a line like:
        ```
        host    all             all             0.0.0.0/0               md5
        ```

        - In `postgresql.conf`, uncomment and set the `listen_addresses` to `'*'`:
        ```
        listen_addresses = '*'
        ```

        - Then restart PostgreSQL again:
        ```bash
        sudo systemctl restart postgresql
        ```

    By following these steps, you should be able to connect to your PostgreSQL database using the new user with password authentication.




### Connecting as another user
```bash
psql --host=localhost --dbname=ecommerce --username=dev
```


### Check where your postgresql.conf file lives - to optimize for performance 
```bash
sudo -u postgres psql -c 'SHOW config_file;'
```

### Link for configs for your machine - https://pgtune.leopard.in.ua/
- For my case: db_version:16, Linux, Web App, RAM: 16GB, CPUs: 4, HDD Storage:
```text
# postgresql.conf

# DB Version: 16
# OS Type: linux
# DB Type: web
# Total Memory (RAM): 16 GB
# CPUs num: 4
# Data Storage: hdd

max_connections = 200 
shared_buffers = 4GB
effective_cache_size = 12GB
maintenance_work_mem = 1GB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 4
effective_io_concurrency = 2
work_mem = 10485kB
min_wal_size = 1GB
max_wal_size = 4GB
max_worker_processes = 4
max_parallel_workers_per_gather = 2
max_parallel_workers = 4
max_parallel_maintenance_workers = 2
```

- SQL commands with alter
```sql
-- DB Version: 16
-- OS Type: linux
-- DB Type: web
-- Total Memory (RAM): 16 GB
-- CPUs num: 4
-- Data Storage: hdd

ALTER SYSTEM SET
 max_connections = '200';
ALTER SYSTEM SET
 shared_buffers = '4GB';
ALTER SYSTEM SET
 effective_cache_size = '12GB';
ALTER SYSTEM SET
 maintenance_work_mem = '1GB';
ALTER SYSTEM SET
 checkpoint_completion_target = '0.9';
ALTER SYSTEM SET
 wal_buffers = '16MB';
ALTER SYSTEM SET
 default_statistics_target = '100';
ALTER SYSTEM SET
 random_page_cost = '4';
ALTER SYSTEM SET
 effective_io_concurrency = '2';
ALTER SYSTEM SET
 work_mem = '10485kB';
ALTER SYSTEM SET
 huge_pages = 'off';
ALTER SYSTEM SET
 min_wal_size = '1GB';
ALTER SYSTEM SET
 max_wal_size = '4GB';
ALTER SYSTEM SET
 max_worker_processes = '4';
ALTER SYSTEM SET
 max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET
 max_parallel_workers = '4';
ALTER SYSTEM SET
 max_parallel_maintenance_workers = '2';
```