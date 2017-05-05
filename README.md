# usid

Unique sortable identifier

## Example

```go
id := usid.MustNew(usid.Timestamp(time.Now()), usid.SecRndEntropy())
```
