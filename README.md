# pasta_bin
Paste bin in Golang!

Can run locally on your machine within the `/src` directory with `go run .`
Must provide a .env file with variable `CRED` that has path the the Firebase SDK json you are using.
Will get this hosted publically soon.

# Notes - Still in development
* Using Firestore for post tracking, not currently supporting user profiles.
* `CMD + SHIFT + R` is the way to hard reset the stylesheet when working on it, duh.
* Scheduled job seems to be working with the expiry timer.
* Next step is provisioning cloud resources
  * [ ] IaC (Pulumi)
  * [x] Kubernetes config
    * [ ] Containerization
      * [ ] Trying to get docker volume to store json config
    * [ ] Kubernetes (minikube) onboarding (EKS IS WAY TOO EXPENSIVE BRO) 
  * [ ] CI/CD
