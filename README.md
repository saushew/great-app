## Как попробовать

```bash
# Клонируем репы
git clone https://github.com/saushew/great-app.git
git clone https://github.com/saushew/great-module-implementation.git

# Собираем модуль
cd great-module-implementation
go build -buildmode=plugin -o executor.so ./executor.go

# Собираем приложение
cd ../great-app
go build -o myapp ./cmd

# Запускаем с указанием пути до модуля через флаг '--module'
./myapp --module ../great-module-implementation/executor.so
