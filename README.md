# go-base
A base service to start new services from

## Getting Started

* Clone to new service folder
* Change the remote to your new service's URL
* Change serviceName in run.sh
* Push

```
git clone git@github.com:divideandconquer/go-base.git new-service
cd new-service
git remote set-url origin git://new.url.here
sed -i '' 's/"base"/"test"/' build/run.sh
git add *
git commit -m "initial clone"
git push -u origin master
```
