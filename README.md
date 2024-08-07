# pasta_bin
Paste bin in Golang!

# Notes
* Using Firestore for post tracking, not currently supporting user profiles.
* `CMD + SHIFT + R` is the way to hard reset the stylesheet when working on it, duh.
* Scheduled job seems to be working with the expiry timer.
* Next step is provisioning cloud resources
  * [ ] IaC (Pulumi)
  * [x] Kubernetes config
    * [ ] Containerization
    * [ ] Kubernetes onboarding
  * [ ] CI/CD