# qdrant - vector database

https://content.aidevs.pl/play/lekcja2_ibrfgh

This demo calls python program to create embedings on local machine. Suggested by chatgpt4.

## Run

```sh
make install # may take a while, installs lots of Python's packages
make run
```

Result:
```log
2024/07/26 10:27:33 INFO add collection name=knowledge dimensions=384
2024/07/26 10:27:41 INFO add point text="Python is kind of snake"
2024/07/26 10:27:49 INFO add point text="C++ is programming language that produces fast programs"
2024/07/26 10:27:57 INFO add point text="Rust is programming language that produces robust programs"
2024/07/26 10:28:05 INFO add point text="Python is lame programming language"
2024/07/26 10:28:13 INFO add point text="Monty Python is a comedy show"
2024/07/26 10:28:21 INFO search text="Which programming language is fast?"
2024/07/26 10:28:21 INFO result payload="C++ is programming language that produces fast programs" score=0.7733528
2024/07/26 10:28:29 INFO search text="Which programming language is robust?"
2024/07/26 10:28:29 INFO result payload="Rust is programming language that produces robust programs" score=0.6896283
2024/07/26 10:28:37 INFO search text="What is Python?"
2024/07/26 10:28:37 INFO result payload="Python is kind of snake" score=0.7788878
2024/07/26 10:28:45 INFO search text="Who is lame?"
2024/07/26 10:28:45 INFO result payload="Python is lame programming language" score=0.4776154
```