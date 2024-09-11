# Elasticsearch Go CRUD Application

## Overview

This Go application provides a RESTful API to perform CRUD (Create, Read, Update, Delete) operations on book documents stored in Elasticsearch. It utilizes the `olivere/elastic` package for Elasticsearch interactions and `gorilla/mux` for routing.

## Features

- **Create**: Add a new book document by ID
- **Read**: Retrieve a book document by ID
- **Update**: Modify a book document by ID
- **Delete**: Remove a book document by ID

## Prerequisites

- **Go**: Version 1.18 or higher
- **Docker**: For running Elasticsearch
- **Elasticsearch**: Running on `http://localhost:[PORT]`
- **.env File**: For environment variables

## Dependency
- github.com/gorilla/mux v1.8.1
- github.com/olivere/elastic/v7 v7.0.32
- github.com/joho/godotenv v1.5.1
- github.com/josharian/intern v1.0.0
- github.com/mailru/easyjson v0.7.7
- github.com/pkg/errors v0.9.1

## Add .env File

- elastic_url=http://localhost:[Port]

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd <repository-directory>
