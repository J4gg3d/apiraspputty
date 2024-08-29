#!/bin/bash

# Set variables
TARGET_HOST="admin@raspputty.local"    # Replace with your Raspberry Pi's SSH address
TARGET_DIR="/home/admin/Downloads"      # Directory on Raspberry Pi to store the image and scripts
IMAGE_NAME="restapi-image"                # Name of the Docker image
TAR_FILE="restapi-image.tar"              # Name of the tar file
DOCKER_COMPOSE_FILE="docker-compose.yaml" # Name of your Docker Compose file

# Step 1: Build the Docker image
#docker build -t ${IMAGE_NAME} .

# Step 2: Save the Docker image to a tar file
#docker save -o ${TAR_FILE} ${IMAGE_NAME}

# Step 3: Transfer the tar file and docker-compose.yml to the Raspberry Pi
scp ${TAR_FILE} ${TARGET_HOST}:${TARGET_DIR}/
scp ${DOCKER_COMPOSE_FILE} ${TARGET_HOST}:${TARGET_DIR}/

# Step 4: Create a script on the Raspberry Pi to load the image and run it
ssh -tt ${TARGET_HOST} << 'EOF'
#!/bin/bash

# Set variables
TARGET_DIR="/home/admin/Downloads"
TAR_FILE="restapi-image.tar"
IMAGE_NAME="restapi-image"
DOCKER_COMPOSE_FILE="docker-compose.yaml"

# Navigate to the target directory
cd ${TARGET_DIR}

# Step 1: Load the Docker image
docker load -i restapi-image

# Step 2: Start the Docker container using docker-compose
docker compose up

# Cleanup: Remove the tar file if you want to save space
#rm ${TAR_FILE}

EOF

echo "Deployment script completed."
