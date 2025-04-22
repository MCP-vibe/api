@echo off

:: Подгружаем переменные окружения из .env
for /f "tokens=1,2 delims==" %%i in (.env) do set %%i=%%j

:: Запускаем миграцию
goose -dir ./internal/infrastructure/database/migrations postgres "host=%DATABASE_HOST% user=%DATABASE_USER% password=%DATABASE_PASSWORD% dbname=%DATABASE_DATABASE% port=%DATABASE_PORT% sslmode=disable" up
