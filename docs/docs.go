// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API [Roman] Support"
        },
        "license": {
            "name": "romandnk",
            "url": "https://github.com/romandnk/shortener"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/urls": {
            "post": {
                "description": "Create short new URL alias if not exists.",
                "tags": [
                    "URL"
                ],
                "summary": "Create short URL alias",
                "parameters": [
                    {
                        "description": "Required JSON body with original url",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/urlroute.CreateURLAliasRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "URL alias was created successfully",
                        "schema": {
                            "$ref": "#/definitions/urlroute.CreateURLAliasResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    }
                }
            }
        },
        "/urls/:alias": {
            "get": {
                "description": "Get original URL by its alias.",
                "tags": [
                    "URL"
                ],
                "summary": "Get original URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Required path param with url alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Original URL was received successfully",
                        "schema": {
                            "$ref": "#/definitions/urlroute.GetOriginalByAliasResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "httpresponse.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "urlroute.CreateURLAliasRequest": {
            "type": "object",
            "properties": {
                "original_url": {
                    "type": "string"
                }
            }
        },
        "urlroute.CreateURLAliasResponse": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                }
            }
        },
        "urlroute.GetOriginalByAliasResponse": {
            "type": "object",
            "properties": {
                "original_url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1/",
	Schemes:          []string{},
	Title:            "URL shortener project",
	Description:      "Swagger API for Golang Project URL Shortener.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
