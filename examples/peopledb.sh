#!/usr/bin/env bash

#mode=text
mode=json

db=/tmp/people

main() {

   for keyval in $*;do
      key=$(echo ${keyval}|awk -F= '{print $1}')
      val=$(echo ${keyval}|awk -F= '{print $2}')
      if [[ "${key}" == "name" ]];then
         name="${val}"
      elif [[ "${key}" == "age" ]];then
         age="${val}"
      elif [[ "${key}" == "fn" ]];then
         function="${val}"
      else
         error "unknown param ${keyval}"
         exit 1
      fi
   done

   [[ ! -f "${db}" ]] && touch "${db}"
   if [[ "${function}" == "update" || "${function}" == "add" ]];then
      if [[ -z "${age}" ]];then
         error "No age passed with update function"
         exit 2
      elif [[ -z "${name}" ]];then
         error "No name passed with update function"
         exit 2
      else
         do_update ${name} ${age}
      fi
   elif [[ "${function}" == "remove" ]];then
      if [[ -z "${name}" ]];then
         error "No name passed with remove function"
         exit 3
      else
         do_remove ${name}
      fi
   elif [[ "${function}" == "list" ]];then
      if [[ -z "${name}" ]];then
	 [[ "${mode}" == "json" ]] && do_json_list || do_list
      else
	 [[ "${mode}" == "json" ]] && do_json_single_list ${name} || do_single_list ${name}
      fi
   else
      error "unknown function ${function}"
      exit 4
   fi
}

do_update() {
  do_remove "${1}" >/dev/null
  echo "${1}:${2}" >> "${db}"
  message "updated ${1}"
} 

do_remove() {
  grep -v "${1}" "${db}" > "/tmp/t.$$"
  mv "/tmp/t.$$" "${db}"
  message "removed ${1}"
}

do_list() {
  concatenated_data=""
  for each_name in $(awk -F: '{print $1}' "${db}");do
    record=$(grep "${each_name}" "${db}")
    record_name=$(echo "${record}"|awk -F: '{print $1}')
    record_age=$(echo "${record}"|awk -F: '{print $2}')
    echo "${record_name} is ${record_age} years old"
  done
}

do_single_list() {
  record=$(grep ^"${1}:" "${db}")
  record_name=$(echo "${record}"|awk -F: '{print $1}')
  record_age=$(echo "${record}"|awk -F: '{print $2}')
  if [[ "x${record_name}x" != "xx" ]];then
     echo "${record_name} is ${record_age} years old"
  else
     echo "${1} not found in peopledb"
  fi
}

do_json_list() {
  concatenated_data=""
  echo -n "'message': ["
  i=0
  num=$(cat ${db}|wc -l)
  for each_name in $(awk -F: '{print $1}' "${db}");do
    i=$(( $i + 1 ))
    record=$(grep "${each_name}" "${db}")
    record_name=$(echo "${record}"|awk -F: '{print $1}')
    record_age=$(echo "${record}"|awk -F: '{print $2}')
    echo -n "{"
    echo -n "'name': '${record_name}'",
    echo -n "'age': '${record_age}'"
    if [[ $i -eq $num ]];then
        echo -n "}"
    else
        echo -n "},"
    fi
  done
  echo -n "]"
}

do_json_single_list() {
  record=$(grep ^"${1}:" "${db}")
  record_name=$(echo "${record}"|awk -F: '{print $1}')
  record_age=$(echo "${record}"|awk -F: '{print $2}')
  if [[ "x${record_name}x" != "xx" ]];then
    echo -n "["
    echo -n "{"
    echo -n "'name': '${record_name}'",
    echo -n "'age': '${record_age}'"
    echo -n "}"
    echo -n "]"
  else
    error "user ${1} not present in peopledb"
  fi
}

error() {
   if [[ "${mode}" == "json" ]];then
      formatted_json_print "error" "${*}"
   else
      echo "ERROR: ${*}"
   fi
   exit 1
}

message() {
   if [[ "${mode}" == "json" ]];then
     formatted_json_print "message" ${*}
   else
     echo "${*}"
   fi
}

formatted_json_print() {
   echo -n "'${1}': '${2}'"
}

main $*

exit 0
