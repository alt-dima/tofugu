# Pre-configured Jenkins for TofuGu and Tofugu Toaster Demo

This project provides a way to deploy a pre-configured Jenkins instance to any Kubernetes cluster (especially local ones like Kind or Minikube). It's designed to demonstrate an end-to-end workflow using TofuGu and Tofugu Toaster with a demo OpenTofu unit/tofie.

## Steps to setup

1.  Execute the deployment script:
    ```bash
    bash deploy.sh
    ```

2.  Get the password for the Jenkins admin user:
    ```bash
    kubectl exec --namespace jenkins -it svc/jenkins-dev -c jenkins -- /bin/cat /run/secrets/additional/chart-admin-password && echo
    ```

3.  Start port-forwarding to access the Jenkins UI:
    ```bash
    kubectl port-forward --namespace jenkins svc/jenkins-dev 8080:8080
    ```

4.  Navigate your browser to [http://127.0.0.1:8080](http://127.0.0.1:8080).

5.  Log in with the username `admin` and the password obtained in step 2.

6.  Go to the `tofugu-pipeline` job page: [http://127.0.0.1:8080/job/tofugu-pipeline/](http://127.0.0.1:8080/job/tofugu-pipeline/).

7.  On the left menu, click on **Build**. This will download and execute the `Jenkinsfile` from the repository.

8.  Build #1 should start. Click on the build link to open the console output: [http://127.0.0.1:8080/job/tofugu-pipeline/1/console](http://127.0.0.1:8080/job/tofugu-pipeline/1/console).

9.  Scroll to the end of the console output. You should see a plan output from OpenTofu and a question: "Do you want to apply the changes?".

10. Click **Proceed** and watch as the apply stage is performed.
