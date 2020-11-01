## fleet simple-demo app

simple-demo app to be used for fleet demo.

The app will generate a different version based on the VERSION build argument pass to the docker build.

`docker build --build-arg VERSION=dev . -t fleet-simple-app:dev`

The manifests for the same are available in the manifests folder.

These can be setup as a gitjob as follows:

```
apiVersion: fleet.cattle.io/v1alpha1
kind: ClusterGroup
  metadata:
    name: simple-app-dev
    namespace: fleet-default
  spec:
    selector:
      matchLabels:
        clusterType: dev
```