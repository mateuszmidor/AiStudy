# RAG - Retrieval Augmented Generation


## Run

```sh
ollama run llama3
make install # may take a while, installs lots of Python's packages
make run
```

, output:
```log
(...)
2024/09/14 13:35:31 INFO retrieving information regarding: Which programming language is robust?
2024/09/14 13:35:40 INFO retrieved results="[{Score:0.68962824 Text:Rust is programming language that produces robust programs} {Score:0.5946623 Text:C++ is programming language that produces fast programs} {Score:0.46246603 Text:Python is lame programming language}]"
2024/09/14 13:35:40 INFO prompt: 
Instruction: Based only on the provided information, answer the question in one short sentence.
Information: Rust is programming language that produces robust programs
Information: C++ is programming language that produces fast programs
Information: Python is lame programming language
Question: Which programming language is robust?
2024/09/14 13:35:40 INFO sending prompt to ollama...
2024/09/14 13:35:45 INFO response: Rust is the programming language that produces robust programs.
(...)
```