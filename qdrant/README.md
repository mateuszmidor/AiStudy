# qdrant - vector database

https://content.aidevs.pl/play/lekcja2_ibrfgh

This demo calls python program to create embedings on local machine. Suggested by chatgpt4.

## Run

```sh
make install
make run
```

Result:
```log
2024/07/25 23:59:01 INFO add collection name=dane dimensions=384
2024/07/25 23:59:10 INFO add point collection_name=dane text="Python is kind of snake"
2024/07/25 23:59:19 INFO add point collection_name=dane text="Python is lame programming language"
2024/07/25 23:59:27 INFO add point collection_name=dane text="C++ is programming language that produces fast programs"
2024/07/25 23:59:36 INFO add point collection_name=dane text="Rust is programming language that produces robust programs"
2024/07/25 23:59:45 INFO search collection_name=dane text="Which programming language is fast?"
{"result":[{"id":"9b31733d-aa7a-07e9-71a1-dd8110a83374","version":2,"score":0.7733528,"payload":{"text":"C++ is programming language that produces fast programs"}}],"status":"ok","time":0.001371117}
2024/07/25 23:59:53 INFO search collection_name=dane text="Which programming language is robust?"
{"result":[{"id":"21c04ccd-902c-8d2a-dd2b-459e977c76f6","version":3,"score":0.68962824,"payload":{"text":"Rust is programming language that produces robust programs"}}],"status":"ok","time":0.000554891}
2024/07/26 00:00:02 INFO search collection_name=dane text="What is Python?"
{"result":[{"id":"827dc662-14cc-e55e-a1a5-8ca7644a8fc0","version":0,"score":0.77888787,"payload":{"text":"Python is kind of snake"}}],"status":"ok","time":0.000523922}
```