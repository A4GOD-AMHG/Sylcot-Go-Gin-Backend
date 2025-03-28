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
            "email": "alexismhgarcia@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/categories": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all categories",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Get categories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Category"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/tasks": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get tasks with optional filters for category, status, and priority",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get filtered tasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Filter by category ID",
                        "name": "categoryId",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter by completion status (true/false)",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "high",
                            "medium",
                            "low"
                        ],
                        "type": "string",
                        "description": "Filter by priority (high/medium/low)",
                        "name": "priority",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Task"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new task for the authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Create a new task",
                "parameters": [
                    {
                        "description": "Task creation data",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TaskRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "details": {
                                    "type": "object"
                                },
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/tasks/{id}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update an existing task's details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Update a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task update data",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Permanently delete a task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Delete a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/tasks/{id}/complete": {
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Toggle a task's completion status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Toggle task status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/auth/forgot-password": {
            "post": {
                "description": "Send a password reset email with a token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Request password reset",
                "parameters": [
                    {
                        "description": "Email for password reset",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Reset email sent",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error: Invalid email",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "error: Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "Registration data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "message: User registered successfully...",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error: Validation failed, details: field errors",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "409": {
                        "description": "error: User already exists",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "error: Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "description": "Validate token and update the password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "Token, new password and confirmation",
                        "name": "resetData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Password updated successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error: Invalid data or token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "error: Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "get": {
                "description": "Validate email verification token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Verify user email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Verification token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Verification success message",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "error: Token required",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "error: Invalid token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "error: Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            }
        },
        "models.Category": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "icon_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Priority": {
            "type": "string",
            "enum": [
                "high",
                "medium",
                "low"
            ],
            "x-enum-varnames": [
                "High",
                "Medium",
                "Low"
            ]
        },
        "models.Task": {
            "type": "object",
            "required": [
                "category_id",
                "title"
            ],
            "properties": {
                "category": {
                    "$ref": "#/definitions/models.Category"
                },
                "category_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "priority": {
                    "$ref": "#/definitions/models.Priority"
                },
                "status": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 3
                },
                "updated_at": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/models.User"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.TaskRequest": {
            "type": "object",
            "required": [
                "category_id",
                "title"
            ],
            "properties": {
                "category_id": {
                    "type": "integer"
                },
                "priority": {
                    "enum": [
                        "high",
                        "medium",
                        "low"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.Priority"
                        }
                    ]
                },
                "title": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3
                }
            }
        },
        "models.User": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_verified": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "refresh_token": {
                    "type": "string"
                },
                "reset_token": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Sylcot API",
	Description:      "Sylcot, that is an acronym for Simplify Your Life by Crossing Out Tasks, it is Task management API to manage your priorities, with a little more functionality and complexity, like JWT authentication",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
