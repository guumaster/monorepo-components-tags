# Monorepo Components Tags

This small tool gets the latest tags of all components on your monorepo. You can pass a specific commit and it will
 get all tags from that commit or before it.


## Usage

You need a `GITLAB_TOKEN` or `GITHUB_TOKEN` as environment variable and the project ID you want to check.

```
NAME:
   monorepo-components-tags - get components tags from a monorepo commits

USAGE:
   monorepo-components-tags [global options] [arguments...]

GLOBAL OPTIONS:
   --project value, -p value  Gitlab ProjectID or Github user/repo (required)

   --commit value             Find tags including only this commit and any older ones. Any short or long SHA are valid

   --base-url value           Gitlab base url (default: https://gitlab.com/api/v4)

   --provider GITLAB|GITHUB   Set this if you want to force an specific provider. Possible values are GITLAB|GITHUB

   --prefix value             Add prefix to exported names

   --suffix value             Change variable names suffix to different value (default: "VERSION")

   --export-shell, -e         Format output as shell variables (default: false)

   --help, -h                 (default: false)

   --version, -v              print the version (default: false)

```

### Gitlab Example:
```
export GITLAB_TOKEN=<YOUR_GITLAB_TOKEN>

monorepo-components-tags --project <YOUR_PROJECT_ID> --export-shell

# Output:
export COMPONENT_A_VERSION="1.7.0"
export COMPONENT_B_VERSION="1.1.106"
export COMPONENT_C_VERSION="1.5.0"
```

### Github Example:
```
export GITHUB_TOKEN=<YOUR_GITHUB_TOKEN>

monorepo-components-tags --project <USER/REPO> --export-shell

# Output:
export COMPONENT_A_VERSION="1.7.0"
export COMPONENT_B_VERSION="1.1.106"
export COMPONENT_C_VERSION="1.5.0"
```


### Exporting environment variables
You can set the environment variables directly doing the following:

```
$> eval $(gitlab-components-tags --project <YOUR_PROJECT_ID> --export-shell)

echo "${COMPONENT_A_VERSION}"
1.7.0
```

### Options

 `--project, -p` Your Gitlab Project ID or Github project in the form `user/repo`

 `--commit` Find tags including only this commit and any older ones. Any short or long SHA are valid.

 `--base-url` Set this if your Gitlab instance is different from `https://gitlab.com/api/v4/` (no need for Github)

 `--provider` Set this if you want to force a specific provider. Possible values are `GITLAB|GITHUB`

 `--prefix` Set this if you want to add a prefix to your environment vars

 `--suffix` Change suffix from 'VERSION' to different value (default: "VERSION")

 `--export-shell, -e` Prepare the output as shell environment variables to use directly on your shell
