
# Who?

[Eliot](https://github.com/ernoaapa/eliot) is lightweight Kubernetes like system to manage containerized applications on IoT devices with an emphasis on simplicity, minimality, and usability. The goal is to give a simple way to develop, deploy and manage applications running in IoT device. Comparing to other traditional platforms, Eliot takes modern technologies from cloud to enable more DevOps way of working in IoT world.

Eliot re-uses technologies and concepts from [Cloud Native technologies](https://www.cncf.io/), like [Docker](https://www.docker.com) and [Kubernetes](https://www.kubernetes.io). So if you're familiar with those, you find really easy to get started with Eliot.

## Motivation

I was building modern connected device product what were distributed around the world.
I have over 10 years of software engineer experience with five years of DevOps and faced the problem that there's no state-of-the-art solution for managing connected devices a way that is common nowadays in cloud solutions. 
Most platforms and services focus heavily to the cloud connectivity, data processing, and analysis, but I needed a solution to manage device Operating System and application deployment to build easy to use, modern service for our customers.

Key features needed:
- No network connectivity requirement
- Simple in-device development
- Fast application deployment
- Over-The-Air device management
- Resource allocation and restriction
- Security and software isolation
- Built for IoT from ground up

And that's the day when [Eliot](https://github.com/ernoaapa/eliot) were born ❤︎

## Use cases

Eliot has a small footprint and minimal requirements why it's suitable for a wide range or use cases.
Linux is used nowadays everywhere; info screens, sensors, factories, home IoT, security, cars, etc. and Eliot can support most of them. And thanks to Golang, there's releases available to wide range of architectures.

## Eliot vs. Other

### Docker

[Docker](https://www.docker.com) is a software technology providing containers, promoted by the company [Docker, Inc.](http://www.docker.com/company), and provides full-blown container platform for the cloud environment, taking care of container distribution, orchestration, authentication, infrastructure, etc. Docker has been playing a big role in pushing out the container technology.

At the heart ❤︎ of Docker is [containerd](https://containerd.io), which provides an additional layer of abstraction and automation of operating-system-level virtualization on Windows and Linux. Docker extracted and open sourced to accelerate the innovation across the ecosystem and donated it to open foundation. _And Eliot is based on the containerd!_

Eliot and Docker are not competing against each other; they are working together, in Open Source, to take the ecosystem forward.

### Kubernetes

[Kubernetes](https://www.kubernetes.io) is a great platform for orchestrating containerized software in the cloud environment. Kubernetes have been one of the fastest growing open source projects in past years and new integrations and 3rd party support are popping out every day.

Keys to Kubernetes great success, addition to the great community, are simple concepts and consistent APIs. That's why Eliot is following the great leader. Goal is to re-use Kubernetes code as much as possible, but at the same time keep focus tightly in IoT limitations and requirements.
Kubernetes is pushing the boundaries of cloud computing and Eliot is working on in IoT.

### AWS, Azure, Google, IBM
All cloud IoT solutions ([AWS](https://aws.amazon.com/iot/), [Azure](https://azure.microsoft.com/en-us/suites/iot-suite/), [Google](https://cloud.google.com/solutions/iot/), [IBM](https://www.ibm.com/internet-of-things)) base to the same practice, you use SDK to implement software that collects data from the sensor and send it to the cloud where data gets processed and analyzed. Analysis result can send a message back to the device to trigger some action.

Eliot doesn't try to provide this kind of features at all, actually, you can use any cloud service with Eliot.
Eliot provides an easy way to deliver your cloud integration to the device and gives you a way to update the software safely and easily.

Even better, might be that you don't need to code anything! There might be available open source implementation made by someone in Docker community or you can share your code with the thousands of Docker users around the world with a single command.

### Ansible, Chef, SaltStack

Configuration management tools (e.g. [Ansible](https://www.ansible.com), [Chef](https://www.chef.io), [SaltStack](https://saltstack.com/)) are configuration management tools for managing servers. Communication models are often point-to-point and are limited up to thousands of nodes with a strong focus to make configuration changes as real-time as possible.

These are show stoppers for IoT solutions because you might have hundred thousands or more devices distributed around the world which have poor connection rarely available.

Eliot builds the solution from the ground up for IoT devices with the focus on things that are important in IoT solutions; security, scalability, stability.