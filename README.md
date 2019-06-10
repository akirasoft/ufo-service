
# UFO Service

This service updates status LEDs on the Dynatrace UFO based on the state of a pipeline run.
The service is subscribed to the following keptn events:

- sh.keptn.events.deployment-finished
- sh.keptn.events.evaluation-done
- sh.keptn.events.tests-finished

UFO row will be set based on environment. Dev and stage events trigger TOP row, production events trigger BOTTOM row.

A new-artefact event will result in the LEDs lighting BLUE (0000ff) in a clockwise pattern
A deployment-finished event will result in the LEDs lighting PURPLE (800080) in a clockwise pattern
A tests-finished event will result in the LEDs flashing green (00ff00)
An evaluation-done event will set LEDs based on the status of the evaluation. Pass will set LEDs to GREEN, Fail will set LEDs to RED (ff0000)

## Installation

To use this service, you must have a Dynatrace UFO accessible from your Kubernetes cluster. Additionally, you must have a secret defined with the IP address of your UFO:
```
kubectl -n keptn create secret generic ufo --from-literal="UFO_ADDRESS=<replacewithip>"
```

Afterwards, to install the service in your keptn installation checkout or copy the `ufo-service.yaml`.

Then apply the `ufo-service.yaml` using `kubectl` to create the Dynatrace service and the subscriptions to the keptn channels.

```
kubectl apply -f ufo-service.yaml
```

Expected output:

```
service.serving.knative.dev/ufo-service created
subscription.eventing.knative.dev/ufo-subscription-deployment-finished created
subscription.eventing.knative.dev/ufo-subscription-tests-finished created
subscription.eventing.knative.dev/ufo-subscription-evaluation-done created
```

## Verification of installation

```
$ kubectl get ksvc ufo-service -n keptn
NAME            DOMAIN                               LATESTCREATED         LATESTREADY           READY     REASON
ufo-service   ufo-service.keptn.x.x.x.x.xip.io   ufo-service-dd9km   ufo-service-dd9km   True
```

```
$ kubectl get subscription -n keptn | grep ufo-subscription
ufo-subscription-deployment-finished          True
ufo-subscription-evaluation-done              True
ufo-subscription-tests-finished               True
$
```

When the next event is sent to any of the subscribed keptn channels the UFO LEDs should be adjusted accordingly.

## Uninstall service

To uninstall the dynatrace service and remove the subscriptions to keptn channels execute this command.

```
kubectl delete -f ufo-service.yaml
````