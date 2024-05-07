# regnode

Register existing bare-metal nodes in Rancher clusters that cannot be managed by Rancher provisioning providers (AWS EKS, Azure AKS, Google GKE etc.). When you create Custom Rancher clusters like k3s you need to register nodes manually by retrieving the register command from Rancher UI or HTTP API. 

### Regnode vs Rancher CLI vs Rancher UI

If possible it's recommended to use either Rancher UI or [rancher-cli](https://ranchermanager.docs.rancher.com/reference-guides/cli-with-rancher/rancher-cli) tool to manage clusters. nodes and other objects within them.

If you have a Rancher server hosted in a datacenter that's far from the downstream cluster, it's recommended to manage the cluster using rancher-cli tool to avoid latency.

To retrieve the command to add the node using `rancher-cli` execute:

```shell
rancher login <SERVER_URL> --token <BEARER_TOKEN>

rancher clusters add-node <CLUSTER> --etcd --controlplane --worker
```

If your cluster was created as a Custom cluster it's possible you will get this error:

```
FATA[0000] a node can only be manually registered to a custom cluster
```

If that's the case then `Regnode` comes to the rescue. 

### Installation

Regnode is intended to be used only once on a node that is not registered to any cluster yet. However it's possible to run regnode on an already registered node and register it to a different cluster or to the same one but with different parameters.

Examples:

- Complie from source and run.

- Install as a linux service.

- Run using ansible.

### Configuration 
**regdnode** reads `$HOME/.regnode` [YAML](https://quickref.me/yaml) file by deafult. If not found you need to configure parameters with ENV variables or flags and arguments.

Loading order (next overrrides the previous):

1. Config file (default: `$HOME/.regnode` or `regnode [-c|--config] <FILE_PATH>`)
2. ENV variables 
3. CLI flags and arguments (see `regnode -h` or `regnode [command] -h`)

Example configurations [here](https://google.com).

*NOTE: If the same parameter is configured in config, ENV var and CLI flag or argument the actual parameter value will be what is passed to the CLI as it's loaded as last.*

### CLI and Environment variables

To avoid clashing with other ENV variables `regnode` reads only ENV vars prefixed with `REGNODE_`.

#### Global variables

They are global so you can pass them both beore and after sub-command `regnode [flags] [command]` or `regnode [command] [flags]`

| Flag | Short | Env variable | Default | Description |
| --- | --- | --- | --- | --- |
| --config | -c | REGNODE_CONFIG | $HOME/.regnode | Changes path to config file |
| --debug | -d | REGNODE_DEBUG | false (Info level) | Makes logging more verbose |
| --api-url | -u | REGNODE_API_URL | n/a | Rancher server URL |
| --api-token | -t | REGNODE_API_TOKEN | n/a | Rancher [HTTP Bearer token](https://ranchermanager.docs.rancher.com/reference-guides/user-settings/api-keys) |

**NOTE**: *Currently to be able to retrieve a node registration command you need to create a token (API access key) with Scope: "none"*

#### regnode register [args] [flags]

See more `regnode register -h`

One of `--worker, --etcd, --controlplane` is required.

| Argument | Env variable | Required | Usage |
| --- | --- | --- | --- |
| name | REGNODE_CLUSTER_NAME | true | regnode register cluster-1 |

| Flag | Env variable | Default | Description |
| --- | --- | --- | --- |
| --worker | REGNODE_CLUSTER_WORKER | false |  Registers node as a worker |
| --etcd | REGNODE_CLUSTER_ETCD | false |  Registers node as etcd |
| --controlplane | REGNODE_CLUSTER_CONTROLPLANE | false | Registers node as a controlplane |

**NOTE:** *It's recommended to register node either as etcd and controlplane or just as a worker node. Your setup always depends on the case and the kind of custom cluster engine you use (RKE or k3s). To achieve HA you should refer to those engines documentations.*

### Thank you

