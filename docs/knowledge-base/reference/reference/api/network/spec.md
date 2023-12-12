---
hide:
- toc
---
# API Reference

Packages:

- [streams.network.edgefarm.io/v1alpha1](#streamsnetworkedgefarmiov1alpha1)

# streams.network.edgefarm.io/v1alpha1

Resource Types:

- [Network](#network)




## Network
<sup><sup>[↩ Parent](#streamsnetworkedgefarmiov1alpha1 )</sup></sup>








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>optional
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>streams.network.edgefarm.io/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Network</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#networkspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#networkstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec
<sup><sup>[↩ Parent](#network)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparameters">parameters</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#networkspecclaimref">claimRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspeccompositionref">compositionRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspeccompositionrevisionref">compositionRevisionRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspeccompositionrevisionselector">compositionRevisionSelector</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspeccompositionselector">compositionSelector</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>compositionUpdatePolicy</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: Automatic, Manual<br/>
            <i>Default</i>: Automatic<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecenvironmentconfigrefsindex">environmentConfigRefs</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecpublishconnectiondetailsto">publishConnectionDetailsTo</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecresourcerefsindex">resourceRefs</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecwriteconnectionsecrettoref">writeConnectionSecretToRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparametersconsumersindex">consumers</a></b></td>
        <td>[]object</td>
        <td>
          List of consumers<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersnats">nats</a></b></td>
        <td>object</td>
        <td>
          NATS config<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersresourceconfig">resourceConfig</a></b></td>
        <td>object</td>
        <td>
          Defines general properties for this resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindex">streams</a></b></td>
        <td>[]object</td>
        <td>
          List of streams<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparameterssubnetworksindex">subNetworks</a></b></td>
        <td>[]object</td>
        <td>
          The subnetworks that are part of this network
<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersusersindex">users</a></b></td>
        <td>[]object</td>
        <td>
          List of users to create<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.consumers[index]
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>



Configuration for a consumer

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparametersconsumersindexconfig">config</a></b></td>
        <td>object</td>
        <td>
          Config is the consumer configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the consumer<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>streamRef</b></td>
        <td>string</td>
        <td>
          The name of the stream the consumer is created for<br/>
          <br/>
            <i>Default</i>: main<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.consumers[index].config
<sup><sup>[↩ Parent](#networkspecparametersconsumersindex)</sup></sup>



Config is the consumer configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>ackPolicy</b></td>
        <td>enum</td>
        <td>
          AckPolicy describes the requirement of client acknowledgements, either Explicit, None, or All. For more information see https://docs.nats.io/nats-concepts/jetstream/consumers#ackpolicy<br/>
          <br/>
            <i>Enum</i>: Explicit, None, All<br/>
            <i>Default</i>: Explicit<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>ackWait</b></td>
        <td>string</td>
        <td>
          AckWait is the duration that the server will wait for an ack for any individual message once it has been delivered to a consumer. If an ack is not received in time, the message will be redelivered. Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.<br/>
          <br/>
            <i>Default</i>: 30s<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>deliverPolicy</b></td>
        <td>enum</td>
        <td>
          DeliverPolicy defines the point in the stream to receive messages from, either All, Last, New, ByStartSequence, ByStartTime, or LastPerSubject. Fore more information see https://docs.nats.io/jetstream/concepts/consumers#deliverpolicy<br/>
          <br/>
            <i>Enum</i>: All, Last, New, ByStartSequence, ByStartTime, LastPerSubject<br/>
            <i>Default</i>: All<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>numReplicas</b></td>
        <td>integer</td>
        <td>
          Replicas sets the number of replicas for the consumer's state. By default, when the value is set to zero, consumers inherit the number of replicas from the stream.<br/>
          <br/>
            <i>Default</i>: 0<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>replayPolicy</b></td>
        <td>enum</td>
        <td>
          ReplayPolicy is used to define the mode of message replay. If the policy is Instant, the messages will be pushed to the client as fast as possible while adhering to the Ack Policy, Max Ack Pending and the client's ability to consume those messages. If the policy is Original, the messages in the stream will be pushed to the client at the same rate that they were originally received, simulating the original timing of messages.<br/>
          <br/>
            <i>Enum</i>: Instant, Original<br/>
            <i>Default</i>: Instant<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>backoff</b></td>
        <td>string</td>
        <td>
          Backoff is a list of time durations that represent the time to delay based on delivery count. Format of the durations is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s where multiple durations are separated by commas. Example: `1s,2s,3s,4s,5s`.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Description is a human readable description of the consumer. This can be particularly useful for ephemeral consumers to indicate their purpose since the durable name cannot be provided.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>filterSubject</b></td>
        <td>string</td>
        <td>
          FilterSubject defines an overlapping subject with the subjects bound to the stream which will filter the set of messages received by the consumer.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactiveThreshold</b></td>
        <td>string</td>
        <td>
          InactiveThreshold defines the duration that instructs the server to cleanup consumers that are inactive for that long. Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxAckPending</b></td>
        <td>integer</td>
        <td>
          MaxAckPending sets the number of outstanding acks that are allowed before message delivery is halted.<br/>
          <br/>
            <i>Default</i>: 1000<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxDeliver</b></td>
        <td>integer</td>
        <td>
          MaxDeliver is the maximum number of times a specific message delivery will be attempted. Applies to any message that is re-sent due to ack policy (i.e. due to a negative ack, or no ack sent by the client).<br/>
          <br/>
            <i>Default</i>: -1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>memStorage</b></td>
        <td>boolean</td>
        <td>
          MemoryStorage if set, forces the consumer state to be kept in memory rather than inherit the storage type of the stream (file in this case).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optStartSeq</b></td>
        <td>integer</td>
        <td>
          OptStartSeq is an optional start sequence number and is used with the DeliverByStartSequence deliver policy.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optStartTime</b></td>
        <td>string</td>
        <td>
          OptStartTime is an optional start time and is used with the DeliverByStartTime deliver policy. The time format is RFC 3339, e.g. 2023-01-09T14:48:32Z<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersconsumersindexconfigpull">pull</a></b></td>
        <td>object</td>
        <td>
          PullConsumer defines the pull-based consumer configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersconsumersindexconfigpush">push</a></b></td>
        <td>object</td>
        <td>
          PushConsumer defines the push-based consumer configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sampleFreq</b></td>
        <td>string</td>
        <td>
          SampleFrequency sets the percentage of acknowledgements that should be sampled for observability.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.consumers[index].config.pull
<sup><sup>[↩ Parent](#networkspecparametersconsumersindexconfig)</sup></sup>



PullConsumer defines the pull-based consumer configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>maxBatch</b></td>
        <td>integer</td>
        <td>
          MaxRequestBatch defines th maximum batch size a single pull request can make. When set with MaxRequestMaxBytes, the batch size will be constrained by whichever limit is hit first. This is a pull consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxBytes</b></td>
        <td>integer</td>
        <td>
          MaxRequestMaxBytes defines the  maximum total bytes that can be requested in a given batch. When set with MaxRequestBatch, the batch size will be constrained by whichever limit is hit first. This is a pull consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxExpires</b></td>
        <td>string</td>
        <td>
          MaxRequestExpires defines the maximum duration a single pull request will wait for messages to be available to pull. This is a pull consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxWaiting</b></td>
        <td>integer</td>
        <td>
          MaxWaiting defines the maximum number of waiting pull requests. This is a pull consumer specific setting.<br/>
          <br/>
            <i>Default</i>: 512<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.consumers[index].config.push
<sup><sup>[↩ Parent](#networkspecparametersconsumersindexconfig)</sup></sup>



PushConsumer defines the push-based consumer configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>deliverGroup</b></td>
        <td>string</td>
        <td>
          DeliverGroup defines the queue group name which, if specified, is then used to distribute the messages between the subscribers to the consumer. This is analogous to a queue group in core NATS. See https://docs.nats.io/nats-concepts/core-nats/queue for more information on queue groups. This is a push consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deliverSubject</b></td>
        <td>string</td>
        <td>
          DeliverSubject defines the subject to deliver messages to. Note, setting this field implicitly decides whether the consumer is push or pull-based. With a deliver subject, the server will push messages to client subscribed to this subject. This is a push consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flowControl</b></td>
        <td>boolean</td>
        <td>
          FlowControl enables per-subscription flow control using a sliding-window protocol. This protocol relies on the server and client exchanging messages to regulate when and how many messages are pushed to the client. This one-to-one flow control mechanism works in tandem with the one-to-many flow control imposed by MaxAckPending across all subscriptions bound to a consumer. This is a push consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>headersOnly</b></td>
        <td>boolean</td>
        <td>
          HeadersOnly delivers, if set, only the headers of messages in the stream and not the bodies. Additionally adds Nats-Msg-Size header to indicate the size of the removed payload.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>idleHeartbeat</b></td>
        <td>string</td>
        <td>
          IdleHeartbeat defines, if set, that the server will regularly send a status message to the client (i.e. when the period has elapsed) while there are no new messages to send. This lets the client know that the JetStream service is still up and running, even when there is no activity on the stream. The message status header will have a code of 100. Unlike FlowControl, it will have no reply to address. It may have a description such "Idle Heartbeat". Note that this heartbeat mechanism is all handled transparently by supported clients and does not need to be handled by the application. Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s. This is a push consumer specific setting.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rateLimitBps</b></td>
        <td>integer</td>
        <td>
          RateLimit is used to throttle the delivery of messages to the consumer, in bits per second.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.nats
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>



NATS config

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>address</b></td>
        <td>string</td>
        <td>
          The address of the NATS server
<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          The name of the operator the account is created for<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.parameters.resourceConfig
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>



Defines general properties for this resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>accountSecret</b></td>
        <td>string</td>
        <td>
          Name of secret containing the account information<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersresourceconfigkubernetes">kubernetes</a></b></td>
        <td>object</td>
        <td>
          Config for provider kubernetes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersresourceconfignatssecrets">natssecrets</a></b></td>
        <td>object</td>
        <td>
          Config for provider natssecrets<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>userSecret</b></td>
        <td>string</td>
        <td>
          Name of secret containing the user information<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.resourceConfig.kubernetes
<sup><sup>[↩ Parent](#networkspecparametersresourceconfig)</sup></sup>



Config for provider kubernetes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>providerConfigName</b></td>
        <td>string</td>
        <td>
          Name of provider config to use for kubernetes<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.resourceConfig.natssecrets
<sup><sup>[↩ Parent](#networkspecparametersresourceconfig)</sup></sup>



Config for provider natssecrets

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>providerConfigName</b></td>
        <td>string</td>
        <td>
          Name of provider config to use for natssecrets<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index]
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>



Configuration for a stream

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfig">config</a></b></td>
        <td>object</td>
        <td>
          Config is the stream configuration.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the stream<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reference</b></td>
        <td>string</td>
        <td>
          When type is mirror, the name of the stream to mirror<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>references</b></td>
        <td>[]string</td>
        <td>
          When type is aggregate, the names of the streams to aggregate<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subNetworkRef</b></td>
        <td>string</td>
        <td>
          The name of the sub network the stream is created for<br/>
          <br/>
            <i>Default</i>: main<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The type of the stream<br/>
          <br/>
            <i>Enum</i>: Standard, Aggregate, Mirror<br/>
            <i>Default</i>: Standard<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config
<sup><sup>[↩ Parent](#networkspecparametersstreamsindex)</sup></sup>



Config is the stream configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>discard</b></td>
        <td>enum</td>
        <td>
          Discard defines the behavior of discarding messages when any streams' limits have been reached. Old (default): This policy will delete the oldest messages in order to maintain the limit. For example, if MaxAge is set to one minute, the server will automatically delete messages older than one minute with this policy. New: This policy will reject new messages from being appended to the stream if it would exceed one of the limits. An extension to this policy is DiscardNewPerSubject which will apply this policy on a per-subject basis within the stream.<br/>
          <br/>
            <i>Enum</i>: Old, New<br/>
            <i>Default</i>: Old<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maxBytes</b></td>
        <td>integer</td>
        <td>
          MaxBytes defines how many bytes the Stream may contain. Adheres to Discard Policy, removing oldest or refusing new messages if the Stream exceeds this size.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: -1<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maxConsumers</b></td>
        <td>integer</td>
        <td>
          MaxConsumers defines how many Consumers can be defined for a given Stream. Define -1 for unlimited.<br/>
          <br/>
            <i>Default</i>: -1<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maxMsgs</b></td>
        <td>integer</td>
        <td>
          MaxMsgs defines how many messages may be in a Stream. Adheres to Discard Policy, removing oldest or refusing new messages if the Stream exceeds this number of messages.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: -1<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retention</b></td>
        <td>enum</td>
        <td>
          Retention defines the retention policy for the stream.<br/>
          <br/>
            <i>Enum</i>: Limits, Interest, WorkQueue<br/>
            <i>Default</i>: Limits<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>storage</b></td>
        <td>enum</td>
        <td>
          Storage defines the storage type for stream data..<br/>
          <br/>
            <i>Enum</i>: File, Memory<br/>
            <i>Default</i>: File<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>allowDirect</b></td>
        <td>boolean</td>
        <td>
          AllowDirect is a flag that if true and the stream has more than one replica, each replica will respond to direct get requests for individual messages, not only the leader.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowRollup</b></td>
        <td>boolean</td>
        <td>
          AllowRollup is a flag to allow the use of the Nats-Rollup header to replace all contents of a stream, or subject in a stream, with a single new message.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>denyDelete</b></td>
        <td>boolean</td>
        <td>
          DenyDelete is a flag to restrict the ability to delete messages from a stream via the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>denyPurge</b></td>
        <td>boolean</td>
        <td>
          DenyPurge is a flag to restrict the ability to purge messages from a stream via the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Description is a human readable description of the stream.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>discardNewPerSubject</b></td>
        <td>boolean</td>
        <td>
          DiscardOldPerSubject will discard old messages per subject.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>duplicates</b></td>
        <td>string</td>
        <td>
          Duplicates defines the time window within which to track duplicate messages.<br/>
          <br/>
            <i>Default</i>: 2m0s<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxAge</b></td>
        <td>string</td>
        <td>
          MaxAge is the maximum age of a message in the stream. Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.<br/>
          <br/>
            <i>Default</i>: 0s<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxMsgSize</b></td>
        <td>integer</td>
        <td>
          MaxBytesPerSubject defines the largest message that will be accepted by the Stream.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: -1<br/>
            <i>Minimum</i>: -1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxMsgsPerSubject</b></td>
        <td>integer</td>
        <td>
          MaxMsgsPerSubject defines the limits how many messages in the stream to retain per subject.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: -1<br/>
            <i>Minimum</i>: -1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigmirror">mirror</a></b></td>
        <td>object</td>
        <td>
          Mirror is the mirror configuration for the stream.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mirrorDirect</b></td>
        <td>boolean</td>
        <td>
          MirrorDirect is a flag that if true, and the stream is a mirror, the mirror will participate in a serving direct get requests for individual messages from origin stream.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>noAck</b></td>
        <td>boolean</td>
        <td>
          NoAck is a flag to disable acknowledging messages that are received by the Stream.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigplacement">placement</a></b></td>
        <td>object</td>
        <td>
          Placement is the placement policy for the stream.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigrepublish">rePublish</a></b></td>
        <td>object</td>
        <td>
          Allow republish of the message after being sequenced and stored.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>replicas</b></td>
        <td>integer</td>
        <td>
          Replicas defines how many replicas to keep for each message in a clustered JetStream.<br/>
          <br/>
            <i>Default</i>: 1<br/>
            <i>Minimum</i>: 1<br/>
            <i>Maximum</i>: 5<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sealed</b></td>
        <td>boolean</td>
        <td>
          Sealed is a flag to prevent message deletion from  the stream  via limits or API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigsourcesindex">sources</a></b></td>
        <td>[]object</td>
        <td>
          Sources is the list of one or more sources configurations for the stream.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subjects</b></td>
        <td>[]string</td>
        <td>
          Subjects is a list of subjects to consume, supports wildcards.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>template</b></td>
        <td>string</td>
        <td>
          Template is the owner of the template associated with this stream.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.mirror
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfig)</sup></sup>



Mirror is the mirror configuration for the stream.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the origin stream to source messages from.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>domain</b></td>
        <td>string</td>
        <td>
          Domain is the JetStream domain of where the origin stream exists. This is commonly used between a cluster/supercluster and a leaf node/cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigmirrorexternal">external</a></b></td>
        <td>object</td>
        <td>
          External is the external stream configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>filterSubject</b></td>
        <td>string</td>
        <td>
          FilterSubject is an optional filter subject which will include only messages that match the subject, typically including a wildcard.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startSeq</b></td>
        <td>integer</td>
        <td>
          StartSeq is an optional start sequence the of the origin stream to start mirroring from.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startTime</b></td>
        <td>string</td>
        <td>
          StartTime is an optional message start time to start mirroring from. Any messages that are equal to or greater than the start time will be included. The time format is RFC 3339, e.g. 2023-01-09T14:48:32Z<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.mirror.external
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfigmirror)</sup></sup>



External is the external stream configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiPrefix</b></td>
        <td>string</td>
        <td>
          APIPrefix is the prefix for the API of the external stream.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>deliverPrefix</b></td>
        <td>string</td>
        <td>
          DeliverPrefix is the prefix for the deliver subject of the external stream.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.placement
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfig)</sup></sup>



Placement is the placement policy for the stream.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cluster</b></td>
        <td>string</td>
        <td>
          Cluster is the name of the Jetstream cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          Tags defines a list of server tags.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.rePublish
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfig)</sup></sup>



Allow republish of the message after being sequenced and stored.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>destination</b></td>
        <td>string</td>
        <td>
          Destination is the destination subject messages will be re-published to. The source and destination must be a valid subject mapping. For information on subject mapping see https://docs.nats.io/jetstream/concepts/subjects#subject-mapping<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Source is an optional subject pattern which is a subset of the subjects bound to the stream. It defaults to all messages in the stream, e.g. >.<br/>
          <br/>
            <i>Default</i>: ><br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>headersOnly</b></td>
        <td>boolean</td>
        <td>
          HeadersOnly defines if true, that the message data will not be included in the re-published message, only an additional header Nats-Msg-Size indicating the size of the message in bytes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.sources[index]
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfig)</sup></sup>



StreamSource dictates how streams can source from other streams.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the origin stream to source messages from.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>domain</b></td>
        <td>string</td>
        <td>
          Domain is the JetStream domain of where the origin stream exists. This is commonly used between a cluster/supercluster and a leaf node/cluster.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersstreamsindexconfigsourcesindexexternal">external</a></b></td>
        <td>object</td>
        <td>
          External is the external stream configuration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>filterSubject</b></td>
        <td>string</td>
        <td>
          FilterSubject is an optional filter subject which will include only messages that match the subject, typically including a wildcard.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startSeq</b></td>
        <td>integer</td>
        <td>
          StartSeq is an optional start sequence the of the origin stream to start mirroring from.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startTime</b></td>
        <td>string</td>
        <td>
          StartTime is an optional message start time to start mirroring from. Any messages that are equal to or greater than the start time will be included. The time format is RFC 3339, e.g. 2023-01-09T14:48:32Z<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.streams[index].config.sources[index].external
<sup><sup>[↩ Parent](#networkspecparametersstreamsindexconfigsourcesindex)</sup></sup>



External is the external stream configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiPrefix</b></td>
        <td>string</td>
        <td>
          APIPrefix is the prefix for the API of the external stream.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>deliverPrefix</b></td>
        <td>string</td>
        <td>
          DeliverPrefix is the prefix for the deliver subject of the external stream.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.subNetworks[index]
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>



Configuration for the sub network

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparameterssubnetworksindexlimits">limits</a></b></td>
        <td>object</td>
        <td>
          Hardware limits for the sub network<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the sub network<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#networkspecparameterssubnetworksindexnodepoolselector">nodepoolSelector</a></b></td>
        <td>object</td>
        <td>
          NodePoolSelector is a label query over nodepool that should match the replica count. It must match the nodepool's labels.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparameterssubnetworksindextolerationsindex">tolerations</a></b></td>
        <td>[]object</td>
        <td>
          Indicates the tolerations the pods under this pool have. A pool's tolerations is not allowed to be updated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.subNetworks[index].limits
<sup><sup>[↩ Parent](#networkspecparameterssubnetworksindex)</sup></sup>



Hardware limits for the sub network

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fileStorage</b></td>
        <td>string</td>
        <td>
          How much disk space is available for data that is stored on disk<br/>
          <br/>
            <i>Default</i>: 1G<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>inMemoryStorage</b></td>
        <td>string</td>
        <td>
          How much memory is available for data that is stored in memory<br/>
          <br/>
            <i>Default</i>: 1G<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.parameters.subNetworks[index].nodepoolSelector
<sup><sup>[↩ Parent](#networkspecparameterssubnetworksindex)</sup></sup>



NodePoolSelector is a label query over nodepool that should match the replica count. It must match the nodepool's labels.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparameterssubnetworksindexnodepoolselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.subNetworks[index].nodepoolSelector.matchExpressions[index]
<sup><sup>[↩ Parent](#networkspecparameterssubnetworksindexnodepoolselector)</sup></sup>



A label selector requirement is a selector that contains values, a key, and an operator that relates the key and values.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.subNetworks[index].tolerations[index]
<sup><sup>[↩ Parent](#networkspecparameterssubnetworksindex)</sup></sup>



The pod this Toleration is attached to tolerates any taint that matches the triple <key,value,effect> using the matching operator <operator>.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>effect</b></td>
        <td>string</td>
        <td>
          Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator represents a key's relationship to the value. Valid operators are Exists and Equal. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tolerationSeconds</b></td>
        <td>integer</td>
        <td>
          TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index]
<sup><sup>[↩ Parent](#networkspecparameters)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparametersusersindexlimits">limits</a></b></td>
        <td>object</td>
        <td>
          The limits for the user<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the user<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersusersindexpermissions">permissions</a></b></td>
        <td>object</td>
        <td>
          The pub/sub permissions for the user<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersusersindexwritetosecret">writeToSecret</a></b></td>
        <td>object</td>
        <td>
          The secret to write the user credentials to<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index].limits
<sup><sup>[↩ Parent](#networkspecparametersusersindex)</sup></sup>



The limits for the user

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>data</b></td>
        <td>integer</td>
        <td>
          Specifies the maximum number of bytes<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>payload</b></td>
        <td>integer</td>
        <td>
          Specifies the maximum message payload<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subscriptions</b></td>
        <td>integer</td>
        <td>
          Specifies the maximum number of subscriptions<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index].permissions
<sup><sup>[↩ Parent](#networkspecparametersusersindex)</sup></sup>



The pub/sub permissions for the user

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#networkspecparametersusersindexpermissionspub">pub</a></b></td>
        <td>object</td>
        <td>
          Specifies the publish permissions<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecparametersusersindexpermissionssub">sub</a></b></td>
        <td>object</td>
        <td>
          Specifies the subscribe permissions<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index].permissions.pub
<sup><sup>[↩ Parent](#networkspecparametersusersindexpermissions)</sup></sup>



Specifies the publish permissions

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allow</b></td>
        <td>[]string</td>
        <td>
          Specifies allowed subjects<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deny</b></td>
        <td>[]string</td>
        <td>
          Specifies denied subjects<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index].permissions.sub
<sup><sup>[↩ Parent](#networkspecparametersusersindexpermissions)</sup></sup>



Specifies the subscribe permissions

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allow</b></td>
        <td>[]string</td>
        <td>
          Specifies allowed subjects<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deny</b></td>
        <td>[]string</td>
        <td>
          Specifies denied subjects<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.parameters.users[index].writeToSecret
<sup><sup>[↩ Parent](#networkspecparametersusersindex)</sup></sup>



The secret to write the user credentials to

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the secret<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.claimRef
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.compositionRef
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.compositionRevisionRef
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.compositionRevisionSelector
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.compositionSelector
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.spec.environmentConfigRefs[index]
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.publishConnectionDetailsTo
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#networkspecpublishconnectiondetailstoconfigref">configRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: map[name:default]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkspecpublishconnectiondetailstometadata">metadata</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.publishConnectionDetailsTo.configRef
<sup><sup>[↩ Parent](#networkspecpublishconnectiondetailsto)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.publishConnectionDetailsTo.metadata
<sup><sup>[↩ Parent](#networkspecpublishconnectiondetailsto)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.resourceRefs[index]
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.spec.writeConnectionSecretToRef
<sup><sup>[↩ Parent](#networkspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Network.status
<sup><sup>[↩ Parent](#network)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>account</b></td>
        <td>string</td>
        <td>
          The UID of the account resource
<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#networkstatusconnectiondetails">connectionDetails</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          The operator for the NATS server
<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>system</b></td>
        <td>string</td>
        <td>
          The UID of the secret user resource
<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.status.conditions[index]
<sup><sup>[↩ Parent](#networkstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Network.status.connectionDetails
<sup><sup>[↩ Parent](#networkstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastPublishedTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
