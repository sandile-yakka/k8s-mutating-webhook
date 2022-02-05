# Example Mutating webhook for K8s

## Introduction
This project aims to implement a basic api whose purpose will be to receive admission review requests from the kubernetes API.
Specifically, it'll be expecting to admissions review requests for kubernetes deployments; upon receiving the request it would add resource requests to the first container specified within the deployment. Similar to how you would use kubectl patch

```
kubectl patch deployment example-deployment --type='json' -p='[{"op": "add", "path":"/spec/template/spec/containers/0/resources", "value": {"requests": {"cpu":"200m", "memory": "100Mi"}}}]'
```
After that it will then give an admission review response containing the modified kubernetes object.

### References
- https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
- https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/
- https://www.youtube.com/watch?v=1mNYSn2KMZk&t=414s