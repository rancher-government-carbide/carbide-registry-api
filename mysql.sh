#!/bin/sh
# dependencies: kubectl, helm
NAMESPACE="${NAMESPACE:-carbide-api}"
kubectl create ns "$NAMESPACE"
kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-crds.yaml && \
kubectl apply -f https://raw.githubusercontent.com/mysql/mysql-operator/trunk/deploy/deploy-operator.yaml

kubectl create secret generic mypwds \
        --from-literal=rootUser=root \
        --from-literal=rootHost=% \
        --from-literal=rootPassword="carbidecarbidecarbide" -n "$NAMESPACE"

kubectl apply -n "$NAMESPACE" -f - <<EOF
apiVersion: mysql.oracle.com/v2
kind: InnoDBCluster
metadata:
  name: mycluster
spec:
  secretName: mypwds
  tlsUseSelfSigned: true
  instances: 3
  router:
    instances: 1
EOF
