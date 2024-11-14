# the pound

## local development

```
$ # Load environment variables from .env file
$ source load_env.sh .env
$ # Start auxiliary backends (postgres, rabbitmq, etc)
$ docker-compose up -d
$ make core
```
