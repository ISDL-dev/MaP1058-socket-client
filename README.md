# MaP1058-socket-client
Go socket client capturing bio-signals sent from [MaP1058](https://wp.santeku-map.com/%e3%82%bd%e3%83%95%e3%83%88%e3%82%a6%e3%82%a7%e3%82%a2/%e5%8f%8e%e9%8c%b2%e3%83%bb%e8%a7%a3%e6%9e%90%e3%82%bd%e3%83%95%e3%83%88%e3%82%a6%e3%82%a7%e3%82%a2/map1058/)


## Utility Commands
### transpose command
run the following command to relocate signals in a recorded csv to the right
```bash
go run cmd/transpose/main.go -i <input csv file path> -o <output csv file path> --trim-index true
```
input csv file should be like the following:
```csv
pnt,ch1,ch2,ch3,ch4,ch5,ch6,ch7,ch8
1,123,123,123,123,123,123,123,123
2,456,456,456,456,456,456,456,456
3,789,789,789,789,789,789,789,789
```

the result will be like the following:
```csv
1,2,3
123,456,789
123,456,789
123,456,789
123,456,789
123,456,789
123,456,789
123,456,789
123,456,789
```