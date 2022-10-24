
### 1. Run e2e

go test
```
# go test -v ./e2e/
```

ginkgo
```
# cd e2e
# ginkgo 
```
```
# cd e2e
# ginkgo build
# ./e2e.test
```

ginkgo: run spec with regex
```
# ginkgo -focus=Statefulset
```

### 2. Report

```
# cat e2e/junit_01.xml
```