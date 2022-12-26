# tasks-web-service

## Motivation
The goal through this simple application is to train the best practises in terms of developing a kubernetes-ready servcie following the clean architecture pattern and up to date dev standards. 
This will be shown through a simple task management service offering basic CRUD operations to manage your tasks. 
## Usage
### Running the server
1- Setup the environment variables specified in the ".env" file, an example is given in the source code.

In the next steps, we will make use of the Makefile in the code to facilitate the commands.
You can see the different possibilities by running the command:
```
make help
```

2- Spin up a Postgres DB. There is a docker compose file already existing for this purpose in case you have docker already installed. You can run:
```bash 
make run_postgres
```
3- Run the server:
```bash
make run
```
### Testing
Using the following you can run the unit tests from the root of the project:
```bash
make test
```

### API documentation
An extensive API specification is provided using [Swagger](https://github.com/swaggo/swag), you can find it in the `docs`folder.
Using the `docs/swagger.yaml` or `docs/swagger.json` files you can easily form your API requests.  \
Also, you can import the file to Postman, and it will automatically set up the API requests.

For an up-to-date version, you can always regenerate the updated Swagger 2.0 docs using the following command:
```bash
make swagger_generate 
```
Swagger also offers a nice UI which allows you to visualize and interact with the API’s resources without having any of the implementation logic in place.
You can access the UI by running the separate swagger server in `cmd/swagger/main.go` using the following command:
```bash
make swagger_ui
```
After running the server, browse to http://localhost:5000/swagger/index.html and will see Swagger 2.0 Api documentation UI.



### Consuming an endpoint
When the server is running, you can consume the endpoints through your favourite http client, an example would be:
```bash
curl --location --request POST 'http://localhost:8080/v1/api/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title":       "some-title",
    "description": "a task that is much needed ",
    "priority":    8,
    "status":      "New"
}'
```
## Running the app on k8s:
The easiest way to do this is using helm, with the already configured charts in the repo. A prerequisite of course is that you have already spinned up a kubernetes cluster.
Then to deploy the app and its dependencies (the database) to kubernetes it suffises to run those steps:

#### Build the Docker Image
- Build the docker image of the app, through the already present Dockerfile:
```
docker build -t username/image_name:tag_name .
```
- Push the image to your docker registry.
```
docker push username/image_name:tag_name
```
#### Deploy it on k8s
- Update the /k8s/helm/task-service-app/Values.yaml with your parameters, most importantly the image you pushed:
This yould be an example snippet:
```yamlreplicaCount: 1
namespace: task-service-app

image:
  repository: fy984/tasks-web-service
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v1.1.0"
```
- Install the helm charts for the DB and the app, by going into the `./k8s/helm` directory and running those commands:
```bash
# to deploy the database needed by the app
helm install postgresql postgresql --namespace <your-chosen-ns>
# deploy the app
helm install task-service-app task-service-app
```


