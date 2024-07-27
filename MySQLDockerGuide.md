# Pull the MySQL Docker image
docker pull mysql:latest

# Create a Docker volume for persistent data storage
docker volume create mysql-data

# Run MySQL container
docker run --name mysql-container \
    -e MYSQL_ROOT_PASSWORD=your_root_password \
    -e MYSQL_DATABASE=your_database_name \
    -e MYSQL_USER=your_username \
    -e MYSQL_PASSWORD=your_password \
    -p 3306:3306 \
    -v mysql-data:/var/lib/mysql \
    -d mysql:latest

# Check if the container is running
docker ps

# Start the container
docker start mysql-container

# Connect to MySQL from within the container
docker exec -it mysql-container mysql -u root -p

# Connect to MySQL from your host machine
mysql -h 127.0.0.1 -P 3306 -u your_username -p your_password your_database_name

# Stop the container
docker stop mysql-container

# Remove the container (warning: this will delete the container)
docker rm mysql-container