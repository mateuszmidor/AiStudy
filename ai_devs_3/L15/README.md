# L15

Find shortest connection between 2 persons

## Run

```sh
make db
make run # another terminal
```
, access neo4j at http://localhost:7474, authentication: None, then query to see connections graph:
```neo4j
match(n) return n
```