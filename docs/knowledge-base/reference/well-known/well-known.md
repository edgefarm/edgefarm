# Well-Known Labels, Annotations and Taints

## Labels

### openyurt.io/is-edge-worker

Type: Label<br>
Example: `openyurt.io/is-edge-worker=true`<br>
Used on: nodes<br>

If this label is set to true on a node, it will be considered an Edge Node. This will enable special features like node autonomy.

### openyurt.io/node-pool-type

Type: Label<br>
Example: `openyurt.io/node-pool-type=edge`<br>
Used on: nodepools<br>

If this label is set to `edge` on a nodepool, it will be considered an Edge Nodepool. The value can also be set to `cloud` as the nodepool controller also supports cloud nodepools.

### apps.openyurt.io/desired-nodepool
Type: Label<br>
Example: `apps.openyurt.io/desired-nodepool=edge0`<br>
Used on: nodes<br>

If this label is set on a node the nodepool controller will try to move the node to the nodepool specified in the label. This is part of the mapping between nodes and nodepools.


### apps.openyurt.io/nodepool
Type: Label<br>
Example: `apps.openyurt.io/nodepool=edge0`<br>
Used on: nodes<br>

If this label is set on a node it means that the mapping to the corresponding nodepool is active.

