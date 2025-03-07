
# kafkactl

A command-line interface for interaction with Apache Kafka

[![Build Status](https://github.com/deviceinsight/kafkactl/workflows/Lint%20%2F%20Test%20%2F%20IT/badge.svg?branch=master)](https://github.com/deviceinsight/kafkactl/actions)
| [![command docs](https://img.shields.io/badge/command-docs-blue.svg)](https://deviceinsight.github.io/kafkactl/)  

## Features

- command auto-completion for bash, zsh, fish shell including dynamic completion for e.g. topics or consumer groups.
- support for avro schemas
- Configuration of different contexts
- directly access kafka clusters inside your kubernetes cluster

[![asciicast](https://asciinema.org/a/vmxrTA0h8CAXPnJnSFk5uHKzr.svg)](https://asciinema.org/a/vmxrTA0h8CAXPnJnSFk5uHKzr)

## Installation

You can install the pre-compiled binary or compile from source.

### Install the pre-compiled binary

**snap**:

```bash
snap install kafkactl
```

**homebrew**:
```bash
# install tap repostory once
brew tap deviceinsight/packages
# install kafkactl
brew install deviceinsight/packages/kafkactl
# upgrade kafkactl
brew upgrade deviceinsight/packages/kafkactl
```

**deb/rpm**:

Download the .deb or .rpm from the [releases page](https://github.com/deviceinsight/kafkactl/releases) and install with dpkg -i and rpm -i respectively.

**yay (AUR)**

There's a kafkactl [AUR package](https://aur.archlinux.org/packages/kafkactl/) available for Arch. Install it with your AUR helper of choice (e.g. [yay](https://github.com/Jguer/yay)):

```bash
yay -S kafkactl
```

**manually**:

Download the pre-compiled binaries from the [releases page](https://github.com/deviceinsight/kafkactl/releases) and copy to the desired location.

### Compiling from source

```bash
go get -u github.com/deviceinsight/kafkactl
```

**NOTE:** make sure that `kafkactl` is on PATH otherwise auto-completion won't work.

## Configuration

If no config file is found, a default config is generated in `$HOME/.config/kafkactl/config.yml`.
This configuration is suitable to get started with a single node cluster on a local machine. 

### Create a config file

Create `$HOME/.config/kafkactl/config.yml` with a definition of contexts that should be available

```yaml
contexts:
  default:
    brokers:
    - localhost:9092
  remote-cluster:
    brokers:
    - remote-cluster001:9092
    - remote-cluster002:9092
    - remote-cluster003:9092

    # optional: tls config
    tls:
      enabled: true
      ca: my-ca
      cert: my-cert
      certKey: my-key
      # set insecure to true to ignore all tls verification (defaults to false)
      insecure: false

    # optional: sasl support
    sasl:
      enabled: true
      username: admin
      password: admin
      # optional configure sasl mechanism as plaintext, scram-sha256, scram-sha512 (defaults to plaintext)
      mechanism: scram-sha512
  
    # optional: access clusters running kubernetes
    kubernetes:
      enabled: false
      binary: kubectl #optional
      kubeConfig: ~/.kube/config #optional
      kubeContext: my-cluster
      namespace: my-namespace

    # optional: clientID config (defaults to kafkactl-{username})
    clientID: my-client-id
    
    # optional: kafkaVersion (defaults to 2.0.0)
    kafkaVersion: 1.1.1

    # optional: timeout for admin requests (defaults to 3s)
    requestTimeout: 10s

    # optional: avro schema registry
    avro:
      schemaRegistry: localhost:8081
    
    # optional: changes the default partitioner
    defaultPartitioner: "hash"


current-context: default
```

The config file location is resolved by
 * checking for a provided commandline argument: `--config-file=$PATH_TO_CONFIG`
 * or by evaluating the environment variable: `export KAFKA_CTL_CONFIG=$PATH_TO_CONFIG`
 * or as default the config file is looked up from one of the following locations:
   * `$HOME/.config/kafkactl/config.yml`
   * `$HOME/.kafkactl/config.yml`
   * `$SNAP_REAL_HOME/.kafkactl/config.yml`
   * `$SNAP_DATA/kafkactl/config.yml`
   * `/etc/kafkactl/config.yml`

### Auto completion

#### bash

**NOTE:** if you installed via snap, bash completion should work automatically.

```
source <(kafkactl completion bash)
```

To load completions for each session, execute once:
Linux:
```
kafkactl completion bash > /etc/bash_completion.d/kafkactl
```
 
MacOS:
```
kafkactl completion bash > /usr/local/etc/bash_completion.d/kafkactl
```

#### zsh

```
source <(kafkactl completion zsh)
```

To load completions for each session, execute once:
```
kafkactl completion zsh > "${fpath[1]}/_kafkactl"
```

#### Fish

```
kafkactl completion fish | source
```

To load completions for each session, execute once:
```
kafkactl completion fish > ~/.config/fish/completions/kafkactl.fish
```

## Running in docker

Assuming your Kafka broker is accessible as `kafka:9092`, you can list topics by running: 

```bash
docker run --env BROKERS=kafka:9092 deviceinsight/kafkactl:latest get topics
```

If a more elaborate config is needed, you can mount it as a volume:

```bash
docker run -v /absolute/path/to/config.yml:/etc/kafkactl/config.yml deviceinsight/kafkactl get topics
``` 

## Configuration via environment variables

Every key in the `config.yml` can be overwritten via environment variables. The corresponding environment variable
for a key can be found by applying the following rules:

1. replace `.` by `_`
1. replace `-` by `_`
1. write the key name in ALL CAPS

e.g. the key `contexts.default.tls.certKey` has the corresponding environment variable `CONTEXTS_DEFAULT_TLS_CERTKEY`.

If environment variables for the `default` context should be set, the prefix `CONTEXTS_DEFAULT_` can be omitted.
So, instead of `CONTEXTS_DEFAULT_TLS_CERTKEY` one can also set `TLS_CERTKEY`.
See **root_test.go** for more examples.

## Running in Kubernetes

> :construction: This feature is still experimental.

If your kafka cluster is not directly accessible from your machine, but it is accessible from a kubernetes cluster
which in turn is accessible via `kubectl` from your machine you can configure kubernetes support:

```$yaml
contexts:
  kafka-cluster:
    brokers:
      - broker1:9092
      - broker2:9092
    kubernetes:
      enabled: true
      binary: kubectl #optional
      kubeContext: k8s-cluster
      namespace: k8s-namespace
```

Instead of directly talking to kafka brokers a kafkactl docker image is deployed as a pod into the kubernetes
cluster, and the defined namespace. Standard-Input and Standard-Output are then wired between the pod and your shell
running kafkactl. 

There are two options:
1. You can run `kafkactl attach` with your kubernetes cluster configured. This will use `kubectl run` to create a pod
in the configured kubeContext/namespace which runs an image of kafkactl and gives you a `bash` into the container.
Standard-in is piped to the pod and standard-out, standard-err directly to your shell. You even get auto-completion.

2. You can run any other kafkactl command with your kubernetes cluster configured. Instead of directly
querying the cluster a pod is deployed, and input/output are wired between pod and your shell.

The names of the brokers have to match the service names used to access kafka in your cluster. A command like this should
 give you this information:
```bash
kubectl get svc | grep kafka
```

> :bulb: The first option takes a bit longer to start up since an Ubuntu based docker image is used in order to have
a bash available. The second option uses a docker image build from scratch and should therefore be quicker.
Which option is more suitable, will depend on your use-case. 

> :warning: when _kafkactl_ was installed via _snap_ make sure to configure the absolute path to your `kubectl` binary. 
Snaps run with a different $PATH and therefore are unable to access binaries on $PATH. 

## Command documentation

The documentation for all available commands can be found here:

[![command docs](https://img.shields.io/badge/command-docs-blue.svg)](https://deviceinsight.github.io/kafkactl/)


## Examples

### Consuming messages

Consuming messages from a topic can be done with:
```bash
kafkactl consume my-topic
```

In order to consume starting from the oldest offset use:
```bash
kafkactl consume my-topic --from-beginning
```

The following example prints message `key` and `timestamp` as well as `partition` and `offset` in `yaml` format:
```bash
kafkactl consume my-topic --print-keys --print-timestamps -o yaml
```

Headers of kafka messages can be printed with the parameter `--print-headers` e.g.:
```bash
kafkactl consume my-topic --print-headers -o yaml
```

If one is only interested in the last `n` messages this can be achieved by `--tail` e.g.:
```bash
kafkactl consume my-topic --tail=5
```

The consumer can be stopped when the latest offset is reached using `--exit` parameter e.g.:
```bash
kafkactl consume my-topic --from-beginning --exit
```

The following example prints keys in hex and values in base64:
```bash
kafkactl consume my-topic --print-keys --key-encoding=hex --value-encoding=base64
```

### Producing messages

Producing messages can be done in multiple ways. If we want to produce a message with `key='my-key'`,
`value='my-value'` to the topic `my-topic` this can be achieved with one of the following commands:

```bash
echo "my-key#my-value" | kafkactl produce my-topic --separator=#
echo "my-value" | kafkactl produce my-topic --key=my-key
kafkactl produce my-topic --key=my-key --value=my-value
```

If we have a file containing messages where each line contains `key` and `value` separated by `#`, the file can be
used as input to produce messages to topic `my-topic`:

```bash
cat myfile | kafkactl produce my-topic --separator=#
```

The same can be accomplished without piping the file to stdin with the `--file` parameter:
```bash
kafkactl produce my-topic --separator=# --file=myfile
```

If the messages in the input file need to be split by a different delimiter than `\n` a custom line separator can be provided:
 ```bash
 kafkactl produce my-topic --separator=# --lineSeparator=|| --file=myfile
 ```

**NOTE:** if the file was generated with `kafkactl consume --print-keys --print-timestamps my-topic` the produce
command is able to detect the message timestamp in the input and will ignore it. 

the number of messages produced per second can be controlled with the `--rate` parameter:

```bash
cat myfile | kafkactl produce my-topic --separator=# --rate=200
```

It is also possible to specify the partition to insert the message:
```bash
kafkactl produce my-topic --key=my-key --value=my-value --partition=2
```

Additionally, a different partitioning scheme can be used. When a `key` is provided the default partitioner
uses the `hash` of the `key` to assign a partition. So the same `key` will end up in the same partition: 
```bash
# the following 3 messages will all be inserted to the same partition
kafkactl produce my-topic --key=my-key --value=my-value
kafkactl produce my-topic --key=my-key --value=my-value
kafkactl produce my-topic --key=my-key --value=my-value

# the following 3 messages will probably be inserted to different partitions
kafkactl produce my-topic --key=my-key --value=my-value --partitioner=random
kafkactl produce my-topic --key=my-key --value=my-value --partitioner=random
kafkactl produce my-topic --key=my-key --value=my-value --partitioner=random
```

Message headers can also be written:
```bash
kafkactl produce my-topic --key=my-key --value=my-value --header key1:value1 --header key2:value\:2
```

The following example writes the key from base64 and value from hex:
```bash
kafkactl produce my-topic --key=dGVzdC1rZXk= --key-encoding=base64 --value=0000000000000000 --value-encoding=hex
```

### Avro support

In order to enable avro support you just have to add the schema registry to your configuration:
```$yaml
contexts:
  localhost:
    avro:
      schemaRegistry: localhost:8081
```

#### Producing to an avro topic

`kafkactl` will lookup the topic in the schema registry in order to determine if key or value needs to be avro encoded.
If producing with the latest `schemaVersion` is sufficient, no additional configuration is needed an `kafkactl` handles
this automatically.

If however one needs to produce an older `schemaVersion` this can be achieved by providing the parameters `keySchemaVersion`, `valueSchemaVersion`.

##### Example

```bash
# create a topic
kafkactl create topic avro_topic
# add a schema for the topic value
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \
--data '{"schema": "{\"type\": \"record\", \"name\": \"LongList\", \"fields\" : [{\"name\": \"next\", \"type\": [\"null\", \"LongList\"], \"default\": null}]}"}' \
http://localhost:8081/subjects/avro_topic-value/versions
# produce a message
kafkactl produce avro_topic --value {\"next\":{\"LongList\":{}}}
# consume the message
kafkactl consume avro_topic --from-beginning --print-schema -o yaml
```

#### Consuming from an avro topic

As for producing `kafkactl` will also lookup the topic in the schema registry to determine if key or value needs to be
decoded with an avro schema.

The `consume` command handles this automatically and no configuration is needed.

An additional parameter `print-schema` can be provided to display the schema used for decoding.

### Altering topics

Using the `alter topic` command allows you to change the partition count, replication factor and topic-level
configurations of an existing topic.

The partition count can be increased with:
```bash
kafkactl alter topic my-topic --partitions 32
```

The replication factor can be altered with:
```bash
kafkactl alter topic my-topic --replication-factor 2
```

> :information_source: when altering replication factor, kafkactl tries to keep the number of replicas assigned to each
> broker balanced. If you need more control over the assigned replicas use `alter partition` directly.

The topic configs can be edited by supplying key value pairs as follows:
```bash
kafkactl alter topic my-topic --config retention.ms=3600000 --config cleanup.policy=compact
```

> :bulb: use the flag `--validate-only` to perform a dry-run without actually modifying the topic 

### Altering partitions

The assigned replicas of a partition can directly be altered with:
```bash
# set brokers 102,103 as replicas for partition 3 of topic my-topic
kafkactl alter topic my-topic 3 -r 102,103
```

### Consumer groups

In order to get a list of consumer groups the `get consumer-groups` command can be used:
```bash
# all available consumer groups
kafkactl get consumer-groups 
# only consumer groups for a single topic
kafkactl get consumer-groups --topic my-topic
# using command alias
kafkactl get cg
```

To get detailed information about the consumer group use `describe consumer-group`. If the parameter `--partitions`
is provided details will be printed for each partition otherwise the partitions are aggregated to the clients.

```bash
# describe a consumer group
kafkactl describe consumer-group my-group 
# show partition details only for partitions with lag
kafkactl describe consumer-group my-group --only-with-lag
# show details only for a single topic
kafkactl describe consumer-group my-group --topic my-topic
# using command alias
kafkactl describe cg my-group
```

### Reset consumer group offsets

in order to ensure the reset does what it is expected, per default only
the results are printed without actually executing it. Use the additional parameter `--execute` to perform the reset. 

```bash
# reset offset of for all partitions to oldest offset
kafkactl reset offset my-group --topic my-topic --oldest
# reset offset of for all partitions to newest offset
kafkactl reset offset my-group --topic my-topic --newest
# reset offset for a single partition to specific offset
kafkactl reset offset my-group --topic my-topic --partition 5 --offset 100
```

### ACL Management

Available ACL operations are documented [here](https://docs.confluent.io/platform/current/kafka/authorization.html#operations).

#### Create a new ACL

```bash
# create an acl that allows topic read for a user 'consumer'
kafkactl create acl --topic my-topic --operation read --principal User:consumer --allow
# create an acl that denies topic write for a user 'consumer' coming from a specific host
kafkactl create acl --topic my-topic --operation write --host 1.2.3.4 --principal User:consumer --deny
# allow multiple operations
kafkactl create acl --topic my-topic --operation read --operation describe --principal User:consumer --allow
# allow on all topics with prefix common prefix
kafkactl create acl --topic my-prefix --pattern prefixed --operation read --principal User:consumer --allow
```

#### List ACLs

```bash
# list all acl
kafkactl get acl
# list all acl (alias command)
kafkactl get access-control-list
# filter only topic resources
kafkactl get acl --topics
# filter only consumer group resources with operation read
kafkactl get acl --groups --operation read
```

#### Delete ACLs

```bash
# delete all topic read acls
kafkactl delete acl --topics --operation read --pattern any
# delete all topic acls for any operation
kafkactl delete acl --topics --operation any --pattern any
# delete all cluster acls for any operation
kafkactl delete acl --cluster --operation any --pattern any
# delete all consumer-group acls with operation describe, patternType prefixed and permissionType allow
kafkactl delete acl --groups --operation describe --pattern prefixed --allow
```

### Getting Brokers

To get the list of brokers of a kafka cluster use `get brokers`

```bash
# get the list of brokers
kafkactl get brokers
```
