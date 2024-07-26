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
2024/07/26 12:31:17 INFO add collection name=knowledge dimensions=384
2024/07/26 12:31:28 INFO add point text="Python is kind of snake"
2024/07/26 12:31:37 INFO add point text="C++ is programming language that produces fast programs"
2024/07/26 12:31:46 INFO add point text="Rust is programming language that produces robust programs"
2024/07/26 12:31:56 INFO add point text="Python is lame programming language"
2024/07/26 12:32:05 INFO add point text="Monty Python is a comedy show"
2024/07/26 12:32:14 INFO search text="Which programming language is fast?"
2024/07/26 12:32:14 INFO result payload="C++ is programming language that produces fast programs" score=0.7733528
2024/07/26 12:32:24 INFO search text="Which programming language is robust?"
2024/07/26 12:32:25 INFO result payload="Rust is programming language that produces robust programs" score=0.68962824
2024/07/26 12:32:34 INFO search text="What is Python?"
2024/07/26 12:32:34 INFO result payload="Python is kind of snake" score=0.77888775
2024/07/26 12:32:44 INFO search text="Who is lame?"
2024/07/26 12:32:44 INFO result payload="Python is lame programming language" score=0.4776155
2024/07/26 12:32:53 INFO search text="Do you know any comedy shows?"
2024/07/26 12:32:53 INFO result payload="Monty Python is a comedy show" score=0.43849182
```