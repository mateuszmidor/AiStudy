# L19

Interpret a map.

## Run on macos

```bash
make macos  
> Forwarding  https://a4b0-37-31-42-104.ngrok-free.app
go run . <url>
```

```
2024/11/29 20:51:33 Received: Słuchaj kolego. Lecimy na maksa w prawo, a później ile wlezie w dół. Co tam widzisz?
2024/11/29 20:51:35 openai responded with: {
        "_thinking" : "startuję z pozycji 1,1. przesuwam się na maksa w prawo, czyli na 4,1. Potem ile wlezie w dół, czyli na 4,4.",
        "x": "4",
        "y": "4"
}
2024/11/29 20:51:35 responding with: jaskinia
2024/11/29 20:51:36 Received: Dobra. To co? zaczynamy? Odpalam silniki. Czas na kolejny lot. Jesteś moimi oczami. Lecimy w dół, albo nie! nie! czekaaaaj. Polecimy wiem jak. W prawo i dopiero teraz w dół. Tak będzie OK. Co widzisz?
2024/11/29 20:51:38 openai responded with: {
        "_thinking": "startuję z pozycji 1,1. przesuwam się o 1 w prawo, czyli na pozycję 2,1. Następnie przesuwam się o 1 w dół, czyli na pozycję 2,2.",
        "x": "2",
        "y": "2"
}
2024/11/29 20:51:38 responding with: wiatrak
2024/11/29 20:51:39 Received: Polecimy na sam dół mapy, a później o dwa pola w prawo. Co tam jest?
2024/11/29 20:51:41 openai responded with: {
        "_thinking": "startuję z pozycji 1,1. przesuwam się na maksa w dół, czyli na 1,4. Potem o dwa w prawo, co daje mi pozycję 3,4.",
        "x": "3",
        "y": "4"
}
2024/11/29 20:51:41 responding with: auto
{{FLG:DARKCAVE}}
```