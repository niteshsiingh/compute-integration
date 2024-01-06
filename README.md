# Compute Integration

## Initial Setup

1. **Create MongoDB Cluster or use an existing one:**
   - If you don't have a MongoDB cluster, create one on [MongoDB online](https://cloud.mongodb.com/).

2. **Connect Cluster to Controller:**
   - Copy the connection string provided by MongoDB online.
   - Paste the connection string into the `connectionString` variable in `controller/controller.go`.

     ```go
     const connectionString = "YOUR_CONNECTION_STRING"
     ```

3. **Specify Variables in `controller/controller.go`:**
   - Replace `"Your_DATABASE_NAME"` with your desired database name in the `dbName` variable.

     ```go
     const dbName = "Your_DATABASE_NAME"
     ```

   - Replace `"Your_COLLECTION_NAME"` with your desired collection name in the `colName` variable.

     ```go
     const colName = "Your_COLLECTION_NAME"
     ```

   - Replace `"Your_DESIRED_ZONE"` with your desired zone name in the `zone` variable.

     ```go
     const zone = "Your_DESIRED_ZONE"
     ```

   - Replace `"PPROJECT-ID"` with your projectId in the `project` variable.

     ```go
     const project = "PPROJECT-ID"
     ```

   - Replace `"YOUR API-KEY"` with your Google Cloud API Key in the `API_KEY` variable.

     ```go
     const API_KEY = "YOUR API-KEY"
     ```

## Flow

In `main.go`, a router is created using the Gin framework, handling 4 types of requests:

- **GET "/instances":**
  - Fetches all the running and available instances from the cloud platform.

- **PUT "/instance":**
  - Takes a query parameter of `instance_type` and provides all the details of instances of that type.

- **PUT "/instances":**
  - Takes a query parameter of `instance_id` and updates the status key of that instance in the database to terminate and delete the instance.

- **POST "/instance":**
  - Takes a JSON request with keys `"types"` and `"name"` asking for the value of the instance type and instance name. 
  - Creates an instance with that name and type, storing it in the database.

### Functionality Details:

- **GetInstances:**
  - Checks all the details of instances corresponding to the instance type using `getInstanceDetail`.
  - Lists instances whose status is not running.

- **GetInstanceDetail:**
  - Checks if the instance corresponding to the instance type is present in the database.
  - Updates its values if it is present, otherwise puts it into the database.

- **TerminateInstance:**
  - Terminates the instance with the provided instance id.
  - Calculates the cost of the instance and feeds that information into the database.

- **CreateInstance:**
  - Creates an instance with a port between 9000-9100.
  - Stores the instance in the MongoDB database.
