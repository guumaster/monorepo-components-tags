# Gitlab Components Tags

This small tool gets the latest tags of all components on your mono-repo. You can pass a specific commit and it will
 get all tags from that commit or before it.
 
## Usage

You need a `GITLAB_TOKEN` environment variable and the project ID you want to check.

```
gitlab-components-tags --project <YOUR_PROJECT_ID> [--commit <COMMIT_SHORT_ID>] [--export-shell]  [--prefix <SOME_STRING>] [--base-url <YOUR_GITLAB_URL>]
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

 `--project` Your Gitlab Project ID

 `--commit` Start the process from this commit anb before
 
 `--export-shell` Prepare the output as shell environment variables to use directly on your shell
 
 `--base-url` Set this if your Gitlab instance is different from `https://gitlab.com/api/v4/`
 
 `--prefix` Set this if you want to add a prefix to your environment vars
