# Manually Disable THP on Nodes Where NuoDB will Run

These instructions cover how to disable THP manually on nodes where NuoDB will be run.

```bash
HOSTS=( cafebabe-node-0 cafebabe-node-1 cafebabe-node-2 )
for HOST in "${HOSTS[@]}"
do
    ssh -T $HOST << EOSSH
#!/usr/bin/env bash

sudo mkdir -p /etc/tuned/nuodb

sudo cat > /etc/tuned/nuodb/tuned.conf <<'EOF'
[main]
summary=Optimize for NuoDB Distributed RDBMS
include=openshift-node
[vm]
transparent_hugepages=never
EOF

sudo chmod 0644 /etc/tuned/nuodb/tuned.conf
sudo tuned-adm profile nuodb
sudo systemctl restart tuned
EOSSH
done
```

And to verify that transparent huge pages is disabled on all nodes, a simple heredoc bash script to run in your terminal:

```bash
HOSTS=( cafebabe-node-0 cafebabe-node-1 cafebabe-node-2 )
for HOST in "${HOSTS[@]}"
do
    ssh -T $HOST << EOSSH
cat /sys/kernel/mm/transparent_hugepage/enabled
EOSSH
done
```

The output will look something like:

```bash
always madvise [never]
always madvise [never]
always madvise [never]
```
