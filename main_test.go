package main

import (
	"encoding/json"
	"os"
	"testing"

	"deli/alice"
)

func BenchmarkHandler(b *testing.B) {
	_ = os.Setenv("GEOCODER_API_KEY", "278f85ec-4fdc-4107-aa20-9832af3d91e2")
	reqJSON := []byte(`
{
  "meta": {
    "locale": "ru-RU",
    "timezone": "UTC",
    "client_id": "ru.yandex.searchplugin/7.16 (none none; android 4.4.2)",
    "interfaces": {
      "screen": {},
      "payments": {},
      "account_linking": {}
    }
  },
  "session": {
    "message_id": 2,
    "session_id": "cf6ceab0-1e42-4067-a665-1e857213c086",
    "skill_id": "187ea0d1-bb4b-4004-9e41-6d456f6d0627",
    "user": {
      "user_id": "D295AFBA42C611E6812CBF856E8AB0BEEFF127E3AA7974EA0026BADC97C4762C"
    },
    "application": {
      "application_id": "EB094B5353930448BD85E68DF3F9B73C391F49E8A71477C814A1AF4FC4264234"
    },
    "user_id": "EB094B5353930448BD85E68DF3F9B73C391F49E8A71477C814A1AF4FC4264234",
    "new": false
  },
  "request": {
    "command": "найти каршеринг проспект мира в радиусе 300 м",
    "original_utterance": "найти каршеринг проспект мира в радиусе 300м",
    "nlu": {
      "tokens": [
        "найти",
        "каршеринг",
        "проспект",
        "мира",
        "в",
        "радиусе",
        "300",
        "м"
      ],
      "entities": [],
      "intents": {
        "cars.nearby": {
          "slots": {
            "address": {
              "type": "YANDEX.STRING",
              "tokens": {
                "start": 2,
                "end": 4
              },
              "value": "проспект мира"
            },
            "measure": {
              "type": "MeasureUnit",
              "tokens": {
                "start": 7,
                "end": 8
              },
              "value": "meter"
            },
            "radius": {
              "type": "YANDEX.NUMBER",
              "tokens": {
                "start": 6,
                "end": 7
              },
              "value": 300
            }
          }
        }
      }
    },
    "markup": {
      "dangerous_context": false
    },
    "type": "SimpleUtterance"
  },
  "version": "1.0"
}`)

	req := new(alice.Request)
	_ = json.Unmarshal(reqJSON, &req)
	var res *alice.Response
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ = Handler(req)
	}
	b.ReportAllocs()
	b.StopTimer()
	resJSON, _ := json.MarshalIndent(res, "", "   ")
	b.Log(string(resJSON))
}
