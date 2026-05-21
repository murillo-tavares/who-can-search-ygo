package httpapi

import "net/http"

const openAPISpec = `{
  "openapi": "3.1.0",
  "info": {
    "title": "Who Can Search YGO API",
    "version": "0.1.0",
    "description": "Public API for finding Yu-Gi-Oh! cards whose stored effects can add a selected target card from the Deck to the hand."
  },
  "servers": [
    {
      "url": "http://localhost:18080",
      "description": "Local development"
    }
  ],
  "paths": {
    "/cards": {
      "get": {
        "summary": "Search cards",
        "operationId": "searchCards",
        "parameters": [
          {
            "name": "query",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "minLength": 1
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "minimum": 1,
              "maximum": 50,
              "default": 20
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Matching cards",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/CardSearchResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      }
    },
    "/cards/{id}": {
      "get": {
        "summary": "Get a card",
        "operationId": "getCard",
        "parameters": [
          {
            "$ref": "#/components/parameters/CardID"
          }
        ],
        "responses": {
          "200": {
            "description": "Card details",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Card"
                }
              }
            }
          },
          "404": {
            "$ref": "#/components/responses/CardNotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      }
    },
    "/cards/{id}/searchers": {
      "get": {
        "summary": "Get cards that can search the selected card",
        "description": "Returns active resolved add_deck_to_hand effects whose stored selector matches the selected target card.",
        "operationId": "getCardSearchers",
        "parameters": [
          {
            "$ref": "#/components/parameters/CardID"
          }
        ],
        "responses": {
          "200": {
            "description": "Searcher results",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SearchersResponse"
                }
              }
            }
          },
          "404": {
            "$ref": "#/components/responses/CardNotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      }
    },
    "/healthz": {
      "get": {
        "summary": "Health check",
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "API is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponse"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "parameters": {
      "CardID": {
        "name": "id",
        "in": "path",
        "required": true,
        "schema": {
          "type": "string",
          "format": "uuid"
        }
      }
    },
    "responses": {
      "BadRequest": {
        "description": "Invalid request",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      },
      "CardNotFound": {
        "description": "Card not found",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      },
      "InternalError": {
        "description": "Internal error",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      }
    },
    "schemas": {
      "HealthResponse": {
        "type": "object",
        "required": ["status"],
        "properties": {
          "status": {
            "type": "string",
            "example": "ok"
          }
        }
      },
      "CardSearchResponse": {
        "type": "object",
        "required": ["results"],
        "properties": {
          "results": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/CardSummary"
            }
          }
        }
      },
      "CardSummary": {
        "type": "object",
        "required": ["id", "name"],
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "name": {
            "type": "string"
          },
          "image_url": {
            "type": ["string", "null"],
            "format": "uri"
          }
        }
      },
      "Alias": {
        "type": "object",
        "required": ["alias", "normalized_alias", "alias_kind", "applies_in_zone_codes", "condition_json", "source"],
        "properties": {
          "alias": {
            "type": "string"
          },
          "normalized_alias": {
            "type": "string"
          },
          "alias_kind": {
            "type": "string"
          },
          "applies_in_zone_codes": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "condition_json": {
            "type": "object",
            "additionalProperties": true
          },
          "source": {
            "type": "string"
          }
        }
      },
      "Card": {
        "type": "object",
        "required": ["id", "name", "normalized_name", "aliases", "normalized_aliases", "card_type", "race", "attribute", "monster_categories", "spell_trap_type", "atk", "def", "level", "rank", "link_rating", "archetype", "mentions"],
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "upstream_source": {
            "type": "string"
          },
          "upstream_id": {
            "type": "string"
          },
          "passcode": {
            "type": ["string", "null"]
          },
          "konami_id": {
            "type": ["string", "null"]
          },
          "name": {
            "type": "string"
          },
          "normalized_name": {
            "type": "string"
          },
          "aliases": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/Alias"
            }
          },
          "normalized_aliases": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "description": {
            "type": "string"
          },
          "card_type": {
            "type": ["string", "null"]
          },
          "frame_type": {
            "type": ["string", "null"]
          },
          "race": {
            "type": ["string", "null"]
          },
          "attribute": {
            "type": ["string", "null"]
          },
          "monster_categories": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "spell_trap_type": {
            "type": ["string", "null"]
          },
          "atk": {
            "type": ["integer", "null"]
          },
          "def": {
            "type": ["integer", "null"]
          },
          "level": {
            "type": ["integer", "null"]
          },
          "rank": {
            "type": ["integer", "null"]
          },
          "link_rating": {
            "type": ["integer", "null"]
          },
          "archetype": {
            "type": ["string", "null"]
          },
          "mentions": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "image_url": {
            "type": ["string", "null"],
            "format": "uri"
          }
        }
      },
      "SearchersResponse": {
        "type": "object",
        "required": ["target_card", "effect_code", "results"],
        "properties": {
          "target_card": {
            "$ref": "#/components/schemas/TargetCard"
          },
          "effect_code": {
            "type": "string",
            "enum": ["add_deck_to_hand"]
          },
          "results": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/SearcherResult"
            }
          }
        }
      },
      "TargetCard": {
        "type": "object",
        "required": ["id", "name"],
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "name": {
            "type": "string"
          }
        }
      },
      "SearcherResult": {
        "type": "object",
        "required": ["effect_id", "source_card", "source_text", "action_text"],
        "properties": {
          "effect_id": {
            "type": "string",
            "format": "uuid"
          },
          "source_card": {
            "$ref": "#/components/schemas/CardSummary"
          },
          "source_text": {
            "type": "string"
          },
          "condition_text": {
            "type": ["string", "null"]
          },
          "cost_text": {
            "type": ["string", "null"]
          },
          "action_text": {
            "type": "string"
          },
          "restriction_text": {
            "type": ["string", "null"]
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code", "message"],
            "properties": {
              "code": {
                "type": "string"
              },
              "message": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  }
}`

const docsHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Who Can Search YGO API</title>
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '/openapi.json',
        layout: 'modern',
        theme: 'default',
        agent: {
          disabled: true
        },
        metaData: {
          title: 'Who Can Search YGO API',
          description: 'OpenAPI reference for the Who Can Search YGO backend.'
        }
      })
    </script>
  </body>
</html>`

func (h *Handler) openAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(openAPISpec))
}

func (h *Handler) docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(docsHTML))
}
