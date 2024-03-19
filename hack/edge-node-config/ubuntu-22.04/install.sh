#!/bin/bash

green='\033[0;32m'
red='\033[0;31m'
yellow='\033[0;33m'
yellowBold='\033[1;33m'
nc='\033[0m'

DEB_PACKAGES_RELEASE_URL=https://github.com/edgefarm/edgefarm/releases/download/k8s-1.22.17-deb
KUBEADM_DEB_PKG=kubeadm_1.22.17-00_amd64_7b7456beaf364ecf5c14f4d995bc49985cd23273ebf7610717961e2575057209.deb
KUBECTL_DEB_PKG=kubectl_1.22.17-00_amd64_b3bcd8e4a64fded2873e873301ef68c6c3787dbc5e68f079a2f9c7c283180709.deb
KUBELET_DEB_PKG=kubelet_1.22.17-00_amd64_3488568197f82b8b8c267058ea7165968560a67daa5cea981ac6bcff43fe0966.deb
KUBERNETS_CNI_DEB_PKG=kubernetes-cni_1.2.0-00_amd64_0c2be3775ea591dee9ce45121341dd16b3c752763c6898adc35ce12927c977c1.deb
CRI_TOOLS_DEB_PKG=cri-tools_1.26.0-00_amd64_5ba786e8853986c7f9f51fe850086083e5cf3c3d34f3fc09aaadd63fa0b578df.deb

if [ "$EUID" -ne 0 ]; then
  echo "Please run this script as root."
  exit 1
fi

# Function to check reachability
check_reachability() {
    # Extract host and port from the input
    input=$1
    host=$(echo "$input" | cut -d: -f1)
    port=$(echo "$input" | cut -d: -f2)

    # Check if the host and port are reachable
    if nc -zv -w 2 "$host" "$port" >/dev/null 2>&1; then
        return 0  # Success
    else
        return 1  # Failure
    fi
}

is_valid_format() {
    input=$1

    # Use a regular expression to check the format
    if [[ $input =~ ^[a-zA-Z0-9.-]+:[0-9]+$ ]]; then
        return 0  # Valid format
    else
        return 1  # Invalid format
    fi
}


###### PARAMETERS HANDLING
NAME=${NAME:-$(hostname)}
TOKEN=${TOKEN:-}
ARCH=${ARCH:-$(uname -m)}
ADDRESS=${ADDRESS:-}
JOIN=${JOIN:-false}
PRECHECKS_ONLY=${PRECHECKS_ONLY:-false}
NODE_IP=${NODE_IP:-}
NODE_TYPE=${NODE_TYPE:-"kubeadm"}
CONVERT=${CONVERT:-false}
NO_INSTALL=${NO_INSTALL:-false}
VERSION=${VERSION:-"1.22.17"}

options=$(getopt -o "" -l "prechecks-only,join,node-ip:,address:,name:,token:,arch:,node-type:,convert,no-download,version:,help" -- "$@")

if [ $? -ne 0 ]; then
 echo "Invalid arguments."
 exit 1
fi

eval set -- "$options"
Help() {
 echo "Usage: script --address ADDRESS --token TOKEN --node-ip IP [--name NAME] [--arch ARCH] [--node-type TYPE] [--convert] [--join] [--no-download] [--help]"
 echo
 echo "Options:"
 echo "--address ADDRESS   Set the API server address. This option is mandatory."
 echo "--token TOKEN       Set the token. This option is mandatory."
 echo "--name NAME         Set the node name. Default is the hostname. (optional)"
 echo "--arch ARCH         Set the architecture. Allowed values are arm64, amd64, and arm. Default is the current architecture. (optional)"
 echo "--node-ip IP        Set the node ip."
 echo "--join              Run the join command rather than printing it after setting everything up. (optional)"
 echo "--prechecks-only    Run prechecks only. (optional)"
 echo "--node-type         Set the type of the node. Allowed values are 'kubeadm' or 'yurtadm'. Default is 'kubeadm'. (optional)"
 echo "--convert           Sets the label 'node.edgefarm.io/to-be-converted' to true. Only used for kubeadm type. Default is false. (optional)"
 echo "--no-install        Disables installation of the components (kubeadm, kubelet, kubectl, cni plugins). Default is false (optional)"
 echo "--version           Set the kubernetes version. Used to download the components. Default is 1.22.17 (optional)"
 echo "--help              Display this help message."
 echo
}

