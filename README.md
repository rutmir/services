# services

Required architecture:</br>
- MongoDB data storage</br>
- Memcached cache server</br>
- RabbitMQ message server


Required environment parameters

for chat-rest:</br>
- LOG_PATH_SEPARATOR = /services</br>
- MEMCACHE_URL = </br>
- MONGO_URL = </br>

for chat-endpoint:</br>
- LOG_PATH_SEPARATOR = /services</br>
- AMQP_HOST =</br>
- AMQP_PORT =</br>
- AMQP_USERNAME =</br>
- AMQP_PASS =</br>
- MEMCACHE_URL = </br>


for chat-worker:</br>
- LOG_PATH_SEPARATOR = /services</br>
- AMQP_HOST =</br>
- AMQP_PORT =</br>
- AMQP_USERNAME =</br>
- AMQP_PASS =</br>