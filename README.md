Calculation Parser with [Pratt Parsing](https://ja.wikipedia.org/wiki/Pratt%E3%83%91%E3%83%BC%E3%82%B5)

### Example

```
$ go run main.go
>> 1 + 2
(1 + 2)
>> 1 * 2 + 3
((1 * 2) * 3)
>> -1 + 2 * 3
((-1) + (2 * 3))
>> 1 + -2 * 3
(1 + ((-2) * 3))
```