while [ $# -gt 0 ]; do
 case "$1" in
  --address) ADDRESS="$2"; shift;;
  --name) NAME="$2"; shift;;
  --token) TOKEN="$2"; shift;;
  --arch) ARCH="$2"; shift;;
  --node-ip) NODE_IP="$2"; shift;;
  --node-type) NODE_TYPE="$2"; shift;;
  --convert) CONVERT="true";;
  --no-install) NO_INSTALL="true";;
  --help) Help; exit;;
  --join) JOIN="true";;
  --version) VERSION="$2"; shift;;
  --prechecks-only) PRECHECKS_ONLY="true";;
  --) shift;;
 esac
 shift
done

if [ -z "$TOKEN" ]; then
 echo -e "${red}Token must be set.${nc}"
 exit 1
fi

if [ -z "$ADDRESS" ]; then
 echo -e "${red}Address must be set.${nc}"
 exit 1
fi

if [ -z "$NODE_IP" ]; then
 echo -e "${red}node-ip must be set.${nc}"
 exit 1
fi

# Check if the address format is valid
if ! is_valid_format "$ADDRESS"; then
    echo "$ADDRESS is invalid. Maybe you forgot to add the port?"
    exit 1
fi

# Check reachability
if ! check_reachability "$ADDRESS"; then
    echo "$ADDRESS is not reachable. Maybe you made a typo? Format must be 'host:port'."
    exit 1
fi

# Map uname architecture to specific values
case "$ARCH" in
 "x86_64") ARCH="amd64" ;;
 "aarch64") ARCH="arm64" ;;
 "armv7l") ARCH="arm" ;;
 *) echo -e "${red}Invalid architecture. Allowed values are arm64, amd64 and arm.${nc}"; exit 1 ;;
esac

TMP=$(mktemp -d)

###### PRECHECKS
PRECHECK_ERRORS=0
echo "Running prechecks..."
INTERFACE="wt0"
if ip link show "$INTERFACE" | grep -qs "state UP"; then
   echo -e "  ${red}Interface $INTERFACE is not up or does not exist. Make sure that netbird is installed, up and running.${nc}"
   PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
else
    echo -e "  ${green}Interface $INTERFACE is up.${nc}"
fi

# Try to run a Docker command
docker info > /dev/null 2>&1
# Check the exit status
if [ $? -eq 0 ]; then
   echo -e "  ${green}Docker is running${nc}"
else
   echo -e "  ${red}Docker is not running${nc}"
   PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
fi

if ! $NO_INSTALL; then
  if [[ $NODE_TYPE == *"yurtadm"* ]]; then
    echo -e "${yellowBold}Downloading openyurt components${NC}"
    wget -q --show-progress https://github.com/openyurtio/openyurt/releases/download/v1.4.0/yurtadm-v1.4.0-linux-${ARCH}.tar.gz -P ${TMP}
    tar xfz ${TMP}/yurtadm-v1.4.0-linux-${ARCH}.tar.gz  -C ${TMP} && mv ${TMP}/linux-${ARCH}/yurtadm /usr/local/bin/yurtadm && chmod +x /usr/local/bin/yurtadm
  else
    INSTALL_KUBEADM=false
    kubeadm version > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "  ${red}kubeadm missing${nc}"
        INSTALL_KUBEADM=true
        PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
    else
        KUBEADM_VERSION=$(kubeadm version | awk -F "GitVersion:\"v" '{print $2}' | awk -F "\"" '{print $1}')
        if [ "$KUBEADM_VERSION" != "$VERSION" ]; then
            echo -e "  ${red}kubeadm version mismatch. Found $KUBEADM_VERSION, expected $VERSION${nc}"
            INSTALL_KUBEADM=true
            PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
        fi
    fi

    INSTALL_KUBELET=false
    kubelet --version > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "  ${red}kubelet missing${nc}"
        INSTALL_KUBELET=true
        PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
    else
        KUBELET_VERSION=$(kubelet --version | awk -F "Kubernetes v" '{print $2}')
        if [ "$KUBELET_VERSION" != "$VERSION" ]; then
            echo -e "  ${red}kubelet version mismatch. Found $KUBELET_VERSION, expected $VERSION${nc}"
            INSTALL_KUBELET=true
            PRECHECK_ERRORS=$((PRECHECK_ERRORS+1))
        fi
    fi

    if [ "$INSTALL_KUBEADM" == "true" ] || [ "$INSTALL_KUBELET" == "true" ]; then
        echo -e "${yellowBold}Installing k8s packages for version $VERSION${NC}"
        wget -q --show-progress ${DEB_PACKAGES_RELEASE_URL}/${KUBEADM_DEB_PKG} -P ${TMP}
        wget -q --show-progress ${DEB_PACKAGES_RELEASE_URL}/${KUBECTL_DEB_PKG} -P ${TMP}
        wget -q --show-progress ${DEB_PACKAGES_RELEASE_URL}/${KUBELET_DEB_PKG} -P ${TMP}
        wget -q --show-progress ${DEB_PACKAGES_RELEASE_URL}/${KUBERNETS_CNI_DEB_PKG} -P ${TMP}
        wget -q --show-progress ${DEB_PACKAGES_RELEASE_URL}/${CRI_TOOLS_DEB_PKG} -P ${TMP}
        apt install -y --fix-broken ${TMP}/${KUBEADM_DEB_PKG} ${TMP}/${KUBECTL_DEB_PKG} ${TMP}/${KUBELET_DEB_PKG} ${TMP}/${KUBERNETS_CNI_DEB_PKG} ${TMP}/${CRI_TOOLS_DEB_PKG} --reinstall --allow-downgrades
        apt-mark hold kubelet kubectl kubeadm
        rm ${TMP}/${KUBEADM_DEB_PKG} ${TMP}/${KUBECTL_DEB_PKG} ${TMP}/${KUBELET_DEB_PKG} ${TMP}/${KUBERNETS_CNI_DEB_PKG} ${TMP}/${CRI_TOOLS_DEB_PKG}
    fi
  fi
