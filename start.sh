#!/bin/bash

socat_from_link() {
	sed 's/.*_PORT_\([0-9]*\)_TCP=tcp:\/\/\(.*\):\(.*\)/socat TCP4-LISTEN:\1,fork,reuseaddr TCP4:\2:\3 \&/'
}

links() {
	echo "Running in links mode..."
	env | grep _TCP= | socat_from_link | echo -e "$(cat)\nwait" | bash
}

setup_iptables() {
	echo "Setting up iptables for container interface..."
	iptables -t nat -A PREROUTING -p tcp -j REDIRECT --to-ports 10000
}

main() {
	case "$1" in
	--links)            links;;   
	--setup-iptables)   setup_iptables;;
	--omnimode)         ambassadord;;
	*)                  ambassadord "$1";;
	esac
}

main "$@"