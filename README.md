# Secretary

A tool to sync kubernetes secrets and also set image pull secrets for default service accounts in namespaces

## ENV variables
| ENV Variable          | Description                                   | Default |
|-----------------------|-----------------------------------------------|---------|
| `EXCLUDED_NAMESPACES` | A comma separated list of excluded namespaces | Empty   |
| `SYNC_DELAY`          | Delay between each sync in seconds            | `15`    |
|                       |                                               |         |