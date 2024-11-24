#!/bin/sh

echo "Waiting for MySQL Database..."

while ! nc -z mysql 3306; do
  sleep 0.1
done