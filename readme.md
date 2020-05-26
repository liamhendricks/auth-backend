Starter repo for golang projects using cobra / viper.

Step 1: Update Makefile image name, and Dockerfile cmd path.

Step 2: Update docker-compose.yml `volume` path and `working_dir` path.

Step 3: Update import path for cmd package in main.go.

Step 4: Run  `make setup` then `make local`. You should see "init" before service stops.

After adding any import, run `make deps` to pull / vendor it.
