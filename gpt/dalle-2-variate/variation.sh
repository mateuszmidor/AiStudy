#!/usr/bin/env bash

# https://platform.openai.com/docs/guides/images/variations-dall-e-2-only
# input must be square PNG, max 4MB size
curl https://api.openai.com/v1/images/variations \
  -H "Authorization: Bearer $GPT_APIKEY" \
  -F image='@corgi.png' \
  -F model="dall-e-2" \
  -F n=1 \
  -F size="1024x1024"