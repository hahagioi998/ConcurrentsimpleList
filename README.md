# ConcurrentsimpleList
一个并发安全的有序链表

```
➜  ConcurrentsimpleList git:(main) ✗ go test -v -race ConcurrentList_test.go
=== RUN   TestIntSet
--- PASS: TestIntSet (25.31s)
PASS
ok      command-line-arguments  26.606s
```

使用 Go 语言完成一个 并发安全的有序链表（数据严格有序并且没有重复元素）

key point
- 一写多读的场景读只需要 atomic
- 一写多读的场景写需要   atomic+lock