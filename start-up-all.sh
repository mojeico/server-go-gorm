
docker rm -f $(docker ps -a -q)
echo "Docker containers were deleted"
docker rmi $(docker images -a -q)
docker image prune -a --force
echo "Docker images were deleted"


rm -rf /home/trucktrace-user-api/
mkdir  /home/trucktrace-user-api/
cd /home/trucktrace-user-api
git clone git@bitbucket.org:trucktrace/trucktrace-user-api.git .
echo "trucktrace-user-api was cloned"
docker-compose up -d
echo "Docker compose was up"
docker build -t trucktrace-user-api .
echo "Trucktrace-user was built"
docker run --name trucktrace-user-api --network host -d --rm trucktrace-user-api
echo "Trucktrace-user was run"


rm -rf /home/trucktrace-api/
mkdir  /home/trucktrace-api/
cd /home/trucktrace-api
git clone git@bitbucket.org:trucktrace/trucktrace-api.git .
echo "trucktrace-api was cloned"
docker build -t trucktrace-api .
echo "Trucktrace-api was built"
docker run --name trucktrace-api --network host -d --rm trucktrace-api
echo "Trucktrace-api was run"


rm -rf /home/trucktrace-notifications/
mkdir  /home/trucktrace-notifications/
cd /home/trucktrace-notifications/
git clone git@bitbucket.org:trucktrace/trucktrace-notifications.git .
echo "trucktrace-notifications was cloned"
docker build -t trucktrace-notifications .
echo "Trucktrace-notification was built"
docker run --name trucktrace-notifications --network host -d --rm trucktrace-notifications
echo "Trucktrace-notification was run"


rm -rf /home/trucktrace-service/
mkdir  /home/trucktrace-service/
cd /home/trucktrace-service/
git clone git@bitbucket.org:trucktrace/trucktrace-service.git .
echo "trucktrace-service was cloned"
docker build -t trucktrace-service .
echo "Trucktrace-service was built"
docker run --name trucktrace-service --network host -d --rm trucktrace-service
echo "Trucktrace-service was run"

screen -S golang-ui -X quit
cd /home/go-ui/trucktrace-nextjs
git pull
yarn build
screen -dmS golang-ui yarn start
echo "Front was started im screen"
