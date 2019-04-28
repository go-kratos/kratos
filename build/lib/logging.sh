#!/bin/bash

# Copyright 2014 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Controls verbosity of the script output and logging.
KRATOS_VERBOSE="${KRATOS_VERBOSE:-5}"

# Handler for when we exit automatically on an error.
# Borrowed from https://gist.github.com/ahendrix/7030300
kratos::log::errexit() {
  local err="${PIPESTATUS[@]}"

  # If the shell we are in doesn't have errexit set (common in subshells) then
  # don't dump stacks.
  set +o | grep -qe "-o errexit" || return

  set +o xtrace
  local code="${1:-1}"
  # Print out the stack trace described by $function_stack  
  if [ ${#FUNCNAME[@]} -gt 2 ]
  then
    kratos::log::error "Call tree:"
    for ((i=1;i<${#FUNCNAME[@]}-1;i++))
    do
      kratos::log::error " $i: ${BASH_SOURCE[$i+1]}:${BASH_LINENO[$i]} ${FUNCNAME[$i]}(...)"
    done
  fi  
  kratos::log::error_exit "Error in ${BASH_SOURCE[1]}:${BASH_LINENO[0]}. '${BASH_COMMAND}' exited with status $err" "${1:-1}" 1
}

kratos::log::install_errexit() {
  # trap ERR to provide an error handler whenever a command exits nonzero  this
  # is a more verbose version of set -o errexit
  trap 'kratos::log::errexit' ERR

  # setting errtrace allows our ERR trap handler to be propagated to functions,
  # expansions and subshells
  set -o errtrace
}

# Print out the stack trace
#
# Args:
#   $1 The number of stack frames to skip when printing.
kratos::log::stack() {
  local stack_skip=${1:-0}
  stack_skip=$((stack_skip + 1))
  if [[ ${#FUNCNAME[@]} -gt $stack_skip ]]; then
    echo "Call stack:" >&2
    local i
    for ((i=1 ; i <= ${#FUNCNAME[@]} - $stack_skip ; i++))
    do
      local frame_no=$((i - 1 + stack_skip))
      local source_file=${BASH_SOURCE[$frame_no]}
      local source_lineno=${BASH_LINENO[$((frame_no - 1))]}
      local funcname=${FUNCNAME[$frame_no]}
      echo "  $i: ${source_file}:${source_lineno} ${funcname}(...)" >&2
    done
  fi
}

# Log an error and exit.
# Args:
#   $1 Message to log with the error
#   $2 The error code to return
#   $3 The number of stack frames to skip when printing.
kratos::log::error_exit() {
  local message="${1:-}"
  local code="${2:-1}"
  local stack_skip="${3:-0}"
  stack_skip=$((stack_skip + 1))

  if [[ ${KRATOS_VERBOSE} -ge 4 ]]; then
    local source_file=${BASH_SOURCE[$stack_skip]}
    local source_line=${BASH_LINENO[$((stack_skip - 1))]}
    echo "!!! Error in ${source_file}:${source_line}" >&2
    [[ -z ${1-} ]] || {
      echo "  ${1}" >&2
    }

    kratos::log::stack $stack_skip

    echo "Exiting with status ${code}" >&2
  fi

  exit "${code}"
}

# Log an error but keep going.  Don't dump the stack or exit.
kratos::log::error() {
  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo "!!! $timestamp ${1-}" >&2
  shift
  for message; do
    echo "    $message" >&2
  done
}

# Print an usage message to stderr.  The arguments are printed directly.
kratos::log::usage() {
  echo >&2
  local message
  for message; do
    echo "$message" >&2
  done
  echo >&2
}

kratos::log::usage_from_stdin() {
  local messages=()
  while read -r line; do
    messages+=("$line")
  done

  kratos::log::usage "${messages[@]}"
}

# Print out some info that isn't a top level status line
kratos::log::info() {
  local V="${V:-0}"
  if [[ $KRATOS_VERBOSE < $V ]]; then
    return
  fi

  for message; do
    echo "$message"
  done
}

# Just like kratos::log::info, but no \n, so you can make a progress bar
kratos::log::progress() {
  for message; do
    echo -e -n "$message"
  done
}

kratos::log::info_from_stdin() {
  local messages=()
  while read -r line; do
    messages+=("$line")
  done

  kratos::log::info "${messages[@]}"
}

# Print a status line.  Formatted to show up in a stream of output.
kratos::log::status() {
  local V="${V:-0}"
  if [[ $KRATOS_VERBOSE < $V ]]; then
    return
  fi

  timestamp=$(date +"[%m%d %H:%M:%S]")
  echo "+++ $timestamp $1"
  shift
  for message; do
    echo "    $message"
  done
}
