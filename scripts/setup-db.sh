#!/bin/bash

# Database setup script for the Gin API template

set -e

echo "🗄️  Setting up database..."

# Create data directory if it doesn't exist
mkdir -p ./data

# Set default environment
ENV=${APP_ENV:-development}

echo "Environment: $ENV"

case $ENV in
  "production")
    echo "Production environment - database should be set up externally"
    ;;
  "test")
    echo "Setting up test database..."
    export DB_DRIVER=sqlite
    export DB_DSN=./data/test.db
    rm -f ./data/test.db
    echo "Test database cleared"
    ;;
  *)
    echo "Setting up development database..."
    export DB_DRIVER=sqlite
    export DB_DSN=./data/app.db
    
    if [ -f "./data/app.db" ]; then
      echo "Development database already exists"
    else
      echo "Creating new development database..."
    fi
    ;;
esac

echo "✅ Database setup completed!"
