# the pound

## local development

a makefile is present to prop up the core server and a docker-compose.yaml file is present prop up auxiliary backends (postgresql, rabbitmq, elasticsearch etc)

```
$ # load environment variables from .env file
$ source load_env.sh .env
$ # start auxiliary backends
$ docker-compose up -d
$ make core
```

## api

for the first iteration of this backend, i went with a monolithic architecture for simplicity and reduced overhead. my initial intention was to self host and then stress test this in order to see how a distributed, microservices architecture could be beneficial, but i'm putting that off for now

absence of response body link indicates empty response body

#### GET /health

#### POST /register
- creates a user given a username and password
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/register.go#L12-L15)

#### POST /login
- creates a user session, identified by a token, given a username and password
- user authentication is handled using a JWT token sent in the Authorization header of all requests with the eexception of the /health, /register, and /login endpoints. each token has a short lived expiration date defined in the token itself and a longer lived session expiration date defined in the database upon first login. when the short lived expiration date is passed, an updated token is sent in the Authorization header of the response of the corresponding request. as such, it is the client's responsibility to check if the token is present in every response it receives, refreshing their own if present. otherwise, the session is lost and a new one must be created by logging in again. when the longer lived expiration date is passed, the client is forced to create a new session by logging in again
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/login.go#L13-L16)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/login.go#L18-L20)

#### GET /bark
- gets all information for bark given the bark ID
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/bark.go#L94-L96)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/bark.go#L98-L100)

#### POST /bark
- creates a bark with the given text content
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/bark.go#L35-L37)

#### DELETE /bark
- deletes a bark with the given ID
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/bark.go#L94-L96)

#### GET /barks
- gets the barks & rebarks for the dog with the given ID, ordered by descending creation date, using the given offset pagination parameters
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/barks.go#L14-L18)

#### POST /protect
- updates the requesting dog's profile visibility according to the given boolean
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/protect.go#L11-L13)

#### POST /approve
- approves a follow request from the given dog and sets the given notification to read
- i understand that needing the notification ID is a big design flaw. it felt odd to denormalise the notification ID into the treat, rebark, and following tables
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/approve.go#L11-L15)

#### POST /reject
- rejects a follow request from the given dog and sets the given notification to read
- see POST /approve for the explanation of the clear design flaw that the notification ID in the request body schema is
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/reject.go#L11-L14)

#### GET /notifications
- gets unread notifications for the requesting dog using the given offset pagination parameters
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/notifications.go#L14-L17)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/notifications.go#L19-L21)

#### POST /notification_read
- reads the notification with the given ID. an idempotent action
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/notification_read.go#L11-L13)

#### GET /timeline
- gets the barks and rebarks in the requesting dog's timeline, ordered by descending creation date, using the given offset pagination parameters
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/timeline.go#L14-L17)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/timeline.go#L19-L21)

#### POST /paw
- creates a paw to the bark with the given ID with the given text content
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/paw.go#L11-L14)

#### GET /paws
- gets all of the paws to the bark with the given ID, ordered by descending creation date
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/paws.go#L12-L14)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/paws.go#L16-L18)

#### ANY /validate
- validates the JWT token in the Authorization header. like every other protected endpoint, sets the Authorization header of the response with the updated token if it's short lived expiration date has passed

#### GET /thread
- gets the thread for the bark with the given ID. considering a bark as a tree with itself as the root, a thread is the shortest path from any given bark in the tree to the route, ordered by ascending level of depth in the tree. the thread includes (and ends with) the bark with the given ID
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/thread.go#L12-L14)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/thread.go#L16-L18)

#### GET /dog
- gets a dog profile information, either identified by username or ID. if both are present in the query string parameters, the ID is ignored
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/dog.go#L12-L15)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/dog.go#L17-L19)

#### GET /does_follow
- gets the following relationship between the requesting dog and the dog with the given ID whether it exists or not
- [query string parameters](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/does_follow.go#L13-L15)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/does_follow.go#L17-L19)

#### POST /treat
- gives the bark with the given ID a treat if it hasn't already, removes it if it has already
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/treat.go#L12-L14)

#### POST /rebark
- gives the bark with the given ID a rebark if it hasn't already, removes it if it has already
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/rebark.go#L12-L14)

#### POST /follow
- requests to follow the dog with the given ID if a follow relationship doesn't exist. the request is immediately granted if the dog with the given ID is public. removes the follow relationship/request if it exists otherwise
- [request body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/follow.go#L13-L15)
- [response body schema](https://github.com/ayaviri/the-pound/blob/main/apps/backend/cmd/core/follow.go#L17-L19)
