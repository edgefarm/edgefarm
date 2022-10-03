# Seamless edge computing

![](../assets/EdgeFarmLogo_with_text.png)

## What is EdgeFarm?

EdgeFarm is an open source cloud native development platform for hybrid applications where application assets can be freely moved between edge and cloud.

The basis of edgefarm is kubernetes. It is extended with a lot of great open source projects. EdgeFarm combines and extends these selectively to provide a platform that is hardly inferior to the comfort of native cloud development.

## Why EdgeFarm?

How great would it be if I could develop edge software just like cloud software for my kubernetes based cloud backend? I'd be free to try out a new piece of software, I'd have access to a huge pool of open source software, I could use my existing ci/cd system to roll out my edge software, and so on.

But edge computing differs from cloud computing in one fundamental way. While compute in the cloud can be added or replaced automatically at any time, edge devices are tied to specific locations, and replacements or upgrades must be done manually on site. This means that network failures or disconnections cannot simply be bridged by redundancies and taken over by other compute.

This results in the requirement that egde devices must be able to run autonomously over a longer period of time and that the acquired data must be buffered until a connection is available again.

All software used on the edge devices must be able to handle unreliable network connections and synchronize with the backend system when the connection is restored.

If this was solved, my edge device behaving like another kubernetes node, everything needed to deal with unreliable connections already integrated, it would make my programming day a lot nicer.

And that is the reason why EdgeFarm is being developed.