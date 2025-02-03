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
            "name": "Ray Toh",
            "url": "https://www.raythx.com",
            "email": "raythx98@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/url": {
            "post": {
                "description": "Given a URL and custom settings, shorten it with a unique alias",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Shorten URL",
                "operationId": "shorten-url",
                "parameters": [
                    {
                        "description": "Shorten URL Request",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ShortenUrlRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Validation Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/url/redirect/{alias}": {
            "post": {
                "description": "Given an alias, redirects request to the full URL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Redirects to full URL",
                "operationId": "redirect-alias",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "303": {
                        "description": "Redirected",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Validation Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.ShortenUrlRequest": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "custom_expiry": {
                    "type": "integer"
                },
                "is_no_expiry": {
                    "type": "boolean"
                },
                "shortened_url": {
                    "type": "string"
                },
                "url": {
                    "description": "TODO: Add custom validation for URL Required",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5051",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL Shortener Server",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
