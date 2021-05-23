echo "k8s install process"
echo "=========================="
echo $1
echo $2

echo "[Step 1] init k8s control-panel"
kubeadm init --apiserver-advertise-address=$1 --apiserver-cert-extra-sans=$1  --node-name k8s-master --pod-network-cidr=$2

echo "[Step 2] move admin config to current user"
mkdir -p /home/vagrant/.kube
cp -i /etc/kubernetes/admin.conf /home/vagrant/.kube/config
chown vagrant:vagrant /home/vagrant/.kube/config

echo "[Step 3] install clico"
export KUBECONFIG=/home/vagrant/.kube/config
# su - vagrant -c "kubectl create -f https://docs.projectcalico.org/v3.14/manifests/calico.yaml"
su - vagrant -c "kubectl create -f /home/vagrant/share/resource/calico.yaml"

echo "[Step 4] expose config to local"
sudo rm /home/vagrant/share/config/admin.conf
cp -i /etc/kubernetes/admin.conf /home/vagrant/share/config

echo "[Step 5] expose join master command"
sudo rm /home/vagrant/share/token/token.txt
touch /home/vagrant/share/token/token.txt
sudo kubeadm token create --print-join-command > /home/vagrant/share/token/token.txt
cat /home/vagrant/share/token/token.txt

