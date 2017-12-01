# go-base

**Note** this has been replaced with a revised version in the [healthimation](https://github.com/healthimation/go-base) organization.  The code in this repo is deprecated.

A base service to start new services from

## Getting Started

* Clone to new service folder
* Change the remote to your new service's URL
* Change serviceName in run.sh
* Push

```sh
git clone git@github.com:divideandconquer/go-base.git new-service
cd new-service
git remote set-url origin git://new.url.here
sed -i '' 's/"base"/"new"/' build/run.sh
git add *
git commit -m "initial clone"
git push -u origin master
```

This service infrastructure assumes you are running it in docker and have access to consul.  I recommend
using [coreos-vagrant](https://github.com/coreos/coreos-vagrant) and then running consul as a 
[fleet unit](https://gist.github.com/divideandconquer/08405a4fb597319d3c3e).  You should be able to copy
the unit file at the previous link on to the coreos box and start it with:

```sh
fleetctl start /path/to/consul.service
```
