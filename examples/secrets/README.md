# How to handle secrets in EdgeFarm

[sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) is a
great solution to store kubernetes (the base of
EdgeFarm) secrets in source code management systems like git.

The big advantage is that sealed-secrets can be added to EdgeFarm in
any way, either directly using kubectl or via a CI/CD system, but none
of the systems need a key to the secrets (as in other solutions such as
[SOPS](https://github.com/mozilla/sops)).
Only the cluster owns the key and converts the sealed-secret back into
a kubernetes secret, usable by your applications and processing pipelines.

That is why sealed-secrets is installed by default with EdgeFarm.

In the following steps, a "Sealed Secret" is created from a kubernetes
Secret. This sealed secret contains only encrypted information and can be
checked in without any problems. Only the sealed-secret operator
installed in EdgeFarm has the necessary key to convert the sealed secret
back into a Kubernetes secret. The actual kubernetes secret is not
published anywhere.

It is important to understand that data encrypted using sealed-secrets can
only be decrypted by the EdgeFarm instance that performed the encryption.
It is not possible to move sealed-secrets back and forth between different
EdgeFarm instances.

sealed-secrets does not serve as a "root of trust". sealed-secrets is
rather an exta layer of securitiy and simplifies the deployment process.
Tools such as [SOPS](https://github.com/mozilla/sops) are suitable for
storing also the "raw" secrets in a source code management system.

## Prerequisites

* Your current kube-context must point to the EdgeFarm cluster, where you
  want to store the secrets in.
* actual version of
  [kubeseal](https://github.com/bitnami-labs/sealed-secrets/releases)
  needs to be installed into your system.

## Step 1: Create a kubernetes secret

An initial secret can be created in many ways. A very convenient one is
kustomize, which is used in the following example.

If you want to create your secret from several key-value pairs, placed in
one file, see `secrets.env`.

```bash
$ cat secrets.env
key1=SecretValue1
key2=SecretValue2
key3=SecretValue3
key4=SecretValue4
key5=SecretValue5
```

If you want to create your secrets from different files containing the
whole secret information's, see `secret.file`

```bash
$ cat secret.file
my super secret
multiline
data
```

The two types can also be combined, as illustrated in this example.

Now let's generate the secrets. To do this

To do this, kustomize is executed in the current folder.
First let's look at the definition in `kustomization.yaml`.

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: example-namespace-name

secretGenerator:
- name: example-secret
  type: Opaque
  env: ./secrets.env
  files:
  - ./secret.file
generatorOptions:
 disableNameSuffixHash: true
```

This controls the generation of the secret. It is important that the
correct namespace is selected in which the secret is to exist later.

The secretGenerator creates a kubernetes secret with the name
example-secret and uses the `secrets.env` as well as the `secret.file`
for the creation.

To generate the secret, the following command is executed:

```bash
$kubectl kustomize .
apiVersion: v1
data:
  key1: U2VjcmV0VmFsdWUx
  key2: U2VjcmV0VmFsdWUy
  key3: U2VjcmV0VmFsdWUz
  key4: U2VjcmV0VmFsdWU0
  key5: U2VjcmV0VmFsdWU1
  secret.file: bXkgc3VwZXIgc2VjcmV0Cm11bHRpbGluZQpkYXRh
kind: Secret
metadata:
  name: example-secret
  namespace: example-namespace-name
type: Opaque
```

The secret is only output and not stored anywhere.
This will be done in the next step.

As you can see, all keys under data are taken over into the secret.
additionally another secret, named after the filename, is created. All
data are base64 encoded and therefore decodable by everyone (NOT secure).

## Step 2: Convert secret to sealed secret

To convert the secret to a sealed secret, combine the generation command
with kubeseal and store the results to `example-sealed-secrets.yaml`:

```bash
kubectl kustomize . |  kubeseal --format yaml > example-sealed-secret.yaml
```

the resulting file looks like this:

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: example-secret
  namespace: example-namespace-name
spec:
  encryptedData:
    key1: AgCYsFYdAvfp3VekUjMRaLYZDXybjqwUDAtSbSfUhSkE...
    key2: ...
    key3: ...
    key4: ...
    key5: ...
    secret.file: ...
  template:
    data: null
    metadata:
      creationTimestamp: null
      name: example-secret
      namespace: example-namespace-name
    type: Opaque
```

Now, the data are encrypted by sealed-secret and can be stored and used
everywhere you want.

If you apply the sealed secret to your EdgeFarm instace, a secret will be
created in the given namespace and can be used by your workloads, like
applications or worker.

## Step 3: Improvement: Prepare kustomize to use the sealed secret
