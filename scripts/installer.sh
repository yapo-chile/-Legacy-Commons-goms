#!/bin/bash

# Installer for APIs

#
# Usage:   ./installer.sh [deploy <poya> | copy_to_stg]
#          ./installer.sh copy_to_stg will copy RPM packages to ch9stg
#          ./installer.sh deploy <poya> will deploy those packages to selected poya
# Example: ./installer.sh deploy poya2
# Valid poyas are: poya1, poya2, poya3, poya4
#

# Global variables
api_name=$(grep "Name:" scripts/api.spec | cut -d ":" -f 2 | sed "s/ //g")
packages_path="build/RPMS/x86_64"
api_version=$(grep "Version:" scripts/api.spec | cut -d ":" -f 2 | sed "s/ //g")
repo_path="/opt/repo/yapo"
exec_path="/etc/init.d/${api_name}-api"
temp_dir="/tmp/.${api_name}"
ssh_params="-q -t -t -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"
scp_params="-q -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"
rpm_action=$(if [ $(rpm -qa | grep ${api_name}) ]; then echo "update"; else echo "install"; fi;)
update_params="sudo yum clean all && sudo yum --disablerepo=\"*\" --enablerepo=\"yapo\" ${rpm_action} -y -q \"*${api_name}*\""
poyas_host="lvdev-"
stg_host="ch9stg"

# Given a poya and a server, return a host
# Usage: get_host <poya> <server>
# Example: get_host poya2 www
get_poya_host(){
    # Get number of poya
    case $1 in
        "poya1")
            num="01"
            ;;
        "poya2")
            num="02"
            ;;
        "poya3")
            num="03"
            ;;
        "poya4")
            num="04"
            ;;
        *)
            echo "Error: $i is not a valid poya" 1>&2
            exit 1
            ;;
    esac
    # Get host
    echo jenkins@${poyas_host}${2}${num}
}

# Compress RPM packages and upload to stg
copy_to_stg(){
	tar czvf ${api_name}-${api_version}.tgz ${packages_path}/*.rpm -C ${packages_path}
	scp ${scp_params} ${api_name}-${api_version}.tgz jenkins@${stg_host}:/tmp
	ssh ${ssh_params} jenkins@${stg_host} "rm -f /tmp/${api_name}.tgz && ln -s /tmp/${api_name}-${api_version}.tgz /tmp/${api_name}.tgz"
}

# Deploy to a defined poya
deploy_to_poya(){
    pay_server=$(get_poya_host $1 payment)
    repo_server=$(get_poya_host $1 repo)

    # Step 1: create a tem directory
    ssh ${ssh_params} ${repo_server} "mkdir -p ${temp_dir}"

    # Step 2: upload RPM packages
    scp ${scp_params} ${packages_path}/*.rpm ${repo_server}:${temp_dir}

    # Step 3: move packages to repository
    ssh ${ssh_params} ${repo_server}  "sudo mv ${temp_dir}/*.rpm ${repo_path}/x86_64 && rm -rf ${temp_dir}"

    # Step 4: update repository database
    ssh ${ssh_params} ${repo_server}  "sudo createrepo ${repo_path}"

    # Step 5: update packages in payment host
    ssh ${ssh_params} ${pay_server} ${update_params}

    # Step 6: restart the service
    ssh ${ssh_params} ${pay_server} "sudo ${exec_path} restart"
}

print_help(){
    echo "Usage:   $0 [deploy <poya> | copy_to_stg]"
    echo "         $0 copy_to_stg will tar and copy RPM packages to ${stg_host}"
    echo "         $0 deploy <poya> will deploy those packages to selected poya"
    echo "Example: $0 deploy poya2"
    echo "Valid poyas are: poya1, poya2, poya3 and poya4"
}

# Get command and execute it
command="$1"
selected_poya="$2"
case $command in
    "copy_to_stg")
        copy_to_stg
        ;;
    "deploy")
        deploy_to_poya $selected_poya
        ;;
    "help" | "--help" | "-h")
        print_help
        ;;
    *)
        echo "Error: command ${command} not defined" 1>&2
        exit 1
esac
