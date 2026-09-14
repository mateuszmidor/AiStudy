class City {
  city string @description("Name of a city")
  country string @description("Country the city is located in")
}

function GetRandomCity(country: string) -> City {
  client GPT4oMini
  prompt #"
    {# start a user message #}
    {{ _.role("user") }}

    Pick one random city in {{ country }}.

    {# special macro to print the output schema instructions. #}
    {{ ctx.output_format }}
  "#
}

test RandomCityTest {
  functions [GetRandomCity]
  args { country "China" }
}