

image_date=$(date +'%m/%d/%Y')
container_data=$(date +'%m-%d-%Y')

random=$RANDOM

docker rm $(docker ps -a | grep trucktrace-service-${container_data})
docker rmi -f $(docker images -a | grep trucktrace-service:${image_date})

docker build -t trucktrace-service:${image_date}:${random} .
docker run --name trucktrace-service-${container_data} --network host -d  trucktrace-service:${image_date}:${random}
