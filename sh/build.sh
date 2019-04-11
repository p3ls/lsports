#!/bin/bash

name='lsports'
ansible_role_path='/Users/stef/dbox/src/ansible-roles-p3ls/iris-connector-v2'
device_user='connector'
device_ip='192.168.178.23'

build=`date`
version="master/HEAD  ${build}"
# if working dir has no changes, insert git commit id
if [[ $(git diff --stat) != '' ]]; then
  # echo 'dirty'
  version="master/HEAD  ${build}"

else
  # echo 'clean'
  last_commit=$(git log --pretty=oneline --max-count=1|cut -c-10)
  # commit_date=$(git log --pretty=format:'%aD' --max-count=1)
  current_branch=$(git branch | grep \* | cut -d ' ' -f2)
  version="${current_branch}/${last_commit}  ${build}"

fi

ldflags="-w -s \
-X 'main.version=${version}'"


compile () {
	# packr
	CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm go build --ldflags "${ldflags}" -o bin/armhf
	# packr clean
	# stop on error
	if [ $? != 0 ]; then echo "compile error"; exit 1; fi
}

install_on_device() {
	# LC_CTYPE=C tr -dc 'a-zA-Z0-9' < /dev/urandom | fold -w 32 | head -n 1
	ssh "${device_user}@${device_ip}" sudo rm -fr "/tmp/${name}"
	scp bin/armhf "${device_user}@${device_ip}:/tmp/${name}"
	ssh "${device_user}@${device_ip}" sudo "/tmp/${name}"
	# ssh "${device_user}@${device_ip}" sudo systemctl stop "${name}"
	# ssh "${device_user}@${device_ip}" sudo cp "/tmp/${name}" "/usr/local/bin/${name}"
	# ssh "${device_user}@${device_ip}" sudo chmod +x "/usr/local/bin/${name}"
	# ssh "${device_user}@${device_ip}" sudo systemctl start "${name}"
	# # ssh "${device_user}@${device_ip}" "${name}" -version
	# ssh "${device_user}@${device_ip}" sudo systemctl status "${name}"
}

copy_into_role() {

	# shrink
	upx bin/*

	# copy binary into ansible role
	cp bin/armhf "${ansible_role_path}/files/bin/${name}-armhf"

	# set commit hash in ansible role
	sed -i.bak \
	-e "s|^connector_${name}_version: .*|connector_${name}_version: ${version}|g" \
	"${ansible_role_path}/vars/main.yml"

	rm -f "${ansible_role_path}/vars/main.yml.bak"	
}



run_locally () {
	go run *.go
}


usage () {
	echo "usage:"
	echo "-i|install   install on device [${device_user}@${device_ip}]"	
	echo "-c|copy      build + copy to ansible role"
	echo "-r|run       compile and run locally"
	exit 1
}


while [ "$1" != "" ]; do
    case $1 in
    #     -s  )   shift   
    #     SERVER=$1 ;;  
    #     -d  )   shift
    #     DATE=$1 ;;
    # --paramter|p ) shift
    #     PARAMETER=$1;;

		-i|install) shift
					compile
		          	install_on_device ;;

        -c|copy) 	shift
					compile
          			copy_into_role ;;

		-r|run) 	shift
		  			run_locally ;;                

        -h|help)   	usage # function call
                	exit ;;

        - )     usage # All other parameters
                exit 1
    esac
    shift
done