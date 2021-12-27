username="$(kubectl get secret roo-test-default-user -o jsonpath='{.data.username}' | base64 --decode)"
password="$(kubectl get secret roo-test-default-user -o jsonpath='{.data.password}' | base64 --decode)"
service="$(kubectl get service roo-test -o jsonpath='{.spec.clusterIP}')"
kubectl run perf-test --image=pivotalrabbitmq/perf-test -- --uri amqp://$username:$password@$service -s 4000 -x 2 -y 4
kubectl wait --for=condition=ready pod perf-test
kubectl logs -f perf-test