fi

if [ $PRECHECK_ERRORS -gt 0 ]; then
    echo -e "${red}Prechecks failed.${nc}"
    exit 1
fi

echo "Prechecks passed."
###### INSTALLATION

if $PRECHECKS_ONLY == "true" ; then
    exit 0
fi

mkdir -p /usr/local/bin
mkdir -p /opt/cni/bin
mkdir -p /etc/edgefarm/
mkdir -p /etc/systemd/system
mkdir -p /etc/udev/rules.d



LABELSAPPEND=""
if [[ $NODE_TYPE == *"yurtadm"* ]]; then
  LABELSAPPEND="node.edgefarm.io/converted=true"
else 
  if $CONVERT ; then
    LABELSAPPEND="node.edgefarm.io/to-be-converted=true"
  else
    LABELSAPPEND="node.edgefarm.io/to-be-converted=false"
  fi
fi

cp files/kubeadm-join.conf.template ${TMP}/kubeadm-join.conf
sed -i "s#LABELSAPPEND#$LABELSAPPEND#g" ${TMP}/kubeadm-join.conf
sed -i "s/ADDRESS/$ADDRESS/g" ${TMP}/kubeadm-join.conf
sed -i "s/NODE_NAME/$NAME/g" ${TMP}/kubeadm-join.conf
sed -i "s/BOOTSTRAP_TOKEN/$TOKEN/g" ${TMP}/kubeadm-join.conf

if [ -n "${NODE_IP}" ]; then
  echo "    node-ip: ${NODE_IP}" >> ${TMP}/kubeadm-join.conf
fi

cp ${TMP}/kubeadm-join.conf /etc/edgefarm/
rm -rf ${TMP}

cp files/etc/systemd/system/edge0-device.service /etc/systemd/system/
cp files/etc/udev/rules.d/90-wt0.rules /etc/udev/rules.d/
cp files/usr/local/bin/add-edge0.sh /usr/local/bin/
cp files/usr/local/bin/add-wt0.sh /usr/local/bin/
cp files/usr/local/bin/remove-wt0.sh /usr/local/bin/

systemctl enable edge0-device
systemctl start edge0-device
udevadm control --reload-rules
udevadm trigger
/usr/local/bin/add-wt0.sh

###### JOIN CLUSTER
if $JOIN ; then
  echo -e "${green}Joining the cluster...${nc}"
  if [[ $NODE_TYPE == *"yurtadm"* ]]; then
    yurtadm join ${ADDRESS} --config /etc/edgefarm/kubeadm-join.conf --node-name=${NAME} --token=${TOKEN} --node-type=edge --discovery-token-unsafe-skip-ca-verification --v=9 --reuse-cni-bin --yurthub-image ghcr.io/openyurtio/openyurt/yurthub:v1.4.0 --cri-socket /var/run/dockershim.sock --yurthub-server-addr=192.168.168.1
  else 
    kubeadm join --config /etc/edgefarm/kubeadm-join.conf -v5
  fi
else
  echo -e "${green}Everything is set up. Run the following command to join the cluster:${nc}"
  if [[ $NODE_TYPE == *"yurtadm"* ]]; then
    echo yurtadm join ${ADDRESS} --config /etc/edgefarm/kubeadm-join.conf --node-name=${NAME} --token=${TOKEN} --node-type=edge --discovery-token-unsafe-skip-ca-verification --v=9 --reuse-cni-bin --yurthub-image ghcr.io/openyurtio/openyurt/yurthub:v1.4.0 --cri-socket /var/run/dockershim.sock --yurthub-server-addr=192.168.168.1
  else 
    echo kubeadm join --config /etc/edgefarm/kubeadm-join.conf -v5
  fi
fi

