echo "k8s setup process start"
echo 'vagrant:ubuntu' | sudo chpasswd
echo "=========================="

echo "[Step1] install common packages"
sudo apt-get upgrade && apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common gnupg2

echo "[Step2] install docker"
# Add Docker's official GPG Key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
# Add the Docker apt repository:
add-apt-repository \
  "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable"
# Install Docker CE
apt-get update && apt-get install -y \
  containerd.io=1.2.13-2 \
  docker-ce=5:19.03.11~3-0~ubuntu-$(lsb_release -cs) \
  docker-ce-cli=5:19.03.11~3-0~ubuntu-$(lsb_release -cs)
# Set up the Docker daemon
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF
mkdir -p /etc/systemd/system/docker.service.d
# Restart Docker
systemctl daemon-reload
systemctl restart docker
sudo usermod -aG docker $USER
su - vagrant 
docker ps

echo "[Step3] Stop and Disable firewalld"
systemctl disable firewalld >/dev/null 2>&1
systemctl stop firewalld

echo "[Step4] Add sysctl settings"
cat >>/etc/sysctl.d/kubernetes.conf<<EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
sysctl --system >/dev/null 2>&1

echo "[Step5] Disable and turn off SWAP"
sed -i '/swap/d' /etc/fstab
swapoff -a

echo "[Step6] install kubeadm, kubelet and kubectl"
sudo apt-get update && sudo apt-get install -y apt-transport-https curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet
