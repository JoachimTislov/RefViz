# This script can be used to find out if a string is present in other files or not.
# $1 is the string to be searched
# Example of $1: "func \K\w+" or "method \K\w+"
# $2 is the file to be searched
# $3 is the directory to be searched
grep -oP "$1" "$2" | while read -r function; do
    grep -r "$function" "$3"
done

# For example, CreateUser: bash script.sh "CreateUser" $HOME/RefViz/sample-code/quickfeed/database/database.go $HOME/RefViz/sample-code/quickfeed

# Versions of this can be used to map out references in a codebase.