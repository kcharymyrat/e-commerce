### Connecting as another user
```bash
psql --host=localhost --dbname=ecommerce --username=dev
```


### Showing corrent user in PostgreSQL
```sql
SELECT current_user;
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