#!/bin/bash

###############################################################################################################################
# A helper script that is called by Novus binary to execute sudo commands without a password
# This is possible by defining a sudoers record /etc/sudoers.d/novus

# Arguments
# $1 => name of the action
# $2...$x => action arguments

# This script is owned and editable by root only, so the file modification is not possible without the root password.
###############################################################################################################################

# Define a list of all allowed directories that can be modified by this script
#
# WARNING:
# If the actions below don't respect this list,
# then the script exposes a SERIOUS SECURITY THREAT as ANYONE can run this script
# as SUDO WITHOUT PASSWORD and provide an arbitrary filepath, thus modifying anything in the system
ALLOWED_PATHS=(--ALLOWED-PATHS--)

assert_dir_permission() {
    path="$1"
    matching=false

    # Iterate through all allowed directories
    for allowed_path in "${ALLOWED_PATHS[@]}"; do
    if [[ $path == "$allowed_path"* ]]; then
        matching=true
        break
    fi
    done

    if ! $matching; then
        echo "Cannot perform action on $1. This path is not allowed."
        exit 1
    fi
}

case "$1" in
    "check-ports")
        # $2 => comma-separated list of TCP ports
        if [ -z "$2" ]; then
            echo "Missing ports"
        else
          sudo lsof -nP -i4:"$2"
        fi
        ;;
    "touch")
        # $2 => filepath
        if [ -z "$2" ]; then
            echo "Missing filepath"
        else
            assert_dir_permission "$2"
            sudo touch "$2"
        fi
        ;;
    "mkdir")
        # $2 => filepath
        if [ -z "$2" ]; then
            echo "Missing filepath"
        else
            assert_dir_permission "$2"
            sudo mkdir -p "$2"
        fi
        ;;
    "rm")
        # $2 => filepath
        if [ -z "$2" ]; then
            echo "Missing filepath"
        else
            assert_dir_permission "$2"
            sudo rm "$2"
        fi
        ;;
    "chown")
        # $2 => username
        if [ -z "$2" ]; then
            echo "Missing username"
        # $3 => filepath
        elif [ -z "$3" ]; then
            echo "Missing filepath"
        else
            assert_dir_permission "$3"
            sudo chown "$2" "$3"
        fi
        ;;
    *)
        echo "Undefined action '$1'"
        ;;
esac