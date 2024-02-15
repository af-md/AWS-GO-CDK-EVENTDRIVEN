// what do you need to create a rerource: 
// for a table i will need the table resource and data
// resource should have the data values defined
// dynamodb needs name, pk, and the rest can stay default

resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "EventTable"
 hash_key         = "id"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  
  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name        = "dynamodb-table-1-terra"
    Environment = "production"
  }
}