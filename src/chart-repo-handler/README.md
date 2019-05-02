# Chart Repo Handler - a swift http frontend for helm
This project aims to let unmodified http clients to access Openstack Swift object store.
Furthermore it acts as a helm chart repo with storage capabilities.

## Motivation
The motivation is that usually object stores (like S3) has API for static HTTP content hosting for standard HTTP clients (like web browsers). As Swift uses custom HTTP headers for authentication, it is not working OOTB.
This swift behavior makes impossible to use swift directly with helm so this project is a middle layer between them.

## Usage
The frontend is configured through environment variables. These are the credentials for the object store and the targeted container in it.
To enable TLS for Swift backend, use https in AUTHURL and set the location of RootCA pem file in the TLSCAPATH variable.
See env.sh for a sample.


## Details
* Written in go, targeted to be a 12 factor app
* Uses ncw/swift as the most mature client API towards the Swift object store
* Uses gorilla/mux to decode HTTP requests
* Borrows some of the helm packages to maintain index.yaml file for the charts


This work is licensed under a Creative Commons Attribution 4.0 International License.
