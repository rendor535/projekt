#!/bin/bash

npx concurrently \
"cd backend && go run ./cmd/api" \
"cd frontend && npm run dev -- --host"
