#!/usr/bin/env bash

curl -is -X POST http://localhost:8080/correctText -H "Content-Type: application/json" -d '{"text":"Im Walde steht ein Baum."}' -H "Content-Language:de" 

curl -is -X POST http://localhost:8080/correctText -d "Im Walde steht ein Baum." -H "Content-Language:de" 