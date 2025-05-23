// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Steven Poon",
            "url": "https://github.com/RYANCOAL9999",
            "email": "lmf242003@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/challenges": {
            "get": {
                "description": "Retrieves a list of recent challenges based on the provided limit. Returns the most recent challenge if there are multiple results.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "challenges"
                ],
                "summary": "List recent challenges",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Maximum number of challenges to retrieve",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of recent challenges or the most recent challenge",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Challenge"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid input data",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error during retrieval",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/challenges/join": {
            "post": {
                "description": "Allows a player to join a new challenge, provided they haven't participated in the last minute. It processes the challenge creation within a transaction, updates the prize pool, and starts a background task to calculate the challenge result after 30 seconds. Returns the status of the challenge creation.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "challenges"
                ],
                "summary": "Join a challenge",
                "parameters": [
                    {
                        "description": "Details for joining the challenge",
                        "name": "challenge",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NewChallengeNeed"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Challenge joined successfully, returns the status of the challenge, it represent as number, 1 is joined, 0 is Ready",
                        "schema": {
                            "$ref": "#/definitions/models.JoinChallengeResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request due to invalid input data",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "425": {
                        "description": "Too many requests if attempting to join within a minute",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error during challenge creation or transaction",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Challenge": {
            "type": "object",
            "required": [
                "amount",
                "created_at",
                "player_id",
                "probability",
                "status",
                "won"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "player_id": {
                    "type": "string"
                },
                "probability": {
                    "type": "number"
                },
                "status": {
                    "$ref": "#/definitions/models.Status"
                },
                "won": {
                    "type": "boolean"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.JoinChallengeResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "$ref": "#/definitions/models.Status"
                }
            }
        },
        "models.NewChallengeNeed": {
            "type": "object",
            "required": [
                "amount",
                "player_id"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "player_id": {
                    "type": "integer"
                }
            }
        },
        "models.Status": {
            "type": "integer",
            "enum": [
                0,
                1
            ],
            "x-enum-varnames": [
                "Ready",
                "Joined"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             ":8085",
	BasePath:         "/v2",
	Schemes:          []string{},
	Title:            "Endless Challenge System API",
	Description:      "This is a endless challenge system server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
