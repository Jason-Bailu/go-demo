source ~/.bash_profile

curl https://discovery.etcd.io/new?size=3
https://discovery.etcd.io/94d70ea9076e58e89b12ab094fe0aad6

# grab this token
TOKEN=token-01
CLUSTER_STATE=new
DISCOVERY=https://discovery.etcd.io/94d70ea9076e58e89b12ab094fe0aad6

etcd --data-dir=data.etcd --name n1 \
	--initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 \
	--advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://127.0.0.1:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

etcd --data-dir=data.etcd --name n2 \
	--initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 \
	--advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://127.0.0.1:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

etcd --data-dir=data.etcd --name n3 \
	--initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://127.0.0.1:2380 \
	--advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http:/127.0.0.1:2379 \
	--discovery ${DISCOVERY} \
	--initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}

export ETCDCTL_API=3
HOST_1=10.240.0.17
HOST_2=10.240.0.18
HOST_3=10.240.0.19
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379

etcdctl --endpoints=$ENDPOINTS member list