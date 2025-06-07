# TM1CTL - TM1 CLI Utility

**WORK-IN-PROGRESS!** Whilst functional as-is, it is still early days, the way things work likely will change!.

A lightweight command-line interface (CLI) for interacting with a TM1. This utility is designed to simplify common TM1 administrative and data operations from the terminal, supporting scripting and automation workflows. The tool is especially suited for users managing TM1 environments or integrating with TM1 programmatically via its OData-compliant REST API.

## Features

* Interact with TM1 models via REST API
* Manage TM1 hosts/instances and authentication
* Ability to restore backup-sets

## Installation

You can download a binary from [releases](#) or build from source:

```bash
git clone https://github.com/Hubert-Heijkers/tm1ctl.git
cd tm1ctl
go install
```

## Configuration

By default, the CLI stores its configuration, hosts, users and your preferences, in:

* **Linux/macOS**: `~/.tm1ctl.json`
* **Windows**: `%USERPROFILE%\.tm1ctl.json`

You can override the default configuration file path using the `--config` global flag.

```bash
tm1ctl --config /path/to/custom-config.json instance list
```

## Output Format

By default, command output is shown in a human-friendly **table** format. You can change the output format to **JSON** using the `--output` global flag:

```bash
tm1ctl --output json instance list
```

Available formats:

* `table` (default)
* `json`

> **Note:** The default format can be changed using the `config` command (see below).

## Usage

```bash
tm1ctl [global options] <command> [flags]
```

### Global Options

| Option     | Description                                |
| ---------- | ------------------------------------------ |
| `--config` | Path to the configuration file to use      |
| `--output` | Output format: `table` or `json`           |
| `--help`   | Show help for any command                  |

### Example

```bash
tm1ctl --config ./dev-config.json --output json instance list
```

## Available Commands

* `tm1ctl config` - Manage global tm1ctl configuration
* `tm1ctl database` - Manage the databases of your TM1 v12 service instance
* `tm1ctl host` - Manage host configuration
* `tm1ctl instance` - Manage the instances of a TM1 v12 service
* `tm1ctl restore` - Performs a database restore using the specified backup-set
* `tm1ctl user` - Manage user's credentials and session variables

These are all the root level command supported by the CLI. See the sections below on how to use each of these commands. 

### Global Configuration Management

Manage the global settings for the CLI.

* `tm1ctl config list` - List all configuration values
* `tm1ctl config set <key> <value>` - Set and save a configuration value

The global configuration settings available 

| Key             | Description                                                    | Updateable  |
| --------------- | -------------------------------------------------------------- | ----------- |
| `output-format` | The default output format to be used to represent the result   | Yes         |
| `host`          | The current active/default host to use if none is provided     | No          |
| `instance`      | The current active/default instance to use if none is provided | No          |
| `user`          | The current active/default user to be used to connect with     | No          |

Updateable implies it can be set using `tm1ctl config set` command. `host`, `instance` and `user` can be set with their respective `use` commands (see their respective sections below).


### Host Management

Manage a collection of named TM1 hosts, including their credentials and service root configuration. Hosts act as reusable, named endpoints that can be switched between or configured independently.

```bash
tm1ctl host [subcommand] [flags]
```

#### Subcommands:

##### `tm1ctl host list`

List all configured hosts and indicate which one is currently active (if any).

```bash
tm1ctl host list
```

##### `tm1ctl host delete <hostName>`

Remove a host from the configuration. If the deleted host is currently active, it will also be unset.

```bash
tm1ctl host delete dev
```

##### `tm1ctl host set <hostName> [flags]`

Add or update the configuration for a host. If the host doesn't exist, it will be created. Use the following flags to define or update host-specific properties:

| Flag                           | Description                                      |
| ------------------------------ | ------------------------------------------------ |
| `--root_client_id <value>`     | Set the root user's client ID for this host      |
| `--root_client_secret <value>` | Set the root user's client secret for this host  |
| `--service_root_url <url>`     | Set the TM1 service root URL for this host       |

You can update multiple values in one command:

```bash
tm1ctl host set dev --service_root_url https://localhost:4444 --root_client_id admin --root_client_secret s3cret
```

##### `tm1ctl host use [<hostName>]`

Set the given host as the active/default host used in subsequent commands. If no name is provided, the active host is unset.

```bash
tm1ctl host use dev         # Use 'dev' as default
tm1ctl host use             # Unset default host
```

#### Notes

* You can override the default host for any command using the `--host` flag:

```bash
tm1ctl instance list --host staging
```

* If no `--host` is provided and no active host is set via `host use`, the command will fail.

### Instance Management

Manage TM1 service **instances** (tenants/environments/namespaces) associated with a given host. Each host can expose one or more named TM1 instances. These are individually addressable and configurable within the context of their host.

```bash
tm1ctl instance [subcommand] [flags]
```

#### Subcommands:

##### `tm1ctl instance list`

List all TM1 service instances defined on the currently active host. Requires a host to be set either via `host use` or the `--host` flag.

```bash
tm1ctl instance list
```

You can also list instances for a specific host explicitly:

```bash
tm1ctl instance list --host staging
```

##### `tm1ctl instance create <instanceName>`

Create a new TM1 instance with the specified name on the currently active host. This will provision a new logical TM1 environment.

```bash
tm1ctl instance create dev01
```

##### `tm1ctl instance delete <instanceName>`

Permanently delete a TM1 instance and all associated content (databases, models, etc.). This action is irreversible.

```bash
tm1ctl instance delete dev01
```

##### `tm1ctl instance use [<instanceName>]`

Set the specified instance as the active/default instance on the current host. Commands that require an instance will use this one unless explicitly overridden via `--instance`.

If no name is provided, the default instance is unset for the current host.

```bash
tm1ctl instance use dev01     # Use 'dev01' as the default instance
tm1ctl instance use           # Unset current instance
```

#### Notes

* The combination of host and instance defines the full context in which TM1 operations take place.
* Use `--instance` to override the default instance for any command:

```bash
tm1ctl cube list --host prod --instance finance
```

* If neither `--instance` nor `instance use` has been set for the active host, commands requiring an instance will produce a clear error.

---

Would you like me to stitch the three command sections (`host`, `instance`, and the rest) together into a full `README.md` now?

##### `tm1ctl instance use [<instanceName>]`

Set the given instance as the active/default instance for the current host. If no name is provided, the default instance is unset.

```bash
tm1ctl instance use finance     # Use 'finance' as the default instance
tm1ctl instance use             # Unset the active instance
```

#### Notes

* Instances are always tied to the active host. Use `tm1ctl host use` or `--host` to control which host you're targeting.

* You can override the default instance for any command by using the `--instance` flag:

```bash
tm1ctl cube list --instance finance
```

* If no `--instance` is given and no instance is active for the current host, commands requiring an instance will fail.

### User Management

Manage credentials and session variables for users that can authenticate with TM1 instances. The CLI maintains a list of reusable users and allows one to be selected as the active/default user.

> **Note:** The CLI treats users independently of hosts and instances. It is the responsibility of the CLI user to ensure that the correct user is active or explicitly specified when issuing commands that require authentication.

```bash
tm1ctl user [subcommand] [flags]
```

#### Subcommands:

##### `tm1ctl user list`

List all configured users and indicate which user is currently active (if any).

```bash
tm1ctl user list
```

##### `tm1ctl user delete <user>`

Remove a user from the configuration. If the user is currently active, they will also be unset.

```bash
tm1ctl user delete admin
```

##### `tm1ctl user set <user> [flags]`

Add a new user or update an existing user with one or more credential or session variables.

| Flag                 | Description                                                                    |
| -------------------- | ------------------------------------------------------------------------------ |
| `--name <value>`     | Set the username                                                               |
| `--password <value>` | Set the password (optional)                                                    |
| `--variables <JSON>` | Set session variables as a JSON object (e.g., `'{"ENV":"dev","REGION":"eu"}'`) |

Example:

```bash
tm1ctl user set admin --name admin@example.com --password s3cret
tm1ctl user set tester --name tester@example.com --variables '{"locale": "en", "theme": "dark"}'
```

##### `tm1ctl user use [<user>]`

Set the specified user as the active/default for authentication. If no name is given, the current user will be unset.

```bash
tm1ctl user use admin    # Use 'admin' as the current user
tm1ctl user use          # Unset active user
```

##### `tm1ctl user variable`

Manage session variables for a specific user. These are typically used to store context-specific values that may be passed to the TM1 service during authentication or execution.

###### `tm1ctl user variable list`

List all session variables associated with the currently active user. Requires a user to be set either via `user --use` or the `--user` flag.

```bash
tm1ctl user variable list
```

You can also list the variables for a specific user explicitly:

```bash
tm1ctl user variable list --user admin
```

###### `tm1ctl user variable set <user> <key> <value>`

Set or update a single session variable for the currently active user. Requires a user to be set either via `user --use` or the `--user` flag.

```bash
tm1ctl user variable set REGION eu-west
```

You can also set or update a single session variable for a specific user explicitly:

```bash
tm1ctl user variable set REGION eu-west --user admin
```

---

#### Notes

* A user does **not** have a password configured for it, the `--password` flag can be used to provide it with any request that requires it. Providing the `--password` flag overwrites any password specified for the user.
* You can override the default user for any operation using the `--user` flag:

```bash
tm1ctl databases list --user tester
```

Absolutely — here's the final documentation section for your `database` command, consistent with the earlier sections and clearly explaining its role in managing TM1 databases per instance.

---

### Database Management

Each TM1 **instance** can host one or more **databases**. A database represents a complete TM1 model, including cubes, dimensions, rules, processes, and associated artifacts. Use the `database` command to create, delete, or list these databases for the active or specified instance.

```bash
tm1ctl database [subcommand] [flags]
```

#### Subcommands:

##### `tm1ctl database list`

List all databases associated with the active instance. Requires a host and instance to be set via `host use` and `instance use`, or specified explicitly using `--host` and `--instance`.

```bash
tm1ctl database list
```

Or for a specific context:

```bash
tm1ctl database list --host staging --instance finance
```

##### `tm1ctl database create <databaseName>`

Create a new TM1 database with the specified name under the active instance. The new database will be empty and ready for model design or content import.

```bash
tm1ctl database create PlanningModel
```

##### `tm1ctl database delete <databaseName>`

Permanently delete a TM1 database, including all of its content (cubes, rules, processes, etc.). This operation is irreversible.

```bash
tm1ctl database delete PlanningModel
```

---

#### Notes

* A valid **host**, **instance**, and **user** must be defined or selected before issuing any `database` operation.
* Use `--host`, `--instance`, and `--user` to override the active context for any operation:

```bash
tm1ctl database list --host prod --instance analytics --user admin
```

Of course — here's the documentation for the `restore` command in the same structured and professional style as the rest:

---

### Database Restore

Restore a TM1 **database** from a previously created **backup-set**. This operation rehydrates the entire database, including all its data and artifacts, from the specified backup.

```bash
tm1ctl restore <backup-set-file> [flags]
```

#### Parameters

* `<backup-set-file>` — Path to the backup-set file on disk (e.g., `./backups/Q1-planning-backup.tgz`). The file must be accessible at runtime.

#### Required Flags

| Flag                | Description                                                                      |
| ------------------- | -------------------------------------------------------------------------------- |
| `--database <name>` | The name of the database to restore into (must exist or be created beforehand)   |
| `--host <name>`     | (Optional if set via `host use`) The host where the instance and database reside |
| `--instance <name>` | (Optional if set via `instance use`) The instance that owns the target database  |
| `--user <name>`     | (Optional if set via `user use`) The user performing the restore                 |

#### Example

```bash
tm1ctl restore ./backups/Q1-planning-backup.tgz \
  --host prod \
  --instance finance \
  --database PlanningModel \
  --user admin
```

If a host, instance, or user has been set using the `use` commands, you may omit the corresponding flags:

```bash
tm1ctl restore ./backups/Q1-planning-backup.tgz --database PlanningModel
```

---

#### Notes

* The file path must point to a valid backup-set created by TM1's Backup operation.
* The database specified with `--database` must already exist.
* Existing data in the target database will be **overwritten** during the restore.
* Authentication via a configured or specified user is required.

## Example Use-cas


Perfect — here's the revised and expanded **Example Use Case** section that incorporates:

* The use of `root-cred.json` from the deployment folder
* Using your [tm1-authenticator-service](https://github.com/hubert-Heijkers/tm1-authenticator-service)
* The domain/email requirement for the `admin` user
* The required authentication to create a database
* The restore step using the migrated TM1 v11 backup-set

---

## Example Use Case: First-Time Setup with TM1 v12 Local Installation

After installing the standalone (local) version of **TM1 v12**, you can use this CLI to configure and initialize your TM1 environment. Below is a typical workflow for a user setting up TM1 for the first time.

### 🛠 Step 1: Configure TM1 for HTTP Passthrough Authentication

Follow the instructions in the TM1 v12 installation guide to enable **HTTP Passthrough Authentication**. This configuration is required to allow secure interaction between this CLI and the TM1 REST API.

During installation, you will have provided a path to a **deployment data folder**. Inside this folder, you will find a file named:

```
root-cred.json
```

This file contains the **root client credentials** needed to connect to the TM1 host. Extract the values for `clientID` and `clientSecret` from this file.

### 🧭 Step 2: Register the Local Host

Using the credentials from `root-cred.json`:

```bash
tm1ctl host set local \
  --service_root_url http://localhost:4444 \
  --root_client_id <clientID-value-from-root-cred.json> \
  --root_client_secret <clientSecret-value-from-root-cred.json>
tm1ctl host use local
```

### 🧰 Step 3: Create an Instance (if you didn't have one created by the installer)

Validate if an initial TM1 instance has been created during installation:

```bash
tm1ctl instance list
```

If the default instance named `tm1` was created you'll see a reponse like this:

```
+----------------------------+------+
|             ID             | NAME |
+----------------------------+------+
| tm1-i-csr2pib478bnq94lugt0 | tm1  |
+----------------------------+------+
```

If you don't have a TM1 instance yet, or wish to create a new TM1 instance:

```bash
tm1ctl instance create demo
```

Optional: Make the instance you are going to use the active/default instance:

```bash
tm1ctl instance use demo
```

### 👤 Step 4: Configure a User

> **Note:** The CLI currently only supports HTTP passthrough authentication which limits it's use to standalone/local installations with HTTP passthrough authentication configured or installations with API key support (albeit in those cases PA based setups, which don't expose TM1's management/instance level APIs, you are limited to any actions that are performed against databases).

If you haven't already we suggest you set up HTTP passthrough authentication either pointing at a v11 TM1 server if you have one of those hanging around or to an instance of my [tm1-authenticator-service](https://github.com/hubert-Heijkers/tm1-authenticator-service). If you do the latter, which we'll presume in this example use case, the following constraints would apply (presuming you are playing with this service as is):

* The username must be an **email address** in the `example.com` domain (e.g., `admin@example.com`)
* The password must be set to **`apple`** (the intentionally well known TM1 admin password)

Then register such a user in the CLI:

```bash
tm1ctl user set admin \
  --name admin@example.com \
  --password apple
tm1ctl user use admin
```

> 🔐 **Note:** Authentication is **required** to create a database, as the user performing the operation becomes the administrator of the new database.

### 🗃 Step 5: Create a Database

Now that authentication is configured, create the database under the active instance:

```bash
tm1ctl database create SalesModel
```

### 💾 Step 6: Restore from a TM1 v11 Backup-Set

If you used the **provided migration utility** to migrate a TM1 v11 database, you should now have a backup-set file representing that model.

Restore the backup into the newly created database:

```bash
tm1ctl restore ./backups/v11-migrated-model.tgz --database SalesModel
```

At this point, your TM1 v12 environment is fully initialized and restored with a migrated model — ready for use and integration.

Want to take a quick peek at your TM1 v12 model? In your [ARC](https://code.cubewise.com/downloads/arc-download/) `settings.yaml` add a connection, not an Admin Host!, specify the database root URL for the database just created as well as the name you want Arc to use as in:

```yaml
connections:
- url: http://localhost:4444/tm1/api/v1/Databases('SalesModel')
  name: Sales
```

Have fun!

## Contributing

Contributions welcome! Please open an issue or submit a pull request if you'd like to improve this utility.

## License

MIT License – see [LICENSE](LICENSE) for details.

