# quick-workflow

> The shell command based on jira-cli and github-cli for quick work

## Two quick commands for streamlined workflow:

1. `pr-create`: Automatically create a Git branch, submit a PR (Pull Request), capture PR details, link the PR to Jira, and update its status.
2. `pr-merge`: Automatically merge a PR, delete the Git branch, and update the Jira status in one swift command.

### Why use it

> Are you tired of the hassle of changing the status every time you create a PR, such as 'code review,' 'merged,' and dealing with repetitive tasks?

> Especially when you're busy and have opened multiple browser windows, waiting for pages to load just to complete these mundane tasks can be time-consuming, diverting our focus from more important matters.

Highlighted Benefits:

- Efficiency: Streamline the process by automating status updates and repetitive tasks.
- Time-Saving: Eliminate the need to open multiple browser tabs and wait for pages to load.
- Enhanced Focus: Allow users to concentrate on more crucial tasks.
- Reduced Manual Work: Minimize the need for manual status updates for each PR.
- Standardized Workflow: Ensures a consistent and error-free process.

### How to use it

> Fortunately, we have two powerful tools, github-cli and jira-cli. I use shell scripts to standardize our PR workflow and achieve the above with just one command. You only need to complete the following steps step by step:

1. Install github-cli, here are some verifications, please refer to the official documentation see: https://github.com/cli/cli,

```shell
brew install gh
```

Then need to auth it:

```shell
gh auth login
```

2. use `gh` command to clone this project to your local computer:
```shell
  gh repo clone Wangggym/quick-workflow
```

3. Install jira-cli, here are some verifications, please refer to the official documentation see: https://github.com/ankitpokhrel/jira-cli

```
brew tap ankitpokhrel/jira-cli
brew install jira-cli
```

Then need to auth it:

  -  [Get a Jira API token](https://id.atlassian.com/manage-profile/security/api-tokens) and export it to your shell as
   a `JIRA_API_TOKEN` variable. Add it to your shell configuration file, for instance, `$HOME/.bashrc`, so that the
   variable is always available. Alternatively, you can also use `.netrc` file or `keychain` to set the token. Learn
   more [here](https://github.com/ankitpokhrel/jira-cli/discussions/356).
  - Run `jira init`, select installation type as `Cloud`, and provide required details to generate a config file required
   for the tool.

4. Install `jq` for shell because we need to deal with JSON structure.

```
brew install jq

```

5. add global variables and alias to your `.zshrc` or `.bash_profile`

```shell
vim ~/.zshrc
```

```shell

alias pr-create=/Users/xxx/xxx/quick-workflow/pr-create.sh
alias pr-merge=/Users/xxx/xxx/quick-workflow/pr-merge.sh

export JIRA_SERVICE_ADDRESS=https://xxx.xx # Your Jira network address
export JIRA_API_TOKEN=xxx

# Optional

# Generate the custom branch prefix name
# export GH_BRANCH_PREFIX=xx # xx/jira_ticket--desc

```

6. Make them to be able to execute:

```shell
chmod +x /Users/xxx/xxx/quick-workflow/pr-create.sh
chmod +x /Users/xxx/xxx/quick-workflow/pr-merge.sh
```

Now you can use it just input `pr-create` or `pr-merge` in your command.

### Updates

See `todolist.md`, you can see the latest updates and other important info.

Thank you very much for your valuable feedback and PR, let's work together to make these commands better
