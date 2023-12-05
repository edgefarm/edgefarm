#!/bin/bash
green='\033[0;32m'
red='\033[0;31m'
yellow='\033[0;33m'
nc='\033[0m'
blue='\033[0;36m'

suported_ubuntu_versions=("22.04")

check_ubuntu_version() {
  version=$(lsb_release -rs)
  if [[ " ${suported_ubuntu_versions[*]} " =~ " ${version} " ]]; then
    echo "true"
  else
    echo "false"
  fi
}

check_swap_disabled() {
  if free | awk '/^Swap:/ {exit !$2}'; then
    echo "false"
  else
    echo "true"
  fi
}

check_conntrack() {
  conntrack > /dev/null 2>&1
  if [ $? -eq 127 ]; then
     echo "false"
  else
     echo "true"
  fi

}

check_socat() {
  socat > /dev/null 2>&1
  if [ $? -eq 127 ]; then
     echo "false"
  else
     echo "true"
  fi
}

check_docker() {
  docker info > /dev/null 2>&1
  if [ $? -eq 0 ]; then
     echo "true"
  else
     echo "false"
  fi
}

check_netbird() {
  netbird status > /dev/null 2>&1
  if [ $? -eq 0 ]; then
     echo "true"
  else
     echo "false"
  fi
}


# Check if the script is run as root
if [ "$EUID" -ne 0 ]
 then echo "Please run this script as root"
 exit
fi

MISSING_PACKAGES=""

echo -n "Checking Ubuntu version... "
UBUNTU_VERSION="$(check_ubuntu_version)"
if [[ "$UBUNTU_VERSION" == "false" ]]; then
  echo -e "${yellow}this version of Ubuntu is untested. The script will continue anyway. This might not work. Supported versions are ${suported_ubuntu_versions[@]}${nc}"
else
  echo -e "${green}supported${nc}"
fi

echo -n "Checking conntrack... "
CONNTRACK_PRESENT="$(check_conntrack)"
if [[ "$CONNTRACK_PRESENT" == "false" ]]; then
  echo -e "${red}missing${nc}"
  MISSING_PACKAGES+=" conntrack"
else
  echo -e "${green}installed${nc}"
fi

echo -n "Checking socat... "
SOCAT_PRESENT="$(check_socat)"
if [[ "$SOCAT_PRESENT" == "false" ]]; then
  echo -e "${red}missing${nc}"
  MISSING_PACKAGES+=" socat"
else
  echo -e "${green}installed${nc}"
fi

echo -n "Checking Docker... "
DOCKER_PRESENT="$(check_docker)"
if [[ "$DOCKER_PRESENT" == "false" ]]; then
  echo -e "${red}missing${nc}"
else
  echo -e "${green}installed${nc}"
fi

echo -n "Checking Netbird... "
NETBIRD_PRESENT="$(check_netbird)"
if [[ "$NETBIRD_PRESENT" == "false" ]]; then
  echo -e "${red}missing${nc}"
else
  echo -e "${green}installed${nc}"
fi

echo -n "Checking if swap is disabled... "
SWAP_DISABLED="$(check_swap_disabled)"
if [[ "$SWAP_DISABLED" == "false" ]]; then
  echo -e "${red}enabled{nc}"
else
  echo -e "${green}disabled${nc}"
fi

echo $MISSING_PACKAGES

if [[ "$MISSING_PACKAGES" != "" ]]; then
  apt update
  apt install ${MISSING_PACKAGES}
fi

if [[ "$DOCKER_PRESENT" == "false" ]]; then
  curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
  sh /tmp/get-docker.sh
  rm /tmp/get-docker.sh
fi

if [[ "$NETBIRD_PRESENT" == "false" ]]; then
  echo -e "${red}Netbird missing. Installing...${nc}"
  sudo apt-get update
  sudo apt-get install ca-certificates curl gnupg -y
  curl -sSL https://pkgs.netbird.io/debian/public.key | sudo gpg --dearmor --output /usr/share/keyrings/netbird-archive-keyring.gpg
  echo 'deb [signed-by=/usr/share/keyrings/netbird-archive-keyring.gpg] https://pkgs.netbird.io/debian stable main' | sudo tee /etc/apt/sources.list.d/netbird.list
  sudo apt-get update
  sudo apt-get install netbird
fi

if [[ "$SWAP_DISABLED" == "false" ]]; then
  swapoff -a                               # Disable all devices marked as swap in /etc/fstab
  sed -e '/swap/ s/^#*/#/' -i /etc/fstab   # Comment the correct mounting point
  systemctl mask swap.target               # Completely disabled
fi

echo -e "Everything installed. Next, make sure to connect to netbird using"
echo -e "${blue}netbird up --setup-key <YOUR-SETUP-KEY>${nc}"
