# Ubuntu setup

## golang-migrate
- Download golang-migrate pre-built binary and move it to a location on your system path
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 $GOPATH/bin/migrate
```

- Check the migrate version:
```bash
migrate -version
```

- If there are some problems with migrate -version (does not recognize migrate command):
### 1. **Add Go Bin Directory to Your PATH**
  You'll need to add the `/home/kcharymyrat/go/bin` directory to your `PATH`. This can be done by editing your shell's configuration file (like `~/.bashrc` or `~/.zshrc`).

  - Open your `.bashrc` file (assuming you are using `bash`):
  ```bash
  nano ~/.bashrc
  ```

  - Add the following line at the end of the file:
  ```bash
  export PATH=$PATH:/home/kcharymyrat/go/bin
  ```

  - Save the file and exit (for `nano`, press `CTRL+O`, then `ENTER`, then `CTRL+X`).

  - Apply the changes:
  ```bash
  source ~/.bashrc
  ```

### 2. **Check the Installation**

   After adding the `migrate` command to your `PATH`, you should be able to run:

   ```bash
   migrate -version
   ```

   This command should now return the version of the `migrate` tool.

### 3. **Verify the Go Installation (Optional)**

   If you continue to face issues, you might want to ensure that the Go installation is correctly set up:

   ```bash
   go env GOPATH
   ```


### Notes:
    - You could download from browser
    - Make sure GOPATH Location is Correct
    - Make sure that migrate is the file and NOT the folder that has migrate file in it.