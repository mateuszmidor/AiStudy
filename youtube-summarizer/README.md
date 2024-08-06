# youtube-summarizer

- summarize youtube video based on it's captions
- it may take a minute or two for llama3 to process the prompt
- mind that llama3 context window is max 8k tokens (can handle 15 mins video with polish captions)

## Run

```sh
ollama run llama3
go run .
```

Response:
```text
Prompt:
Summarize the following text, the response MUST be form of bullet points, and MUST be in the same language as source text:
###
Dla niektórych osób jajka są niezastąpionym źródłem składników odżywczych.
Inne zaś osoby unikają ich choćby ze względu na cholesterol.
...<many lines omitted>...
A na dziś z mojej strony to wszystko.
Dziękuję za uwagę i do zobaczenia w kolejnym odcinku.
###

Sending prompt to ollama...
Received response from ollama:
Here is a summary of the text in bullet points:

• Jajka są bogate w kwas pantotenowy, który pomaga utrzymać sprawność umysłową na prawidłowym poziomie i wpływa pozytywnie na urodę.
• W jajkach znajduje się witamina A, która zapewnia prawidłowe widzenie i wspiera odporność.
• Selen jest jednym z kluczowych pierwiastków w jajku, niedoborem którego zmaga się wiele osób. Jeden jajko dziennie pokrywa 30% dobowego zapotrzebowania na selen.
• Selen chroni ciało przed wolnymi rodnikami, usprawnia pracę układu odpornościowego, dba o płodność mężczyzn i pomaga zachować zdrowe włosy i paznokcie.
• Żelazo jest drugim ważnym pierwiastkiem w żółtku jaja, niedoborem którego jest częstą przyczyną anemii.
• Warto spożywać około 5-10 jajek tygodniowo (mniej niż to, jeśli ma się choroby, np. cukrzycę lub choroby układu krążenia).
• Jajka powinny być spożywane w formie ugotowanej na miękko, na twardo lub lekko ściętej jajecznicy.
• Zaleca się unikanie długotrwale smażonych i w głębokim tłuszczu.
• Osobiście autor spożywa średnio 10-12 jajek tygodniowo.
```
