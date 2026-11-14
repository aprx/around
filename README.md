
# Around

Around is a tool for managing environment variables in a simple, shareable,
and reproducible way.

It allows you to define environments as code, load secrets from password
vaults, and manage temporary SSH agents.

Around makes it easy to switch between configurations (such as development,
staging, or production) and share them with your team.

----

## Usage

```sh
$ around -h
Usage of ./around:
  -config string
     Path to configuration (default "/home/$HOME/.config/around")
  -output string
     Output path (default "/home/$HOME/.aroundrc")
```

## Example

If you have a command-line application that requires `LOGIN` and `PASSWORD`
environment variables to run.

You can define these variables in `~/.config/around/configuration.ini` like this:

```
[production]
LOGIN=prod-login
PASSWORD=prod-password
```

byB running  `around`, this will generates a shell function for you in `~/.aroundrc`

```
conf_production() {
  around_set_env production
  around_set_clean_env LOGIN:PASSWORD
  export LOGIN=admin
  export PASSWORD=coucou4242
}
```

You can then run `source ~/.aroundrc` and run `conf_production`. It sets up your
shell with the correct environment variables for the production environment.
`~/.aroundrc` is overwritten at every call and should not be modified by hand.

```sh
$ conf_production
$  echo "->$LOGIN $PASSWORD<-"
->admin coucou4242<-
$ conf_unset
$ echo $LOGIN $PASSWORD
-> <-

```

### Using a shared secret vault

If you store secrets in a shared vault (for example, 1Password), you can
reference them directly:

```
PASSWORD=$(op item get prod-password --format json | jq .id)
```

or gopass

```
PASSWORD=$(gopas show /path/to/secret)
```

This way, sensitive values never need to be stored in plain text.

## Configuration

By default, Around:

- Looks for *every* .ini files in

```
$HOME/.config/around/
```

- Generates its output in

```
$HOME/.aroundrc
```

To enable Around in your shell, add the following line to your shell
configuration file (e.g. .bashrc, .zshrc):

```
source ~/.aroundrc
```

### Ini files

`Around` will use the section name to generate the shell function by
prefixing it with `conf_<SECTION>`. Defining 2 sections with the same
name will result in unknown behavior (well last defined will override
the first previous one).

Also it's not possible to name a section `DEFAULT` or `unset`, the section
will be ignored or generate an error.

That’s it, you’re ready to use Around.

## Shell prompt integration

Around sets the `AROUND_NAME` environment variable to indicate which
configuration is currently active.
You can use this to display the active environment directly in your shell prompt.

```sh
PS1='${AROUND_NAME:+[|$AROUND_NAME|]}\u@\h:\w\$ '
```

Here’s how to show the active Around environment in a [Starship](https://starship.rs/config/#environment-variable) prompt:

```toml
[env_var.AROUND_NAME]
variable = "AROUND_NAME"
format = '\[|[$env_value](red)|\] '
```

### SSH Agents

`Around` can spawn a temporary SSH agent, allowing you to load SSH keys stored
in a secure vault.

- `conf_start_ssh_agent` starts a temporary SSH agent, replacing
  `SSH_AUTH_SOCK` and `SSH_AGENT_PID` if they are already set.
- `conf_stop_ssh_agent` stops the temporary agent and restores the previous
  SSH agent environment.
- `conf_unset` removes the loaded environment variables and stops the
  temporary SSH agent, if one was started.

An Around configuration can define both environment variables and SSH keys to load.

Use `around_ssh_key` to specify a command that outputs an SSH private key.
Its output is piped directly to `ssh-add -`, without writing the key to a
temporary file.

Multiple `around_ssh_key` entries are supported.

```ini
[production]
USER=privilegedUser
around_ssh_key=cat /path/to/key
around_ssh_key=gopass show -o=false --noparsing secure/store/to/keys
```
