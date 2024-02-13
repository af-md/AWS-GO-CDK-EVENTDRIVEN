

- cdk bootstrap: create s3 bucket to store the cloudformation template inside your account

# Testing

- cdk synth: to verify that the stack compilers properly. when doing a cdk synth delete your cdk.out folder.

- cdk deploy: to deploy the stack

- cdl destroy: delete all your stack on aws env



# Track
 - Save the event inside a dynamo table
   - DYNAMODB
   - SEND EVENT
   - LOOK AT LAMBDA LOGS
   - WHAT HAPPENED TO LAMBDA'S
 - Create a new lambda that is able to retrieve the events and display them