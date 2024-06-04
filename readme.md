# Social Network

This project is a social network webpage. The Project is written in Golang and Vue3.
Audit questions are [here](https://01.kood.tech/git/root/public/src/branch/master/subjects/social-network/audit)

## Table of Contents

- [Social Network](#social-network)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Technologies](#technologies)
  - [Installation](#installation)
    - [With Docker](#with-docker)
    - [Without Docker](#without-docker)
  - [Usage](#usage)
  - [Authors](#authors)


## Features

- Chat service: Allows users to communicate with each other in real-time through a chat interface.
- Post creation: Users can create new posts to start discussions on various topics, share their thoughts with even pictures.
- Commenting: Users can comment on existing threads to share their thoughts and engage in conversations.
- Authentication: Provides user authentication and authorization to secure the application.
- Groups: Users can create and manage groups for various topics to communicate better with likeminded people.
- Events: Members can create events for their groups.
- Following: Users can follow eachother and 

## Technologies

- Vue.js 3: Frontend JavaScript framework for building the user interface.
- Gorilla WebSocket: Backend library for implementing WebSocket communication such as listing active clients and messaging.
- Golang: Handle server-side logic and provide the necessary functionality and data to support the Vue.js frontend application.
- SQLite3 : Store and manage user data, messages, comments and threads in the application.
- Docker: Run the project without worrying about combatability or dependencies.

## Installation

### With Docker

1. Install and run "Docker Desktop" from [www.docker.com](https://www.docker.com/products/docker-desktop/) 

2. Start Docker Desktop (open the program from Windows button or type "systemctl --user start docker-desktop" on Linux)

3. run "docker compose up" to get docker containers running.

4. wait until it is built and loaded. Backend will be on port 8000 and frontend will be on port 5173

5. go to http://localhost:5173 to enjoy our social network

6. when you've finished - stop and delete containers by running 'docker compose down'

7. to delete also all unused images and stopped containers - run 'docker system prune'
   or if you wish to delete only selected images - run 'docker images' to show all images and then 'docker rmi image_id' (replace 'image_id' with ID of your image(s))

### Without Docker

1. Run backend server:
   'cd backend/'
   'go run .'

2. Install the necessary dependencies:
   open new terminal
   'cd frontend/'
   'npm install'

3. Build the Vue app:
   'npm run build'

4. Run frontend:
   'npm run preview'

5. Go to http://localhost:5173 to enjoy our social network

6. When you're finished with audit - stop both backend and frontend by pushing 'Ctrl+C' in both open terminals

## Usage

1. Register a new user or log in with existing credentials.
2. Explore the different features of the application, such as chat, thread creation, and commenting.
3. Interact with other users through the real-time-chat interface.
4. Create new threads to start discussions or participate in existing threads by adding comments.
5. Log out of the application when done.

## Authors

- [@KristjanL](https://01.kood.tech/git/KristjanL)
- [@asdfguy](https://01.kood.tech/git/asdfguy)
- [@Rostislav](https://01.kood.tech/git/Rostislav)
- [@hr.tauno](https://01.kood.tech/git/hr.tauno)
- [@Pihtje](https://01.kood.tech/git/Pihtje)