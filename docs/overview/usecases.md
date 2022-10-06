# Use Cases

EdgeFarm is a general purpose development platform for edge scenarios. Accordingly, a variety of use cases can be implemented with EdgeFarm.

What makes EdgeFarm different from other solutions is not only that I can implement my specific use case, but that I can implement it easily and conveniently and maintain it after implementation.

The following (incomplete) list should give an idea for what EdgeFarm can be used for.

## Read out sensors or controll actors

This is one of the simplest and most common use cases. I have sensors or actuators running somewhere, accessible e.g. via local protocols, and need access at a completely different location.

The required application is developed locally on my machine and can then be deployed, packaged in a container, with edgefarm.applications to the appropriate edge device with access to the sensors. The transport of my collected data is handled by edgefarm.network.

## Digitize existing air-gabbed systems

Let's assume I have a fully functional system that has been elaborately developed and is doing its job. This system generates data that I can only read locally, but would be very interesting for me in the meantime.

In this case I could use EdgeFarm to digitize my device afterwards. For this I connect my edgefarm-ready device with the target hardware. From this point I can access the target hardware at any time and digitize it bit by bit.

## Preproduce and reduce application data

Often, changes to an embedded device are quite complex. This can result in a lot of useless data being transferred, because one tries to think about future requirements at the same time.

Due to the separation of EdgeFarm into firmware update and application update, individual application components can be changed or delivered at any time. This allows me to react to new requirements at any time and enrich, convert or expand data, and only what is actually needed is transferred.

## Deploy Multible Applications

EdgeFarm has the ability to run multiple applications on one edge device. These applications are isolated from each other and their resources can be limited.

This means that different groups of people can run different applications on the devices, each of which implements its own use cases, and the edge device becomes an application platform.
