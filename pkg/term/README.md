# Term package
Content of this package were copied from
`github.com/kubernetes/kubernetes/pkg/kubectl/util/term`

but modified to remove dependency to following packages:
- `k8s.io/apimachinery/pkg/util/runtime`
- `k8s.io/client-go/tools/remotecommand`

because those caused a lot of more dependencies but were not
required for the functionality.