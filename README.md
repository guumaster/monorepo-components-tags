# Gitlab Components Tags

This small tool gets the latest tags of all components on your mono-repo. You can pass a specific commit and it will
 get all tags from that commit or before it.
 
## Usage

You need a `GITLAB_TOKEN` environment variable and the project ID you want to check.

```
NAME:
   gitlab-components-tags - get components tags from gitlab

USAGE:
   gitlab-components-tags [global options] [arguments...]

GLOBAL OPTIONS:
   --project value, -p value  Gitlab ProjectID (required)
   --commit value             Short ID commit from where to start the search
   --base-url value           Gitlab base url (default: https://gitlab.com/api/v4)
   --prefix value             Add prefix to exported names
   --export-shell, -e         format output as shell variables (default: false)
   --help, -h                 (default: false)
   --version, -v              print the version (default: false)

```

Example:
```
export GITLAB_TOKEN=<YOUR_GITLAB_TOKEN>

gitlab-components-tags --project <YOUR_PROJECT_ID> --export-shell

# Output:
export COMPONENT_A_VERSION="1.7.0"
export COMPONENT_B_VERSION="1.1.106"
export COMPONENT_C_VERSION="1.5.0"
```

You can set the environment variables directly doing the following:

```
$> eval $(gitlab-components-tags --project <YOUR_PROJECT_ID> --export-shell)

echo "${COMPONENT_A_VERSION}"
1.7.0
```

### Options

 `--project, -p` Your Gitlab Project ID

 `--commit` Start the process from this commit anb before
 
 `--export-shell, -e` Prepare the output as shell environment variables to use directly on your shell
 
 `--base-url` Set this if your Gitlab instance is different from `https://gitlab.com/api/v4/`
 
 `--prefix` Set this if you want to add a prefix to your environment vars
