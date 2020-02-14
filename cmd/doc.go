// Package main for contains the code to execute is as a CLI tool.
//
/*

NAME:
   monorepo-components-tags - get components tags from a git monorepo

USAGE:
   monorepo-components-tags [global options] [arguments...]

VERSION:
   v1.0.0

GLOBAL OPTIONS:
   --project user/repo, -p user/repo  Your Gitlab Project ID or Github project in the form user/repo (required)
   --commit value                     Find tags including only this commit and any older ones. Any short or long SHA are valid
   --base-url value                   Gitlab base url (default: https://gitlab.com/api/v4)
   --provider GITLAB|GITHUB           Set this if you want to force an specific provider. Possible values are GITLAB|GITHUB
   --prefix value                     Add prefix to exported names
   --suffix value                     Change variable names suffix to different value (default: "VERSION")
   --export-shell, -e                 format output as shell variables (default: false)
   --help, -h                         (default: false)
   --version, -v                      print the version (default: false)


*/
package main
