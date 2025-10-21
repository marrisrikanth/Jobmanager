# Jobmanager
Remote job execution framework to execute and manage commands
This program should receive a job execution request over REST endpoint and execution it.
REST API Framework to execute commands as jobs on a system. It has capabilities to get the status of a particular job details or list of all the jobs.

It also has capabilities to que to the jobs when the CPU load on the system reaches a threshold limit and execute them later by checking it periodically.

It has capabilities to cancel the pending or running jobs.

All the jobs will get persisted into SQLite db. It gives the capabilities to restart the pending/queued jobs. 

It purge old job records which are older than x hours.

"Need to develop a job execution and status report management system using REST APIs. List down all the tasks required to implement this program in Go." 

HTTPS communication can be enabled by generating certificate like:

mkdir certs; cd certs
openssl -q -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"

go build -ldflags="-s -w" . to remove symbol table and DWARF debugging information for production.


Testing:
To start the jobmanager
export API_KEY=mysecretkey123
go run .

$curl -X POST http://localhost:8080/jobs -d {"command": "date;sleep 4;date"} -H Content-Type: application/json -H X-API-KEY: mysecretkey123

$curl -X GET http://localhost:8080/jobs -H X-API-KEY: mysecretkey123
curl -X GET http://localhost:8080/job/cancel?id=a4aedba8-b8of-42a7-a80f-5df5f11717d6 -H X-API-KEY: mysecretkey123
