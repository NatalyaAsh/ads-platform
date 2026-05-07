## 🚀 Выбор следующего шага

1. ### Загрузка изображений (MongoDB)
- Добавить мутацию uploadMedia
- Сохранять файлы на диск, а метаданные в MongoDB
- Расширить тип Ad полем images

2. ### Админ-панель и модерация
- Добавить роли admin (уже есть в Auth Service)
- Реализовать мутации blockUser, moderateAd, updateRole
- Проверять права в Gateway перед вызовом gRPC

3. ### Elasticsearch + RabbitMQ
- Развернуть контейнеры Elasticsearch и RabbitMQ
- Ad Service → публикация событий в очередь
- Новый микросервис indexer → обновление индекса
- gRPC метод searchAds в Gateway

4. ### Kubernetes + Helm
- Написать Dockerfile (если нет)
- Создать Helm-чарты для каждого сервиса
- Развернуть в minikube или kind