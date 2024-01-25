# Use Cases

EdgeFarm is a general purpose development platform for edge scenarios. Accordingly, a variety of use cases can be implemented with EdgeFarm.

What makes EdgeFarm different from other solutions is not only that I can implement the users specific use case, but that you can implement it easily and conveniently and maintain it after implementation.

The following (incomplete) list should give an idea what EdgeFarm can be used for.

## Read out sensors or control actors

This is one of the simplest and most common use cases. There are sensors that are accessible by a device using local protocols. However, the data of these sensors need to be distributed to a completely different location, e.g. a cloud system. 

The required application is developed locally on the developers machine and can then be deployed, packaged in a container, with edgefarm.applications to the appropriate edge device with access to the sensors. The transport of my collected data is handled by edgefarm.network.

## Connecting existing air-gapped systems

Let's assume there is a fully functional system that has been developed elaborately and is doing its job. This system generates data that can only be read locally. However, this generated data would be very interesting to know in the meantime.

In this case EdgeFarm can be used to digitalize the device afterwards. For this the edgefarm-ready device  is connected with the target hardware. From this point one can access the target hardware at any time and digitize it bit by bit.

## Optimize application performance

Changes to the firmware of an embedded device are often quite complex. Thus, developers tend to implement as much functionality as possible in the application - even not needed functionality that is used current use case. This leads to a situation where the application simply is too complex for the current use case and tends to transfer data that is not needed.

Due to the separation of EdgeFarm into firmware updates and application updates, individual application components can be changed or delivered at any time. This allows to react to new requirements at any time and enrich, convert or expand data. Only the data that is actually needed should be transferred.

## Deploy Multiple Applications

EdgeFarm has the ability to run multiple applications on one edge device. These applications are isolated from each other and their resources can be limited.

This means that different groups of people can run different applications on the devices, each of which implements its own use cases. The edge device becomes an application platform.

## Preprocess data before sending it to the cloud

Sometimes it is necessary to preprocess data before sending it to the cloud. The preprocessing can be done in the application running on the edge devices and could use a lot of data, preprocess it and send only single events using edgefarm.network. As for preprocessing, a simple algorithm or a complex AI model can be used. That is completely application specific and up to the developer. 

Preprocessing at the edge can drastically reduce the amount of data to be transmitted.