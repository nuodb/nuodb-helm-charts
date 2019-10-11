# Using hostPath Storage with DaemonSet

## Set Up Local Storage

Nuo CE Helm templates require local storage, so set it up; again we're running a bash heredoc script for each target host (change the host name list and run the script in your terminal):

```bash
HOSTS=( cafebabe-node-0 cafebabe-node-1 cafebabe-node-2 )
for HOST in "${HOSTS[@]}"
do
    ssh -T $HOST << EOSSH
sudo mkdir -p /mnt/nuodb/data/disk0
sudo mkdir -p /mnt/nuodb/data/disk1
sudo mkdir -p /mnt/nuodb/data/disk2
sudo mkdir -p /mnt/nuodb/data/disk3

sudo chmod -R 775 /mnt/nuodb/data

# RHEL 7.5, and earlier
# sudo chcon -R unconfined_u:object_r:svirt_sandbox_file_t:s0 /mnt/nuodb/data
# RHEL 7.6, or later
sudo chcon -R system_u:object_r:container_file_t:s0 /mnt/nuodb/data
sudo chown -R root:root /mnt/nuodb/data

sudo ls -Zal /mnt/nuodb/data/
EOSSH
done
```

If you need to clean up after yourself:

```bash
HOSTS=( cafebabe-node-0 cafebabe-node-1 cafebabe-node-2 )
for HOST in "${HOSTS[@]}"
do
    ssh -T $HOST << EOSSH
sudo rm -fr /mnt/nuodb/data/disk0/*
sudo rm -fr /mnt/nuodb/data/disk1/*
sudo rm -fr /mnt/nuodb/data/disk2/*
sudo rm -fr /mnt/nuodb/data/disk3/*
EOSSH
done
```

## Set Up Security Policies

Permit containers to run as any UID, for the default service account, and permit those containers to access host provisioned storage:

```bash
oc adm policy add-scc-to-user hostmount-anyuid -z default
oc adm policy add-scc-to-user anyuid -z default
```

## Set Up Labels

```bash
HOSTS=( cafebabe-node-0 cafebabe-node-1 cafebabe-node-2 )
for HOST in "${HOSTS[@]}"
do
    oc label node ${HOST} nuodb.com/zone=east --overwrite=true
done
oc label node cafebabe-node-0 nuodb.com/node-type=storage
```

Run the last command, above, multiple times for NuoDB Enterprise Edition,
once for each host where an SM will be permitted to run.

Pick one host to start and SM and label it:

```bash
oc label node <node> nuodb.com/nuodb.demo=nobackup
```

To check your list of labels, issue the following command:

```bash
$ oc get nodes -L nuodb.com/zone -L nuodb.com/node-type \
  -L nuodb.com/domain -L nuodb.com/nuodb.demo \
  -L failure-domain.beta.kubernetes.io/region \
  -L failure-domain.beta.kubernetes.io/zone
NAME                STATUS                     ROLES     AGE       VERSION             ZONE      NODE-TYPE   DOMAIN    NUODB.DEMO
cafebabe-infra-0    Ready                      <none>    7d        v1.9.1+a0ce1bc657                                   
cafebabe-infra-1    Ready                      <none>    7d        v1.9.1+a0ce1bc657                                   
cafebabe-infra-2    Ready                      <none>    7d        v1.9.1+a0ce1bc657                                   
cafebabe-master-0   Ready,SchedulingDisabled   master    7d        v1.9.1+a0ce1bc657                                   
cafebabe-master-1   Ready,SchedulingDisabled   master    7d        v1.9.1+a0ce1bc657                                   
cafebabe-master-2   Ready,SchedulingDisabled   master    7d        v1.9.1+a0ce1bc657                                   
cafebabe-node-0     Ready                      compute   7d        v1.9.1+a0ce1bc657   east      storage               
cafebabe-node-1     Ready                      compute   7d        v1.9.1+a0ce1bc657   east                            
cafebabe-node-2     Ready                      compute   7d        v1.9.1+a0ce1bc657   east                            
```
