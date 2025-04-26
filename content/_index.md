---
title: "Runink: Distributed Pipeline Platform"
---

<br>
<br>
<br>
<br>

<table>
  <tr>
    <th><img src="/images/logo.png" width="250"/></th>
    <th><h4>Runink is a Go-native distributed pipeline orchestration and governance platform.</h4></th>
  </tr>
</table>
<br>
<table>


- **Go + Linux primitives** (cgroups, namespaces, pipes)
- **Slice-based execution** (ephemeral, secure workers)
- **Raft-backed metadata stores** for state consistency
- **Schema contracts** for trusted data evolution
- **Built-in observability, security, and lineage**

Runink empowers you to build **declarative, auditable, and efficient** data pipelines — without the complexity of Kubernetes or JVM-based stacks.

Curious?
[Detailed Architecture →](/docs/architecture/) | [Components Overview →](/docs/components/) | [Comparison with other open-source projects →](/docs/benchmark/)

---


---

## Key Concepts
<br>
<br>
<img src="/images/components.png" width="580"/>

<br>
<br>

---

<table>
  <tr>
    <th><img src="/images/runink.png" width="250"/></th>
    <th><h4>The golang code base to deploy features from configurations files deployed by command actions over the CLI/API.</h4></th>
  </tr>
  <tr>
    <th>Runink</th>
  </tr>
</table>
<br>
<table>
  <tr>
    <th><img src="/images/runi.png" width="250"/></th>
    <th><h4>A single instance of a pipeline step running as an isolated <i>Runi Slice Process</i> managed by a <i>Runi Agent</i> within the constraints of a specific <i>Herd</i></h4></th>
  </tr>
  <tr>
    <th>Runi</th>
  </tr>
</table>
<br>
<table>
  <tr>
    <th><img src="/images/herd.png" width="250"/></th>
    <th><h4>A logical grouping construct, similar to a Kubernetes Namespace, enforced via RBAC policies and resource quotas. Provides multi-tenancy and domain isolation.</h4></th>
  </tr>
  <tr>
    <th>Herd</th>
  </tr>  
</table>
<table>
  <tr>
    <th><img src="/images/barn.png" width="350"/></th>
    <th><h4>A distributed, Raft-backed state store that guarantees strong consistency, high availability, and deterministic orchestration. No split-brain, no guesswork — just fault-tolerant operations.</h4></th>
  </tr>
  <tr>
    <th>Barn</th>
  </tr>  
</table>
